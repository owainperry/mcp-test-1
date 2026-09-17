# syntax=docker/dockerfile:1

# ---- Build stage ----
# Runs natively on the builder and cross-compiles, so multi-arch builds do not
# need emulation.
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build

WORKDIR /src

# Download dependencies first so this layer caches independently of source edits.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build \
      -trimpath \
      -ldflags="-s -w -X main.version=${VERSION}" \
      -o /out/mcp-test-1 .

# ---- Runtime stage ----
FROM alpine:3.23

RUN adduser -D -u 10001 nonroot

COPY --from=build /out/mcp-test-1 /usr/local/bin/mcp-test-1

USER nonroot:nonroot

# The server speaks MCP over the streamable HTTP transport at /mcp.
ENV MCP_ADDR=:8080
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/mcp-test-1"]
