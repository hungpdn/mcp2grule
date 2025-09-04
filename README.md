# mcp2grule

A minimal **MCP server** in Go that evaluates **Grule** rules.

- SDK: [`github.com/modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk)
- Rules: [`github.com/hyperjumptech/grule-rule-engine`](https://github.com/hyperjumptech/grule-rule-engine)

## Install

```bash
git clone https://github.com/hungpdn/mcp2grule
cd mcp2grule
go mod tidy
go build -o mcp2grule
```

## Run

As a standalone MCP server over stdio

```bash
./mcp2grule
```

Then connect from an MCP client (e.g., Claude Desktop custom server, Postman, ..) using a stdio command.

## Example Request/Response

Request

```json
{
    "jsonrpc":"2.0",
    "id":"1",
    "method": "tools/call",
    "params": {
        "name": "evaluate_grule",
        "arguments": {
            "facts": {
                "total": 150
            },
            "grl": "rule Discount \"if total > 100\" salience 10 { when In.Get(\"total\") > 100 then Out.Set(\"discount\", true); Retract(\"Discount\"); }"
        }
    }
}
```

Response

```json
{
    "jsonrpc": "2.0",
    "id": "1",
    "result": {
        "discount": true
    }
}
```
