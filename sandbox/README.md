## Sandbox

Run an agent container that reaches `OpenRouter` through `Ziee`, then talk to it over `RPC`.


### Start the Sandbox

```bash
docker run -d --name sandbox \
  -p 8765:8765 \
  -v /path/to/repo:/repo \
  -e RUN_ID=my-run-id \
  -e RPC_API_KEY=my-rpc-key \
  -e PROXY_URL=https://<ziee-host>/api/v1/sandbox/<sandbox-id>/api \
  -e PI_MODEL=openrouter/anthropic/claude-sonnet-4.5 \
  clivern/sandbox:v0.1.0
```


### Call it

```bash
RPC_API_KEY=my-rpc-key python3 client.py --host 127.0.0.1 --port 8765 "what does this repo do?"
```

```bash
RPC_API_KEY=my-rpc-key python3 client.py --host 127.0.0.1 --port 8765
```
