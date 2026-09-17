# mcp-test-1

A deliberately minimal [MCP](https://modelcontextprotocol.io) server, for testing an MCP gateway.

It exposes one tool:

| Tool    | Input            | Output                  |
| ------- | ---------------- | ----------------------- |
| `hello` | `text` (string)  | `Hello, <text>!`        |

Empty or whitespace-only input returns `Hello!`.

## Transport

Streamable HTTP, via [`github.com/modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk) v1.8.0,
which negotiates MCP protocol version `2026-07-28` and falls back to `2025-11-25`,
`2025-06-18`, `2025-03-26` or `2024-11-05` if the client asks for an older one.

| Path       | Purpose                       |
| ---------- | ----------------------------- |
| `/mcp`     | MCP endpoint (point the gateway here) |
| `/healthz` | Plain-text liveness probe     |

The listen address comes from `MCP_ADDR` (default `:8080`).

## Run it

```sh
go run .                                  # listens on :8080
MCP_ADDR=:9000 go run .                   # or pick a port
```

```sh
docker build -t mcp-test-1 .
docker run --rm -p 8080:8080 mcp-test-1
```

Quick smoke test:

```sh
curl -s localhost:8080/healthz
```

## CI

`.github/workflows/build.yml` vets, tests and builds the image on every push and
pull request. On pushes to `main` it also pushes to Docker Hub as:

- `owainperry/mcp-test-1:<github.run_number>` — the build number
- `owainperry/mcp-test-1:latest`

The run number is also compiled into the binary and reported as the server
version over MCP.

### Required secrets

Add these in **Settings → Secrets and variables → Actions**:

| Secret               | Value                                        |
| -------------------- | -------------------------------------------- |
| `DOCKERHUB_USERNAME` | Your Docker Hub username                      |
| `DOCKERHUB_TOKEN`    | A Docker Hub access token (not your password) |

Until they exist, pushes to `main` will fail at the login step; pull request
builds work without them.
