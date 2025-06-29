# openai-orgs - CLI for OpenAI Platform Management API

[![codecov](https://codecov.io/gh/klauern/openai-orgs/graph/badge.svg?token=7ENEG01SSU)](https://codecov.io/gh/klauern/openai-orgs)

`openai-orgs` is a comprehensive command-line interface (CLI) tool and API client for managing OpenAI Platform Administration APIs. It provides comprehensive management capabilities for OpenAI organizations including projects, users, API keys, service accounts, certificates, audit logs, and more. The project also includes a Model Context Protocol (MCP) server for AI assistant integration.

## Installation

### CLI Tool

To install `openai-orgs`, make sure you have Go installed on your system, then run:

```bash
go install github.com/klauern/openai-orgs/cmd/openai-orgs@latest
```

### MCP Server

To install the MCP server:

```bash
go install github.com/klauern/openai-orgs/cmd/mcp@latest
```

### Development Setup

For development, you can use the included Task runner:

```bash
# Install Task if you don't have it
go install github.com/go-task/task/v3/cmd/task@latest

# Build the CLI
task build

# Build the MCP server
task mcp:build

# Run tests
task test
```

## Configuration

Before using `openai-orgs`, you need to set up your OpenAI API key:

1. Log in to your OpenAI account at <https://platform.openai.com/>
2. Navigate to the API keys section
3. Create a new API key
4. Set the API key as an environment variable:

```bash
export OPENAI_API_KEY=your_api_key_here
```

## Usage

`openai-orgs` uses subcommands to organize its functionality. Here are the main commands:

### Organization Level Commands
- `audit-logs`: Manage audit logs
- `invites`: Manage organization invites
- `users`: Manage organization users
- `admin-api-keys`: Manage organization admin API keys
- `certificates`: Manage organization certificates (mutual TLS)

### Project Level Commands
- `projects`: Manage organization projects
- `project-users`: Manage project users
- `project-service-accounts`: Manage project service accounts
- `project-api-keys`: Manage project API keys
- `project-rate-limits`: Manage project rate limits
- `project-certificates`: Manage project certificates

### Output Formats

All commands support multiple output formats via the `--output` flag:
- `pretty` (default): Human-readable formatted output
- `json`: JSON format
- `jsonl`: JSON Lines format

To see available subcommands and options for each command, use the `--help` flag:

```bash
openai-orgs --help
openai-orgs <command> --help
```

### Examples

1. List all users in the organization:

```bash
openai-orgs users list
```

2. Create a new project:

```bash
openai-orgs projects create --name "My New Project"
```

3. List project API keys with JSON output:

```bash
openai-orgs project-api-keys list --project-id <project_id> --output json
```

4. Create an invite:

```bash
openai-orgs invites create --email user@example.com --role member
```

5. Manage certificates for mutual TLS:

```bash
openai-orgs certificates list
openai-orgs certificates create --cert-file ./cert.pem
```

6. View audit logs:

```bash
openai-orgs audit-logs list --limit 10
```

## MCP Server

The project includes a Model Context Protocol (MCP) server that provides AI assistants with tools and resources for managing OpenAI organizations.

### Running the MCP Server

```bash
# Run the MCP server
mcp

# For development with debugging
task mcp:dev
```

### MCP Features

- **Tools**: Complete set of tools for all OpenAI organization and project management operations
- **Resources**: Dynamic resource templates and subscription management
- **Type Safety**: Reflection-based parameter validation
- **Integration**: Works with Claude Desktop and other MCP-compatible AI assistants

### MCP Configuration

Add to your MCP client configuration (e.g., Claude Desktop):

```json
{
  "mcpServers": {
    "openai-orgs": {
      "command": "mcp",
      "env": {
        "OPENAI_API_KEY": "your_api_key_here"
      }
    }
  }
}
```

## Default Settings

- The CLI uses the OpenAI API base URL: `https://api.openai.com/v1`
- Authentication is handled using the `OPENAI_API_KEY` environment variable
- List commands typically have optional `--limit` and `--after` flags to control pagination
- Conservative retry strategy: 20 retries, 5-second wait, max 5-minute backoff

## Error Handling

If an error occurs during command execution, the CLI will display an error message and exit with a non-zero status code.

## Contributing

Contributions to `openai-orgs` are welcome! Please submit issues and pull requests on the GitHub repository.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
