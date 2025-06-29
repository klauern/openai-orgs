# Project Brief: [OpenAI Orgs]

## Overview

[High-level overview of what you're building]

## Core Requirements

- Follow Go error handling and code style guidelines
- All API endpoints must have corresponding CLI commands
- Use helper functions and interfaces for testability
- Mock external APIs using jarcoal/httpmock
- Use standard Go project structure and naming conventions

## Goals

- Maintain high code quality and test coverage
- Ensure clear separation of concerns in code organization
- Provide a robust CLI and API client
- Enable easy testing and maintainability

## Project Scope

- In scope: API client, CLI commands, test suite, documentation, code style enforcement
- Out of scope: Non-Go implementations, unsupported API endpoints

## Development Workflow

- Build: `task build` or `go build -v ./...`
- Install: `task install`
- Lint: `task lint` or `golangci-lint run`
- Format: `task fmt` or `gofmt -s -w -e .`
- Test all: `task test` or `go test -v -coverprofile=coverage.out -timeout=120s -parallel=10 ./...`
- Test single file: `go test -v ./path/to/file_test.go`
- Test specific test: `go test -v ./... -run "TestName"`
- Test coverage: `task cover`

- Always run linter and tests before committing
- Use helper functions for test setup/teardown
- Mock HTTP responses using jarcoal/httpmock
- Never use testify for test generation
