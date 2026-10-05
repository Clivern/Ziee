## Sandbox image

Docker image for the Ziee sandbox agent runtime.

### Published tags

| Tag | Image |
| --- | --- |
| `v0.1.0` | [`clivern/sandbox:v0.1.0`](https://hub.docker.com/r/clivern/sandbox) |
| `latest` | [`clivern/sandbox:latest`](https://hub.docker.com/r/clivern/sandbox) |

```bash
docker pull clivern/sandbox:v0.1.0
docker pull clivern/sandbox:latest
```

### Build and push

From this directory:

```bash
docker build -t clivern/sandbox:v0.1.0 .
docker tag clivern/sandbox:v0.1.0 clivern/sandbox:latest
docker push clivern/sandbox:v0.1.0
docker push clivern/sandbox:latest
```

### Layout

- `Dockerfile` — image definition
- `entrypoint.sh` — container entrypoint
- `bridge.py` — RPC bridge (`/usr/local/bin/rpc`)
