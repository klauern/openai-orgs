# Active Context: [OpenAI Orgs]

> **Note:**
> This file is the primary source for actionable TODOs, engineering work, and current priorities. When reviewing the memory bank, always check the "Current Work Focus" and "Next Steps" sections here first for implementation tasks. Use supporting files (systemPatterns.md, progress.md, etc.) for architectural context, progress tracking, and rationale as needed.
>
> **All work should align with the documented project, workflow, structure, and coding style rules.** Always review these rules in the memory bank before planning or implementing changes to ensure consistency and quality.

## Current Work Focus

- Refactor GenericToolHandler to accept a client factory or ClientProvider interface for dependency injection
- Extract an OpenAIOrgsClient interface for all client methods used by tools (e.g., ListProjects, etc.)
- Update all tool handlers to use the interface, not the concrete client
- Update AddTools to allow injecting a mock client for tests
- Write integration-style tests for GenericToolHandler with a mock client
- Expand and complete test coverage, especially for error and integration scenarios (see Test Coverage section)
- Update documentation in doc.go to reflect the full set of tools and the generic handler/parameter schema pattern
- Continue documentation improvements per the established commenting pattern (see Documentation Progress section)

## Recent Changes

- [Change 1]
- [Change 2]
- [Change 3]

## Next Steps

- Implement dependency injection for GenericToolHandler
- Define and use OpenAIOrgsClient interface throughout tool handlers
- Refactor AddTools for mock client injection
- Develop integration tests using mocks
- Create and update test files: resources_test.go, tools_test.go, uri_test.go, auth_test.go
- Generate and maintain mocks using GoMock (see Test Implementation Plan)
- Update and maintain doc.go and related documentation files
- Continue CLI and usage/audit documentation improvements

## Test Coverage and Implementation Plan

### Test Coverage Status (from mcp-tests.txt)

- **Resources (resources.go):**
  - AddResources: registration, MIME types, URI validation
  - createResourceHandler: auth, param extraction, error/subscription handling
  - Individual handlers: handleActiveProjects, handleCurrentMembers, handleUsageDashboard (pagination, client, formatting, error cases)
  - SubscriptionManager: subscribe/unsubscribe, channel mgmt, concurrency, notification delivery
  - pollForChanges: context cancellation, update triggers, error handling
- **Tools (tools.go):**
  - AddTools: registration, description validation
  - handleListProjects: auth, client, formatting, error cases
- **URI Handling (uri.go):**
  - ParseURI: valid/invalid parsing, resource type extraction
- **Authentication (auth.go):**
  - Context handling: token extraction/validation, error cases

### Test Implementation Phases

- **Phase 1:** Core functionality (infra, resource tests, tool tests)
- **Phase 2:** Advanced features (subscription, polling)
- **Phase 3:** Integration (end-to-end, error recovery, performance)
- **Mocks:** Client and server interfaces using GoMock
- **Helpers:** setupTest, setupResourceTest for common test setup
- **Coverage Goals:** Line >80%, Branch >70%, Function >90%
- **Notes:** Table-driven tests, mock external dependencies, thorough error testing, performance tests for subscriptions, document test assumptions

## Tool Implementation and Framework (from mcp-tools-plan.txt)

- **All major tools implemented in pkg/mcp/tools.go using a generic handler and parameter schema framework**
- **Tool Categories:**
  - Project Management: list/create/retrieve/modify/archive_project
  - Project User Management: list/add/remove/retrieve/modify_project_user
  - Project API Keys: list/retrieve/delete_project_api_key
  - Project Service Accounts: list/create/retrieve/delete_project_service_account
  - User Management: list/retrieve/delete/modify_user_role
  - Invites: list/create/retrieve/delete_invite
  - Usage/Billing: get_usage
- **Framework:**
  - GenericToolHandler with ToolHandlerFunc and ParamSchema for validation, client instantiation, result formatting
  - Parameter schemas registered with mcp.NewTool using type helpers
  - Code structured for testability, with TODOs for further abstraction and easier mocking
- **Testing:**
  - Table-driven unit tests in tools_test.go
  - GoMock for interfaces (mock_interfaces.go)
  - Test plan in mcp-tests.txt
- **Next Steps:**
  - Expand test coverage, refactor for full dependency injection, update documentation, maintain parameter registration and schema

## Documentation Progress and Strategy (from commenting.txt)

- **Centralized package documentation in doc.go**
- **File-level documentation progress:**
  - Core API Client, Authentication & Users, Project Management: documented
  - Usage & Audit, CLI Commands: in progress
- **Guidelines:**
  1. Package docs in doc.go, not repeated elsewhere
  2. Type docs: what the type represents
  3. Field docs: complete sentences, concise
  4. Method docs: what it does, errors, examples
  5. Constants/vars: purpose, formatting, group comments
  6. Examples/tests: add Example functions, document assumptions
  7. General: active voice, complete sentences, focus on what, not how
- **Pattern established in project_api_keys.go:**
  - Type, field, function, and constant documentation
  - Context and usage information
- **Taskfile additions:**
  - Tasks to run go doc and check for missing comments
  - Task to start a local doc server
- **Progress:** 12/24 files documented

## Active Decisions and Considerations

- Prioritize testability and interface-driven design for all new and refactored code
- Maintain and expand documentation to ensure onboarding and maintenance are easy
- Use table-driven tests and GoMock for all external dependencies
- Focus on error handling and integration scenarios in upcoming test phases
- Ensure all new tools and features follow the established handler and parameter schema framework

## Important Patterns and Preferences

- Generic handler and parameter schema for all tools
- Centralized, example-driven documentation in doc.go
- Table-driven tests and interface-based mocks
- Documentation and test coverage tracked in workbooks and memory bank

## Learnings and Project Insights

- Decoupling tool handlers from concrete client implementations improves testability and flexibility
- Centralized documentation and clear commenting patterns accelerate onboarding and reduce errors
- Phased test implementation and coverage goals provide a clear roadmap for quality assurance
