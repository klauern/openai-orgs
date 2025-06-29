# MCP Context: Model Context Protocol (MCP)

## Overview

The Model Context Protocol (MCP) is a specification and implementation pattern for building extensible, testable, and robust server components, tools, and resources. In this project, all MCP-related context, guidelines, and best practices are maintained exclusively in the memory-bank. The memory-bank serves as the single source of truth for MCP implementation, project evolution, and documentation.

## Key Guidelines

### Server Implementation

- Always provide server name and version
- Implement error handling for all server operations
- Use stdio transport for local development and testing

### Resources

- Use URI pattern: `scheme://path`
- Implement both static and dynamic resources
- Provide resource descriptions and MIME types
- Use resource templates for parameterized resources

### Tools

- Each tool must have:
  1. Unique name
  2. Clear description
  3. Well-defined parameter schema
  4. Error handling
- Use `mcp.NewTool()` for tool creation
- Implement tool handlers with context support

### Code Style

- Use descriptive error messages with context
- Wrap errors using `fmt.Errorf("failed to %s: %w", action, err)`
- Validate all input parameters before processing
- Use required, description, and enum helpers for parameters

### Testing

- Test each tool and resource handler independently
- Mock external dependencies
- Test error and edge cases
- Use context cancellation for timeout tests
- Integration tests should verify protocol compliance and error handling

### Best Practices

- Use semantic URI schemes and logical resource paths
- Cache resource results when appropriate
- Implement proper cleanup in resource handlers
- Keep tools focused and single-purpose
- Return structured results when possible
- Document side effects in tool descriptions

### Security

- Validate all input parameters
- Implement authentication as required
- Sanitize resource paths
- Limit tool capabilities appropriately

### Common Commands

- Start server: `server.ServeStdio(s)`
- Add tool: `s.AddTool(tool, handler)`
- Add resource: `s.AddResource(resource, handler)`
- List tools: `s.ListTools()`

### Required Imports

```go
import (
    "context"
    "fmt"
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)
```

## Integration in This Project

- MCP is used for OpenAI org management via CLI and API
- All tools and resources follow MCP handler and parameter schema patterns
- Interfaces are used for testability and dependency injection
- All code and tests must comply with these guidelines
- The memory-bank is the authoritative location for all MCP context, guidelines, and updates. All contributors must consult and update the memory-bank to reflect changes, decisions, and best practices.

---

*Update this file as the MCP context evolves or as new best practices are adopted. The memory-bank is the canonical reference for all project context.*
