# OpenAI Orgs Project Guide

## Common Commands

- Build: `task build` or `go build -v ./...`
- Install: `task install`
- Lint: `task lint` or `golangci-lint run`
- Format: `task fmt` or `gofmt -s -w -e .`
- Test all: `task test` or `go test -v -coverprofile=coverage.out -timeout=120s -parallel=10 ./...`
- Test single file: `go test -v ./path/to/file_test.go`
- Test specific test: `go test -v ./... -run "TestName"`
- Test coverage: `task cover`

## Code Style Guidelines

- **Testing**: Use helper functions for test setup/teardown, mock external APIs
  - NEVER use testify for test generation
- Follow Go error handling pattern (check err != nil)
- **Comments**: Document exported functions and types with meaningful comments
- CLI commands should be organized in subpackages
- All API endpoints should have corresponding CLI commands
- Use consistent naming: lowercase for parameters, camelCase for exported fields
- **Imports**: Standard Go import organization (stdlib first, then external)
  - Core dependencies: resty/v2 for HTTP, urfave/cli/v3 for CLI
- **Error handling**: Use `fmt.Errorf` with context wrapping, e.g., `fmt.Errorf("error making request: %v", err)`
- **Naming**:
  - Types: PascalCase (e.g., `Client`, `ListResponse`)
  - Constants: Use prefix conventions (e.g., `OwnerTypeUser`)
  - Functions: PascalCase for exported, camelCase for internal
- **Types**: Use generics for common operations, strongly type constants with custom types

## Project Structure

- `/cmd` - CLI commands and main entrypoints
- Root package - Core API client implementation
- Tests should mock HTTP responses using jarcoal/httpmock
