package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/klauern/openai-orgs/pkg/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mcpServer := mcp.NewMCPServer(ctx)

	err := server.ServeStdio(mcpServer, server.WithStdioContextFunc(mcp.AuthFromEnvironment))
	if err != nil {
		log.Fatalf("Error serving server: %v", err)
	}
}
