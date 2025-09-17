# mcp2grule

mcp2grule is a production-ready, high-performance Go server that exposes the Grule rule engine to Large Language Models (LLMs) and other AI applications via the Model Context Protocol (MCP).

## Quick start

1. Clone and build:

    ```bash
    git clone https://github.com/hungpdn/mcp2grule.git
    cd mcp2grule
    go mod tidy
    go build -o mcp2grule
    ```

2. Run the server (defaults to stdio transport):

    ```bash
    ./mcp2grule server
    ```

3. Or run with HTTP transport for easy testing (streamable-HTTP):

    ```bash
    MCP_TRANSPORT=streamable-http DATABASE_TYPE=memory ./mcp2grule server
    ```

## Configuration

All configuration is provided via environment variables (see `internal/config/config.go`). The minimal values you may set locally are in `.env.example`:

```bash
MCP_TRANSPORT=streamable-http
DATABASE_TYPE=memory
```

Key env vars

- `MCP_TRANSPORT`: `stdio`, `sse`, or `streamable-http` (default: `stdio`)
- `DATABASE_TYPE`: `memory`, `sqlite`, or `postgresql` (default: `memory`)
- `HTTP_HOST` / `HTTP_PORT`: used for SSE / streamable-http transports

## Project Structure

```text
mcp2grule/
├─ cmd/                # CLI entrypoint (wires services and starts server)
├─ internal/
│  ├─ api/
│  │  ├─ server.go     # MCP transport selection (stdio / sse / streamable-http), tool registration, graceful shutdown
│  │  └─ tool.go       # MCP tool registration (grule.evaluate, grule.create, ...)
│  │  └─ handler/      # MCP handlers that map requests to domain DTOs
│  ├─ grule/
│  │  └─ grule.go      # Domain service: constructs grule engine, Evaluate/Create/Update/Delete logic
│  ├─ storage/
│  │  ├─ storage.go    # IRulesetStorage interface and common errors
│  │  ├─ memory.go     # In-memory ruleset storage (default for local dev)
│  │  └─ postgres.go   # Postgres driver (optional; implements IRulesetStorage)
│  ├─ config/
│  │  └─ config.go     # Environment variable parsing and typed config
│  └─ pkg/
│     ├─ exitcode/     # Canonical exit codes for CLI/startup failures
│     └─ logger/       # Logging helpers and context wiring
├─ ...
└─ README.md           
```

## MCP tools provided

The server registers these MCP tools (exact names used by clients):

- `grule.evaluate` - Evaluate facts against a named ruleset
- `grule.create` - Create a new ruleset
- `grule.update` - Update an existing ruleset
- `grule.delete` - Delete ruleset by name
- `grule.list` - List all rulesets
- `grule.detail` - Get ruleset details by name

See `internal/api/tool.go` for the registration and `internal/api/handler/mcp.go` for request/response handling examples.

## Linters & formatting

This repo uses `golangci-lint`. A starter config is present at `.golangci.yml`. Run:

```bash
golangci-lint run
```

## Docker

A multi-stage `Dockerfile` is included. It builds a static binary and copies it into a minimal distroless runtime image. Example build:

```bash
docker build -t mcp2grule:latest .
docker run -e MCP_TRANSPORT=streamable-http -e DATABASE_TYPE=memory -p 9000:9000 mcp2grule:latest
```

Security note: builder images can contain OS-level CVEs. Consider scanning images in CI and using patched base images.

## Testing

There are currently no unit tests in the repo. Recommended next steps:

- Add unit tests for `internal/storage` and `internal/grule` (happy path + error conditions).
- Add a CI workflow to run `go test ./...` and `golangci-lint run` on PRs.

## TODO

- [ ] Add tests.
- [ ] Add CI.
- [ ] Add pprofing.
- [ ] Add middleware for auth, metrics and logging.
- [ ] Add postgres storage.
- [ ] Add migration.
