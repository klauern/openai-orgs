package mcp

import (
	"context"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func AddTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool(
		"list_projects",
		mcp.WithDescription("Lists all projects for the authenticated user"),
	), handleListProjects)
}

func handleListProjects(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	authToken := ctx.Value(authToken{}).(string)

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)

	projects, err := client.ListProjects(100, "", false)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(projects.String()), nil
}
