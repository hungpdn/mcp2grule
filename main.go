package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Bag is a dynamic key/value store exposed to GRL as "In" (read) and "Out" (write).
type Bag struct{ M map[string]any }

func NewBag(m map[string]any) *Bag {
	if m == nil {
		m = map[string]any{}
	}
	return &Bag{M: m}
}
func (b *Bag) Get(key string) any    { return b.M[key] }
func (b *Bag) Has(key string) bool   { _, ok := b.M[key]; return ok }
func (b *Bag) Set(key string, v any) { b.M[key] = v }
func (b *Bag) Keys() []string {
	ks := make([]string, 0, len(b.M))
	for k := range b.M {
		ks = append(ks, k)
	}
	return ks
}
func (b *Bag) AsMap() map[string]any { return b.M }

// --- Tool input schema & handler types ---

type EvalParams struct {
	GRL   string         `json:"grl"   jsonschema:"Grule rules as a single GRL string"`
	Facts map[string]any `json:"facts" jsonschema:"Input facts: arbitrary JSON object"`
	// Optional: knowledge base identifiers (useful when you later split rules per ruleset/version)
	Ruleset string `json:"ruleset,omitempty" jsonschema:"Optional ruleset name (defaults to 'Session')"`
	Version string `json:"version,omitempty" jsonschema:"Optional ruleset version (defaults to 'v1')"`
}

func evaluateOnce(ctx context.Context, p EvalParams) (map[string]any, error) {
	if p.GRL == "" {
		return nil, errors.New("missing grl")
	}
	rs := p.Ruleset
	if rs == "" {
		rs = "Session"
	}
	ver := p.Version
	if ver == "" {
		ver = "v1"
	}

	// Build rules into a KnowledgeLibrary (per-call for simplicity; cache in prod).
	kbLib := ast.NewKnowledgeLibrary()
	rb := builder.NewRuleBuilder(kbLib)
	res := pkg.NewBytesResource([]byte(p.GRL))
	if err := rb.BuildRuleFromResource(rs, ver, res); err != nil {
		log.Printf("build rule err: %v", err)
		return nil, fmt.Errorf("build rule failed: %w", err)
	}

	// DataContext with In / Out bags.
	dataCtx := ast.NewDataContext()
	in := NewBag(p.Facts)
	out := NewBag(nil)
	if err := dataCtx.Add("In", in); err != nil {
		return nil, err
	}
	if err := dataCtx.Add("Out", out); err != nil {
		return nil, err
	}

	// Execute on a fresh KB instance.
	kb, _ := kbLib.NewKnowledgeBaseInstance(rs, ver)
	gr := engine.NewGruleEngine()
	if err := gr.ExecuteWithContext(ctx, dataCtx, kb); err != nil {
		return nil, fmt.Errorf("execute failed: %w", err)
	}
	return out.AsMap(), nil
}

func main() {
	impl := &mcp.Implementation{Name: "mcp2grule", Version: "0.0.1"}

	// Create the MCP server (stdio transport by default).
	server := mcp.NewServer(impl, nil)

	// Register the evaluate_grule tool with typed arguments.
	mcp.AddTool(server, &mcp.Tool{
		Name:        "evaluate_grule",
		Description: "Evaluate a GRL rule against a JSON facts object; use In.Get/Has and Out.Set in GRL",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args EvalParams) (*mcp.CallToolResult, any, error) {
		result, err := evaluateOnce(ctx, args)
		if err != nil {
			return nil, nil, err
		}

		// Return JSON to the client. The SDK content type supports plain text;
		// we return a JSON string for portability.
		payload, _ := json.MarshalIndent(result, "", "  ")
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(payload)},
			},
		}, nil, nil
	})

	// Run over stdio (works great with Claude Desktop custom servers, or mcp clients).
	// You can also wire your own transport (HTTP/SSE) later with the SDK.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}

}
