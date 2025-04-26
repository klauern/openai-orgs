/*
Package openaiorgs provides a Go client for interacting with the OpenAI Organizations API.

The client handles authentication, rate limiting, and provides type-safe methods for all API operations.
It supports managing organization resources such as users, projects, API keys, and service accounts.

Basic Usage:

	client := openaiorgs.NewClient("", "your-api-key")
	projects, err := client.ListProjects(10, "", false)
	if err != nil {
		log.Fatal(err)
	}

The package is organized into several main components:

Core Client:
  - API client with built-in retries and rate limiting
  - Generic request/response handling
  - Pagination support for list operations

Authentication & Users:
  - Organization API key management
  - User management (invite, modify roles, remove)
  - Organization invitations

Project Management:
  - Project creation and configuration
  - Project user management
  - Project-specific API keys
  - Rate limit configuration
  - Service account management

Usage & Audit:
  - Usage tracking and reporting
  - Audit logging
  - Administrative operations

Each component provides a set of methods for interacting with the corresponding API endpoints.
All operations use strong typing and follow consistent patterns for error handling and response processing.

For detailed examples and documentation of specific types and methods, see the relevant type
and function documentation.
*/
package openaiorgs
