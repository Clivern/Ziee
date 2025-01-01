FROM golang:1.26.5 AS builder

ARG ACTX0_VERSION=0.1.0
ARG ACTX0_COMMIT=none
ARG ACTX0_BUILD_DATE=unknown
ARG ACTX0_BUILT_BY=docker

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${ACTX0_VERSION} -X main.commit=${ACTX0_COMMIT} -X main.date=${ACTX0_BUILD_DATE} -X main.builtBy=${ACTX0_BUILT_BY}" \
    -o /out/actx0 .

RUN mkdir -p /out/configs /out/var/logs && \
    cp config.dist.yml /out/configs/config.dist.yml

FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

COPY --from=builder --chown=65532:65532 /out/actx0 /app/actx0
COPY --from=builder --chown=65532:65532 /out/configs /app/configs
COPY --from=builder --chown=65532:65532 /out/var /app/var

EXPOSE 8080

VOLUME ["/app/configs", "/app/var"]

USER 65532:65532

ENTRYPOINT ["/app/actx0"]
CMD ["server", "-c", "/app/configs/config.dist.yml"]
