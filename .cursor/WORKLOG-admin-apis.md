# Administration APIs Implementation Worklog

## Overview

This worklog tracks the implementation of OpenAI Administration APIs following the established patterns in:

- api_client.go (HTTP client and generic request handling)
- administration.go (common types and utilities)
- projects.go (CLI command structure)

## Implementation Pattern

Each API endpoint will follow this structure:

1. Define types in a dedicated file (e.g., `users.go`, `invites.go`)
2. Implement client methods using the generic `Get`, `Post`, `Delete` functions
3. Create CLI commands in the `cmd/` directory
4. Add tests for both the client methods and CLI commands

## API Endpoints to Implement

### 1. Organization Management

- [ ] List organizations
- [ ] Retrieve organization
- [ ] Update organization settings
- [ ] Delete organization

Files needed:

```
organizations.go
cmd/organizations.go
organizations_test.go
```

### 2. User Management

- [x] List users
- [x] Get user details
- [x] Update user role
- [x] Remove user

Files needed:

```
users.go (existing)
cmd/users.go (existing)
users_test.go (existing)
```

### 3. Project Management

- [x] List projects
- [x] Create project
- [x] Get project details
- [x] Update project
- [x] Archive project

Files needed:

```
projects.go (existing)
cmd/projects.go (existing)
projects_test.go (existing)
```

### 4. Project API Keys

- [x] List project API keys
- [x] Create API key
- [x] Get API key details
- [x] Delete API key

Files needed:

```
project_api_keys.go (existing)
cmd/project_api_keys.go (existing)
project_api_keys_test.go (existing)
```

### 5. Project Service Accounts

- [x] List service accounts
- [x] Create service account
- [x] Get service account details
- [x] Delete service account

Files needed:

```
project_service_accounts.go (existing)
cmd/project_service_accounts.go (existing)
project_service_accounts_test.go (existing)
```

### 6. Project Rate Limits

- [x] List rate limits
- [x] Update rate limits

Files needed:

```
project_rate_limits.go (existing)
cmd/project_rate_limits.go (existing)
project_rate_limits_test.go (existing)
```

### 7. Invites

- [x] List invites
- [x] Create invite
- [x] Get invite details
- [x] Delete invite

Files needed:

```
invites.go (existing)
cmd/invites.go (existing)
invites_test.go (existing)
```

### 8. Audit Logs

- [x] List audit logs
- [x] Get audit log details

Files needed:

```
audit_logs.go (existing)
cmd/audit_logs.go (existing)
audit_logs_test.go (existing)
```

### 9. Admin API Keys

- [ ] List admin API keys
- [ ] Create admin API key
- [ ] Get admin API key details
- [ ] Delete admin API key

Files needed:

```
admin_api_keys.go
cmd/admin_api_keys.go
admin_api_keys_test.go
```

## Implementation Tasks

### Organization Management Implementation

1. [ ] Create `organizations.go`:
   - Define Organization struct
   - Implement CRUD methods
   - Add JSON struct tags

2. [ ] Create `cmd/organizations.go`:
   - Define CLI commands (list, get, update, delete)
   - Implement command handlers
   - Add flags and documentation

3. [ ] Create `organizations_test.go`:
   - Add unit tests for client methods
   - Add integration tests
   - Test error cases

### Admin API Keys Implementation

1. [x] Create `admin_api_keys.go`:
   - Define AdminAPIKey struct
   - Implement CRUD methods
   - Add JSON struct tags

2. [x] Create `cmd/admin_api_keys.go`:
   - Define CLI commands (list, create, get, delete)
   - Implement command handlers
   - Add flags and documentation

3. [x] Create `admin_api_keys_test.go`:
   - Add unit tests for client methods
   - Add integration tests
   - Test error cases

4. [x] Add to `cmd/openai-orgs/main.go`:
   - Register admin API keys command in CLI application

CLI Commands:

```
# List all admin API keys
openai-orgs admin-api-keys list [--limit N] [--after KEY_ID]

# Create a new admin API key
openai-orgs admin-api-keys create --name "Key Name" --scopes organization.read,organization.write

# Get details of a specific admin API key
openai-orgs admin-api-keys retrieve --id key_123

# Delete an admin API key
openai-orgs admin-api-keys delete --id key_123
```

### Testing Strategy

- Use existing test helper in `test_helpers.go`
- Mock HTTP responses using `httpmock`
- Test both success and error cases
- Ensure proper error handling and messages

### Documentation

- [ ] Update README.md with new commands
- [ ] Add examples for organization management and admin API keys
- [ ] Document any new configuration options

## Notes

- Follow existing patterns for error handling
- Use consistent naming conventions
- Maintain backward compatibility
- Keep test coverage high (aim for >70%)
- Use shared utilities from `cmd/utils.go`

## Dependencies

- github.com/go-resty/resty/v2
- github.com/urfave/cli/v3
- github.com/jarcoal/httpmock (for testing)

## Progress Tracking

- [x] Initial worklog creation
- [ ] Organization management implementation
- [ ] Admin API keys implementation
- [x] User management implementation
- [x] Project management implementation
- [x] Project API keys implementation
- [x] Project service accounts implementation
- [x] Project rate limits implementation
- [x] Invites implementation
- [x] Audit logs implementation
