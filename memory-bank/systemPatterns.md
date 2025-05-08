# System Patterns: [OpenAI Orgs]

## System Architecture

- Follows standard Go project structure: /cmd for CLI, root for core API client
- Clear separation of concerns between CLI, API client, and supporting packages
- Interfaces are used to enable dependency injection and testability

## Key Technical Decisions

- All API endpoints have corresponding CLI commands
- Use generics and strongly typed constants where appropriate
- Mock external APIs using jarcoal/httpmock for tests
- Use uber.go/mock for interface mocking
- Never use testify for test generation

## Design Patterns in Use

- Dependency injection via interfaces for easier testing
- Helper functions for test setup/teardown
- Standard Go error handling with context wrapping

## Component Relationships

- CLI commands interact with the API client via interfaces
- API client handles HTTP requests and responses
- Tests interact with mocks and helpers for isolation

## Critical Implementation Paths

- Adding a new API endpoint: implement in API client, expose via CLI, add tests and documentation
- Refactoring: update interfaces, mocks, and ensure all tests pass
