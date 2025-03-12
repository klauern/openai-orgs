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
- Use tabs for indentation, blanks for alignment
- Interfaces preferred to facilitate testing
- NEVER use testify for test generation
- Follow Go error handling pattern (check err != nil)
- Use proper documentation comments for exported functions
- CLI commands should be organized in subpackages
- All API endpoints should have corresponding CLI commands
- Use consistent naming: lowercase for parameters, camelCase for exported fields
- Core dependencies: resty/v2 for HTTP, urfave/cli/v3 for CLI

## Project Structure
- `/cmd` - CLI commands and main entrypoints
- Root package - Core API client implementation
- Tests should mock HTTP responses using jarcoal/httpmock
