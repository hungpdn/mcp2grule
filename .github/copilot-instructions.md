Short summary

This repository implements a small MCP (Model Context Protocol) server in Go that evaluates Grule rules. Key responsibilities:
- Accept MCP tool calls (stdio / SSE / streamable-http transports)
- Store and manage rulesets via a pluggable storage interface
- Evaluate facts against GRL rules using a grule engine

Quick entry points

- Build: `go build -o mcp2grule` (project root). See `go.mod` for Go version and dependencies.
- Run (stdio mode, default): `./mcp2grule server` or just `./mcp2grule` after building. The `server` command is defined in `cmd/server.go` and wired via `cmd/root.go`.
- Env: configuration is driven entirely by environment variables parsed in `internal/config/config.go`. See `.env.example` for the minimal variables (`MCP_TRANSPORT`, `DATABASE_TYPE`).

Architecture / important files

- `cmd/` — CLI entrypoints. `server.go` constructs the dependencies (storage, grule engine, handler) and starts the server.
- `internal/api/server.go` — wraps the `modelcontextprotocol/go-sdk` server and picks transport (stdio / sse / streamable-http). Use this file to understand how the app starts and shuts down.
- `internal/api/tool.go` — registers MCP tools (grule.evaluate, grule.create, grule.update, grule.delete, grule.list, grule.detail). Use these names when composing MCP requests.
- `internal/api/handler/` — handler methods translate MCP requests to domain DTOs and call `internal/grule` service.
- `internal/grule/grule.go` — business logic; constructs grule engine and exposes Create/Update/Evaluate operations.
- `internal/storage/` — storage interfaces and implementations. `storage/storage.go` declares `IRulesetStorage`; `memory.go` is an in-memory implementation used by default.
- `internal/config/config.go` — central config parsed via `github.com/caarlos0/env/v11`. All env var names and defaults live here.
- `internal/pkg/exitcode/exitcode.go` — canonical exit codes used across CLI and startup errors. Prefer these constants when adding scripts or new command code.

Conventions and patterns (specific)

- Dependency injection by constructor: `cmd/server.go` wires concrete implementations (e.g., `storage.NewMemory()` into `grule.New(...)`). When adding a DB backend, implement `IRulesetStorage` and swap here.
- Single responsibility services: the `grule` package contains orchestration and calls into `engine.IGruleEngine` (from `github.com/hungpdn/grule-plus/engine`). Keep rule evaluation and storage logic separated.
- MCP tool handlers always return a `*mcp.CallToolResult` for the transport and a typed DTO for internal flows. See `internal/api/handler/mcp.go` for serialization examples (they marshal DTOs into TextContent).
- Config via env: mutating config at runtime is not supported. Tests and local runs should set env vars (or use `direnv` / `envrc`) before starting the server.

Run & debug tips (project-specific)

- To run using streamable HTTP transport (useful for HTTP clients):

  MCP_TRANSPORT=streamable-http DATABASE_TYPE=memory ./mcp2grule server

- To run in SSE mode (server exposes SSE endpoint): set `MCP_TRANSPORT=sse` and configure `HTTP_HOST` / `HTTP_PORT`.
- To run locally during development: `go run ./cmd` or `go run ./cmd/*.go server` — the `cmd` package is the entrypoint.
- To test MCP flows manually: send JSON-RPC 2.0 calls to the transport used (stdio uses stdio, streamable-http accepts HTTP requests). Example request format is in `README copy.md`.

Files to edit for common changes

- Add a new storage driver: implement `IRulesetStorage` (in `internal/storage`) and update `cmd/server.go` switch on `config.App.DatabaseType`.
- Add a new MCP tool: register it in `internal/api/tool.go` and implement the handler in `internal/api/handler/`.
- Add metrics / auth middleware for HTTP transports: modify `internal/api/server.go` to wrap HTTP handlers (there's a TODO for auth in SSE path).

Examples (copyable snippets)

- MCP tool names (use these exact strings in requests):
  - `grule.evaluate` — evaluate facts
  - `grule.create` — create a ruleset
  - `grule.update` — update an existing ruleset
  - `grule.delete` — delete by name
  - `grule.list` — list all rules
  - `grule.detail` — get by name

- Minimal env for local dev (from `.env.example`):

  MCP_TRANSPORT=streamable-http
  DATABASE_TYPE=memory

Why this layout matters

- Small, composable binary: `cmd/server.go` builds everything at startup and the service runs as one binary. That favors simple deployments (containers that run one process) and local stdio integrations (MCP clients using stdio).
- Clear separation of transport vs domain: `internal/api/server.go` (transport) vs `internal/grule` (domain) lets you add transports or change rule engine implementations with minimal cross-cutting changes.

Notes for AI agents

- Prefer editing or adding code inside `internal/` and `cmd/` — avoid touching top-level scaffolding unless adding CI/Docker changes.
- When changing startup wiring, update `cmd/server.go` and ensure appropriate exit codes from `internal/pkg/exitcode` are used for failure paths.
- Use the `IRulesetStorage` interface and `internal/api/handler/mcp.go` patterns when creating new features that need persistence or MCP tooling.

If any of the above is unclear or you'd like examples for adding a Postgres storage driver or a new MCP tool, tell me which one and I'll add a short patch and example tests.
