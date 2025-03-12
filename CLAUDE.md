# OpenAI-Orgs Project Guide

## Build, Lint, and Test Commands
- Build: `task build` or `go build -v ./...` 
- Install: `task install`
- Lint: `task lint` or `golangci-lint run`
- Format: `task fmt` or `gofmt -s -w -e .`
- Test all: `task test` or `go test -v -coverprofile=coverage.out -timeout=120s -parallel=10 ./...`
- Test single file: `go test -v -coverprofile=coverage.out ./[filename]_test.go`
- Test specific test: `go test -v -run TestFunctionName ./...`
- Coverage report: `task cover` or `go tool cover -html=coverage.out`

## Code Style Guidelines
- **Imports**: Standard Go import organization (stdlib first, then external)
- **Error handling**: Use `fmt.Errorf` with context wrapping, e.g., `fmt.Errorf("error making request: %v", err)`
- **Naming**: 
  - Types: PascalCase (e.g., `Client`, `ListResponse`)
  - Constants: Use prefix conventions (e.g., `OwnerTypeUser`)
  - Functions: PascalCase for exported, camelCase for internal
- **Types**: Use generics for common operations, strongly type constants with custom types
- **Testing**: Use helper functions for test setup/teardown, mock external APIs
- **Comments**: Document exported functions and types with meaningful comments