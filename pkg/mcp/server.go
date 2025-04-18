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

type Server struct {
	*server.MCPServer
}

func NewServer() (*Server, error) {
	mcpServer := server.NewMCPServer(
		serverName,
		version,
		server.WithInstructions(instructions),
		server.WithLogging(),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)

	return &Server{mcpServer}, nil
}
