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

The package provides tools for direct organization management operations:

	list_projects - Lists all projects for the authenticated user

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

For more detailed information about specific components, refer to the individual
type and function documentation.
*/
package mcp
