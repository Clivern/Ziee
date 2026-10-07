// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeBridge mimics bridge.py: auth line first, then newline-delimited JSON both ways.
// handle runs for each command; reply writes one raw line back. A nil handle hangs up.
func fakeBridge(t *testing.T, apiKey string, handle func(cmd Event, reply func(string))) (string, int) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { listener.Close() })

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		scanner := bufio.NewScanner(conn)
		reply := func(line string) { conn.Write([]byte(line + "\n")) }
		next := func() (Event, bool) {
			if !scanner.Scan() {
				return nil, false
			}
			cmd := Event{}
			json.Unmarshal(scanner.Bytes(), &cmd)
			return cmd, true
		}

		auth, ok := next()
		if !ok || auth["apiKey"] != apiKey {
			reply(`{"type":"auth","success":false,"error":"invalid api key"}`)
			return
		}
		reply(`{"type":"auth","success":true}`)

		for {
			cmd, ok := next()
			if !ok || handle == nil {
				return
			}
			handle(cmd, reply)
		}
	}()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.IP.String(), addr.Port
}

func connect(t *testing.T, host string, port int, apiKey string) *Client {
	t.Helper()

	client := NewClient(host, port, apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, client.Connect(ctx))
	t.Cleanup(func() { client.Close() })

	return client
}

func TestUnitSandboxClient(t *testing.T) {
	t.Run("Connect rejects a bad api key", func(t *testing.T) {
		host, port := fakeBridge(t, "secret", nil)

		err := NewClient(host, port, "wrong").Connect(context.Background())
		assert.ErrorContains(t, err, "invalid api key")
	})

	t.Run("Prompt before Connect", func(t *testing.T) {
		_, err := NewClient("127.0.0.1", 1, "k").Prompt(context.Background(), "hi")
		assert.ErrorContains(t, err, "not connected")
	})

	t.Run("Prompt returns the agent text", func(t *testing.T) {
		host, port := fakeBridge(t, "secret", func(cmd Event, reply func(string)) {
			assert.Equal(t, "prompt", cmd.Type())
			assert.Equal(t, "hi", cmd["message"])
			reply(fmt.Sprintf(`{"type":"response","id":%q,"success":true}`, cmd["id"]))
			reply(`not json`)
			reply(``)
			reply(`{"type":"message_update","assistantMessageEvent":{"type":"text_delta","delta":"hello "}}`)
			reply(`{"type":"message_update","assistantMessageEvent":{"type":"thinking_delta","delta":"hmm"}}`)
			reply(`{"type":"message_update","assistantMessageEvent":{"type":"text_delta","delta":"world"}}`)
			reply(`{"type":"agent_settled"}`)
		})
		client := connect(t, host, port, "secret")

		text, err := client.Prompt(context.Background(), "hi")
		require.NoError(t, err)
		assert.Equal(t, "hello world", text)
	})

	t.Run("Prompt rejected", func(t *testing.T) {
		host, port := fakeBridge(t, "secret", func(cmd Event, reply func(string)) {
			reply(fmt.Sprintf(`{"type":"response","id":%q,"success":false,"error":"busy"}`, cmd["id"]))
		})
		client := connect(t, host, port, "secret")

		_, err := client.Prompt(context.Background(), "hi")
		assert.EqualError(t, err, "busy")
	})

	t.Run("Request returns the matching response", func(t *testing.T) {
		host, port := fakeBridge(t, "secret", func(cmd Event, reply func(string)) {
			assert.Equal(t, "get_session_stats", cmd.Type())
			reply(`{"type":"response","id":"other","success":true}`)
			reply(fmt.Sprintf(`{"type":"response","id":%q,"success":true,"data":{"tokens":{"total":18}}}`, cmd["id"]))
		})
		client := connect(t, host, port, "secret")

		reply, err := client.Request(context.Background(), Event{"type": "get_session_stats"})
		require.NoError(t, err)
		assert.NotEqual(t, "other", reply["id"])
		assert.Equal(t, map[string]any{"tokens": map[string]any{"total": float64(18)}}, reply["data"])
	})

	t.Run("Request failure", func(t *testing.T) {
		host, port := fakeBridge(t, "secret", func(cmd Event, reply func(string)) {
			reply(fmt.Sprintf(`{"type":"response","id":%q,"success":false}`, cmd["id"]))
		})
		client := connect(t, host, port, "secret")

		_, err := client.Request(context.Background(), Event{"type": "get_state"})
		assert.EqualError(t, err, "get_state failed")
	})

	t.Run("Request when the bridge hangs up", func(t *testing.T) {
		host, port := fakeBridge(t, "secret", nil)
		client := connect(t, host, port, "secret")

		_, err := client.Request(context.Background(), Event{"type": "get_state"})
		assert.ErrorContains(t, err, "connection closed")
	})

	t.Run("Request times out", func(t *testing.T) {
		host, port := fakeBridge(t, "secret", func(Event, func(string)) {})
		client := connect(t, host, port, "secret")

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		_, err := client.Request(ctx, Event{"type": "get_state"})
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
