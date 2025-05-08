# Tech Context: [OpenAI Orgs]

## Technologies Used

- Go (main language)
- resty/v2 (HTTP client)
- urfave/cli/v3 (CLI framework)
- jarcoal/httpmock (HTTP mocking for tests)
- uber.go/mock (interface mocking)

## Development Setup

- Use Go modules for dependency management
- Common commands are managed via Taskfile.yml
- Standard Go tools for build, test, lint, and format

## Technical Constraints

- Only Go is supported for implementation
- All code must pass lint and formatting checks
- No use of testify for testing

## Dependencies

- resty/v2 for HTTP
- urfave/cli/v3 for CLI
- jarcoal/httpmock for HTTP mocking
- uber.go/mock for interface mocking

## Tool Usage Patterns

- Use helper functions for test setup/teardown
- Mock external APIs in tests
- Use interfaces for dependency injection and testability
