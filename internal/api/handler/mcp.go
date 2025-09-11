package handler

import (
	"context"
	"encoding/json"

	"github.com/hungpdn/mcp2grule/internal/api/dto"
	"github.com/hungpdn/mcp2grule/internal/grule"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPHandler is the handler for MCP API
type MCPHandler struct {
	grule grule.IGrule
}

// NewMCPHandler creates a new MCPHandler
func NewMCPHandler(grule grule.IGrule) *MCPHandler {
	return &MCPHandler{grule: grule}
}

// Evaluate handles the Evaluate API call
func (h *MCPHandler) Evaluate(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in dto.EvaluateIn,
) (*mcp.CallToolResult, *dto.EvaluateOut, error) {

	out, err := h.grule.Evaluate(ctx, in)
	if err != nil {
		return nil, nil, err
	}

	text, _ := json.Marshal(out)

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(text),
			},
		},
	}
	return result, out, nil
}

// Create handles the Create API call
func (h *MCPHandler) Create(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in dto.CreateIn,
) (*mcp.CallToolResult, *dto.CreateOut, error) {

	out, err := h.grule.Create(ctx, in)
	if err != nil {
		return nil, nil, err
	}

	text, _ := json.Marshal(out)

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(text),
			},
		},
	}
	return result, out, nil
}

// Update handles the Update API call
func (h *MCPHandler) Update(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in dto.UpdateIn,
) (*mcp.CallToolResult, *dto.UpdateOut, error) {

	out, err := h.grule.Update(ctx, in.Name, in)
	if err != nil {
		return nil, nil, err
	}

	text, _ := json.Marshal(out)

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(text),
			},
		},
	}
	return result, out, nil
}

// Delete handles the Delete API call
func (h *MCPHandler) Delete(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in dto.DeleteIn,
) (*mcp.CallToolResult, *dto.DeleteOut, error) {

	out, err := h.grule.Delete(ctx, in.Name)
	if err != nil {
		return nil, nil, err
	}

	text, _ := json.Marshal(out)

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(text),
			},
		},
	}
	return result, out, nil
}

// GetAll handles the GetAll API call
func (h *MCPHandler) GetAll(
	ctx context.Context,
	req *mcp.CallToolRequest,
	_ any,
) (*mcp.CallToolResult, *dto.GetAllOut, error) {

	out, err := h.grule.GetAll(ctx)
	if err != nil {
		return nil, nil, err
	}

	text, _ := json.Marshal(out)

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(text),
			},
		},
	}
	return result, out, nil
}

// GetByName handles the GetByName API call
func (h *MCPHandler) GetByName(
	ctx context.Context,
	req *mcp.CallToolRequest,
	in dto.GetByNameIn,
) (*mcp.CallToolResult, *dto.GetByNameOut, error) {

	out, err := h.grule.GetByName(ctx, in.Name)
	if err != nil {
		return nil, nil, err
	}

	text, _ := json.Marshal(out)

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(text),
			},
		},
	}
	return result, out, nil
}
