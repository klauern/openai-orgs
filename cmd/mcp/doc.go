/*
Package main provides an MCP (Model Context Protocol) server implementation for managing OpenAI organizations.

It implements a stdio-based transport following the MCP specification and handles
authentication through environment variables. The server provides tools and resources
for OpenAI organization management while maintaining secure access control.

The server requires appropriate environment variables for authentication as documented
in the project README. To start the server:

	go run cmd/mcp/main.go

The implementation uses mcp.NewMCPServer to initialize the server instance and
serves it over stdio transport with proper error handling and logging. The server
supports standard MCP features including:

  - Stdio-based transport for local development
  - Environment-based authentication
  - Organization management tools and resources
  - Structured error handling with context

See https://modelcontextprotocol.io for MCP specification details.
*/
package main
