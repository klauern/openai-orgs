package main

import (
	"fmt"
	"os"

	"github.com/klauern/openai-orgs/internal/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := mcp.ServerConfig{
		Name:    "OpenAI Orgs MCP Server",
		Version: "1.0.0",
	}

	s, err := mcp.NewServer(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating server: %v\n", err)
		os.Exit(1)
	}

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
