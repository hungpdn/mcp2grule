package api

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (s *Server) AddTools() {

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "grule.evaluate",
		Description: "Evaluate facts against the rule set",
	}, s.mcpHandler.Evaluate)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "grule.create",
		Description: "Create a new rule",
	}, s.mcpHandler.Create)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "grule.update",
		Description: "Update an existing rule",
	}, s.mcpHandler.Update)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "grule.delete",
		Description: "Delete a rule by name",
	}, s.mcpHandler.Delete)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "grule.list",
		Description: "List all existing rules",
	}, s.mcpHandler.GetAll)

	mcp.AddTool(s.server, &mcp.Tool{
		Name:        "grule.detail",
		Description: "Read an existing rule by name",
	}, s.mcpHandler.GetByName)

}
