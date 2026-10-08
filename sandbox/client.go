// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package sandbox

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Event is a JSON record sent to or received from the sandbox RPC bridge.
type Event map[string]any

// Type returns the event type.
func (e Event) Type() string {
	t, _ := e["type"].(string)

	return t
}

// Client talks to a sandbox container over its newline-delimited JSON RPC bridge.
// Calls are synchronous and must not run concurrently. After a call fails,
// the stream may be mid-reply, so Close the client and connect again.
type Client struct {
	host    string
	port    int
	apiKey  string
	conn    net.Conn
	scanner *bufio.Scanner
}

// NewClient returns a client for the bridge at host:port.
func NewClient(host string, port int, apiKey string) *Client {
	return &Client{host: host, port: port, apiKey: apiKey}
}

// Connect dials the bridge and authenticates.
func (c *Client) Connect(ctx context.Context) error {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(c.host, strconv.Itoa(c.port)))
	if err != nil {
		return err
	}

	c.conn = conn
	c.scanner = bufio.NewScanner(conn)
	c.scanner.Buffer(nil, 64<<20)

	reply, err := c.Call(ctx, Event{"type": "auth", "apiKey": c.apiKey}, func(e Event) bool {
		return e.Type() == "auth"
	})
	if err == nil {
		err = replyError(reply, "auth failed")
	}
	if err != nil {
		c.Close()
		return fmt.Errorf("sandbox rpc auth: %w", err)
	}

	return nil
}

// Close closes the connection.
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

// Prompt sends a prompt, waits for the agent to finish and returns its text reply.
func (c *Client) Prompt(ctx context.Context, message string) (string, error) {
	id := uuid.NewString()
	var text strings.Builder
	var failure error

	_, err := c.Call(ctx, Event{"id": id, "type": "prompt", "message": message}, func(e Event) bool {
		switch e.Type() {
		case "response":
			if e["id"] == id {
				failure = replyError(e, "prompt failed")
			}

			return failure != nil
		case "message_update":
			delta, _ := e["assistantMessageEvent"].(map[string]any)
			if delta["type"] == "text_delta" {
				s, _ := delta["delta"].(string)
				text.WriteString(s)
			}
		case "agent_settled":
			return true
		}

		return false
	})
	if err != nil {
		return "", err
	}
	if failure != nil {
		return "", failure
	}

	return text.String(), nil
}

// Request sends a command such as {"type": "get_session_stats"} and returns its response.
// It returns an error if the response reports failure.
func (c *Client) Request(ctx context.Context, command Event) (Event, error) {
	id := uuid.NewString()
	if v, ok := command["id"].(string); ok {
		id = v
	}

	out := Event{"id": id}
	for k, v := range command {
		if k != "id" {
			out[k] = v
		}
	}

	reply, err := c.Call(ctx, out, func(e Event) bool {
		return e.Type() == "response" && e["id"] == id
	})
	if err != nil {
		return nil, err
	}

	return reply, replyError(reply, fmt.Sprintf("%v failed", command["type"]))
}

// Call writes command and reads records until done returns true for one, which it returns.
func (c *Client) Call(ctx context.Context, command Event, done func(Event) bool) (Event, error) {
	if c.conn == nil {
		return nil, errors.New("sandbox rpc: not connected")
	}

	// Unblock reads and writes when ctx ends.
	stop := context.AfterFunc(ctx, func() { c.conn.SetDeadline(time.Now()) })
	defer stop()

	data, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Write(append(data, '\n'))
	if err != nil {
		return nil, ctxErr(ctx, err)
	}

	for c.scanner.Scan() {
		event := Event{}
		if json.Unmarshal(c.scanner.Bytes(), &event) != nil {
			continue
		}
		if done(event) {
			return event, nil
		}
	}

	err = c.scanner.Err()
	if err == nil {
		err = errors.New("sandbox rpc: connection closed")
	}

	return nil, ctxErr(ctx, err)
}

// ctxErr prefers the context error when ctx ending caused err.
func ctxErr(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return err
}

// replyError returns the error in a {"success": false, "error": "..."} reply.
func replyError(reply Event, fallback string) error {
	if success, _ := reply["success"].(bool); success {
		return nil
	}
	if msg, _ := reply["error"].(string); msg != "" {
		return errors.New(msg)
	}

	return errors.New(fallback)
}
