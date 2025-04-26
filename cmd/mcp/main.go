package main

import (
	"log"

	"github.com/klauern/openai-orgs/pkg/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := mcp.NewMCPServer()

	err := server.ServeStdio(mcpServer, server.WithStdioContextFunc(mcp.AuthFromEnvironment))
	if err != nil {
		log.Fatalf("Error serving server: %v", err)
	}
}
