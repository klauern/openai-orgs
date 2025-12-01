# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.
Following the [agents.md](https://agents.md/) standard for universal agent configuration.

# OpenAI Orgs Project Guide

## Project Overview

This is a Go CLI tool and API client for managing OpenAI Platform Administration APIs. It provides comprehensive management capabilities for OpenAI organizations including projects, users, API keys, service accounts, certificates, audit logs, and more. The project also includes an MCP (Model Context Protocol) server for AI assistant integration.

## Common Commands

- Build: `task build` or `go build -v ./...`
- Install: `task install`
- Lint: `task lint` or `golangci-lint run`
- Format: `task fmt` or `gofmt -s -w -e .`
- Test all: `task test` or `go test -v -coverprofile=coverage.out -timeout=120s -parallel=10 ./...`
- Test single file: `go test -v ./path/to/file_test.go`
- Test specific test: `go test -v ./... -run "TestName"`
- Test coverage: `task cover`
- MCP build: `task mcp:build`
- MCP test: `task mcp:test`
- MCP dev: `task mcp:dev` (requires npm/npx for inspector)

## Architecture Overview

### Core Components

- **Root Package**: Generic HTTP client with built-in rate limiting and retry logic using `resty/v2`
- **`/cmd/openai-orgs`**: Main CLI application with subcommands for each OpenAI API resource
- **`/cmd/mcp`**: Model Context Protocol server implementation
- **`/pkg/mcp`**: MCP-specific utilities and tools

### API Client Design

The core API client (`api_client.go`) uses Go generics for type-safe operations:
- `Get[T]`, `GetSingle[T]`, `Post[T]`, `Delete[T]` for common HTTP operations
- Conservative retry strategy: 20 retries, 5-second wait, max 5-minute backoff
- Generic `ListResponse[T]` for paginated API responses

### Resource Coverage

**Organization Level**: users, invites, admin API keys, certificates, audit logs
**Project Level**: projects, project users, project API keys, project service accounts, project rate limits, project certificates

Each resource follows consistent patterns: list, create, retrieve, modify, delete operations.

## Code Style Guidelines

- **Testing**: Use helper functions for test setup/teardown, mock external APIs with `jarcoal/httpmock`
  - NEVER use testify for test generation (per `.cursorrules`)
  - Test helpers are centralized in [test_helpers.go](test_helpers.go)
  - **Test Helper Pattern**: Use `newTestHelper(t)` which provides:
    - Pre-configured mocked client with retries disabled
    - `mockResponse(method, endpoint, statusCode, response)` - registers mock HTTP responses
    - `assertRequest(method, endpoint, times)` - verifies expected request counts
    - `cleanup()` - must be called via `defer` to reset mocks after tests
- Follow Go error handling pattern (check err != nil)
- **Error handling**: Use `fmt.Errorf` with context wrapping, e.g., `fmt.Errorf("error making request: %v", err)`
- **Comments**: Document exported functions and types with meaningful comments
- CLI commands should be organized in subpackages under [cmd/](cmd/)
- All API endpoints should have corresponding CLI commands
- **Imports**: Standard Go import organization (stdlib first, then external)
  - Core dependencies: `resty/v2` for HTTP, `urfave/cli/v3` for CLI, `httpmock` for testing
- **Naming**:
  - Types: PascalCase (e.g., `Client`, `ListResponse`)
  - Constants: Use prefix conventions (e.g., `OwnerTypeUser`)
  - Functions: PascalCase for exported, camelCase for internal
  - Parameters: lowercase, exported fields: camelCase
- **Types**: Use generics for common operations, strongly type constants with custom types
- Use `UnixSeconds` type for timestamp handling

## Project Structure

- `/cmd/openai-orgs` - Main CLI commands and entry point
- `/cmd/mcp` - Model Context Protocol server
- `/pkg/mcp` - MCP implementation details
- Root package - Core API client implementation with individual resource files
- `interfaces.go` - Complete API interface contract for dependency injection
- `types.go` - Central type definitions and constants
- `test_helpers.go` - Shared testing utilities and HTTP mocking patterns

## Development Patterns

- **Interface Design**: `OpenAIOrgsClient` interface in [interfaces.go](interfaces.go) defines all operations for easy testing and dependency injection
- **Generic Operations**: Type-safe HTTP operations using Go generics in [api_client.go](api_client.go)
  - `Get[T]`, `GetSingle[T]`, `Post[T]`, `Delete[T]` for common HTTP operations
  - All use `*resty.Client` as first parameter for testability
- **Consistent CLI Structure**: Each resource has list/create/retrieve/modify/delete subcommands
  - Output formats: `--output pretty|json|jsonl` (default: pretty)
- **Pagination**: Standard `--limit` and `--after` flags across list commands
- **Authentication**: Uses `OPENAI_API_KEY` environment variable
- **Mocking Strategy**: Mock HTTP responses with `httpmock`, not OpenAI API behavior
- **Recent Features**: Certificates API with full mutual TLS lifecycle management

## Key Dependencies

- `github.com/go-resty/resty/v2` - HTTP client with built-in retry logic
- `github.com/urfave/cli/v3` - Modern CLI framework
- `github.com/jarcoal/httpmock` - HTTP mocking for tests
- `go.uber.org/mock` - Mock code generation
- `github.com/mark3labs/mcp-go` - Model Context Protocol support

## MCP Architecture

The MCP server provides AI assistants with tools and resources for managing OpenAI organizations:
- [cmd/mcp/main.go](cmd/mcp/main.go) - Server entry point
- [pkg/mcp/server.go](pkg/mcp/server.go) - Server setup and registration
- [pkg/mcp/tools.go](pkg/mcp/tools.go) - Tool implementations with reflection-based parameter validation
- [pkg/mcp/resources.go](pkg/mcp/resources.go) - Dynamic resource templates and subscriptions
- [pkg/mcp/auth.go](pkg/mcp/auth.go) - Authentication handling
- [pkg/mcp/uri.go](pkg/mcp/uri.go) - URI parsing and routing

MCP server exposes all OpenAI org management operations as tools, with automatic type validation via reflection.