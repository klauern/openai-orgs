# openai-orgs - CLI for OpenAI Platform Management API

[![codecov](https://codecov.io/gh/klauern/openai-orgs/graph/badge.svg?token=7ENEG01SSU)](https://codecov.io/gh/klauern/openai-orgs)

`openai-orgs` is a command-line interface (CLI) tool for interacting with the OpenAI Platform Administration APIs. It provides various commands to manage projects, users, API keys, service accounts, invites, and more.

## Installation

To install `openai-orgs`, make sure you have Go installed on your system, then run:

```
go install github.com/klauern/openai-orgs/cmd/openai-orgs@latest
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

- `audit-logs`: Manage audit logs
- `invites`: Manage organization invites
- `users`: Manage organization users
- `projects`: Manage organization projects
- `project-users`: Manage project users
- `project-service-accounts`: Manage project service accounts
- `project-api-keys`: Manage project API keys

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

3. List project API keys:

```bash
openai-orgs project-api-keys list --project-id <project_id>
```

4. Create an invite:

```bash
openai-orgs invites create --email user@example.com --role member
```

## Default Settings

- The CLI uses the OpenAI API base URL: `https://api.openai.com/v1`
- Authentication is handled using the `OPENAI_API_KEY` environment variable
- List commands typically have optional `--limit` and `--after` flags to control pagination

## Error Handling

If an error occurs during command execution, the CLI will display an error message and exit with a non-zero status code.

## Contributing

Contributions to `openai-orgs` are welcome! Please submit issues and pull requests on the GitHub repository.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
