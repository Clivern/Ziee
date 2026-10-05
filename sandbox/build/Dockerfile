# Copyright 2026 Clivern. All rights reserved.
# License can be found in the LICENSE file.

FROM node:24-bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends git ca-certificates python3 \
    && rm -rf /var/lib/apt/lists/* \
    && npm install -g --ignore-scripts @earendil-works/pi-coding-agent

WORKDIR /repo

COPY entrypoint.sh /usr/local/bin/entrypoint
COPY bridge.py /usr/local/bin/rpc
RUN chmod +x /usr/local/bin/entrypoint /usr/local/bin/rpc

EXPOSE 8765

ENTRYPOINT ["/usr/local/bin/entrypoint"]
