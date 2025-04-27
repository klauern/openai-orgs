/*
Package mcp implements the Model Context Protocol (MCP) server for OpenAI organization management.
It provides a set of tools and resources for managing OpenAI projects, members, and usage statistics
through a standardized protocol interface.

# Server Configuration

The package provides a main entry point through NewMCPServer() which configures and returns an MCP
server with the following capabilities:

  - Tool execution for organization management
  - Static and dynamic resource access
  - Real-time resource updates through subscriptions
  - Authentication token management
  - Logging and instrumentation

# Resources

The package exposes several key resources through URI-based endpoints:

	openai-orgs://active-projects    - Lists currently active projects
	openai-orgs://current-members    - Shows current organization members
	openai-orgs://usage-dashboard    - Displays usage statistics and metrics

Each resource supports pagination and optional real-time updates through subscriptions.
Resource data is returned in specialized MIME types for proper content handling:

	application/vnd.openai-orgs.project-list+json
	application/vnd.openai-orgs.member-list+json
	application/vnd.openai-orgs.usage+json

# Tools

The package provides a comprehensive set of tools for direct organization management operations, including:

- Project management: list_projects, create_project, retrieve_project, modify_project, archive_project
- Project user management: list_project_users, add_project_user, remove_project_user, retrieve_project_user, modify_project_user
- Project API key management: list_project_api_keys, retrieve_project_api_key, delete_project_api_key
- Project service account management: list_project_service_accounts, create_project_service_account, retrieve_project_service_account, delete_project_service_account
- User management: list_users, retrieve_user, delete_user, modify_user_role
- Invite management: list_invites, create_invite, retrieve_invite, delete_invite
- Usage and billing statistics: get_usage

All tools are implemented using a generic handler and parameter schema pattern, ensuring consistent parameter validation, error handling, and testability. Parameters are registered with mcp.NewTool using type helpers (e.g., mcp.WithString, mcp.WithNumber, mcp.WithBoolean), making them visible and enforced in the MCP Inspector and compatible clients.

# Tool Implementation Framework

- Tools use a GenericToolHandler that takes a ToolHandlerFunc and a ParamSchema, handling parameter extraction, validation, client instantiation, and result formatting.
- Each tool defines its parameters using a ParamSchema and registers them with mcp.NewTool.
- The framework is designed for testability, with support for dependency injection and GoMock-based mocks.

# Resource Updates

Resources support real-time updates through a subscription mechanism. When a client
subscribes to a resource, they receive updates when the underlying data changes:

 1. Active project status changes
 2. Member list modifications
 3. Usage statistics updates

The update frequency is managed by an internal polling mechanism that efficiently
checks for changes in the underlying data.

# Authentication

All operations require proper authentication through an OpenAI API token. The token
is managed through context and is required for all tool and resource operations.

# Example Usage

To create and start an MCP server:

	server := mcp.NewMCPServer()
	server.ServeStdio()

Resource subscriptions can be enabled through the subscription parameter:

	{
		"uri": "openai-orgs://active-projects",
		"arguments": {
			"subscribe": true,
			"pagination": {
				"limit": 20,
				"after": "project_id"
			}
		}
	}

# Testing

- Unit tests for tool handlers use the standard Go testing package and GoMock for interface mocking.
- The framework supports dependency injection for easier testability.
- See llm-workbooks/mcp-tests.txt for the current test plan and coverage goals.

For more detailed information about specific components, refer to the individual
type and function documentation.
*/
package mcp
