package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName   = "openai-orgs"
	version      = "0.0.1"
	instructions = `
	You are a helpful assistant that can help with OpenAI organization management, which includes managing projects, members, limits, and billing.
	`
)

func NewMCPServer() *server.MCPServer {
	mcpServer := server.NewMCPServer(
		serverName,
		version,
		server.WithInstructions(instructions),
		server.WithLogging(),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	AddTools(mcpServer)
	AddResources(mcpServer)
	return mcpServer
}
