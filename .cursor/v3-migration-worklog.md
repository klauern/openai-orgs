# urfave/cli v3 Migration Worklog

## Changes Required

- [x] Update import paths from v2 to v3 in all files
- [x] Change `Subcommands` field to `Commands` in command structures
- [x] Update utils.go with new context pattern
- [x] Update projects.go with new context pattern and fix Int() conversion
- [x] Update project_rate_limits.go
  - [x] listProjectRateLimits
  - [x] validateModifyProjectRateLimitContext
  - [x] modifyProjectRateLimit
- [x] Update project_api_keys.go
  - [x] listProjectApiKeys
  - [x] retrieveProjectApiKey
  - [x] deleteProjectApiKey
- [x] Update project_service_accounts.go
  - [x] listProjectServiceAccounts
  - [x] createProjectServiceAccount
  - [x] retrieveProjectServiceAccount
  - [x] deleteProjectServiceAccount
- [x] Update audit_logs.go
  - [x] listAuditLogs
- [x] Update users.go
  - [x] listUsers
  - [x] retrieveUser
  - [x] deleteUser
  - [x] modifyUserRole
- [x] Update org_invites.go
  - [x] listInvites
  - [x] createInvite
  - [x] deleteInvite
  - [x] retrieveInvite
- [x] Update project_users.go
  - [x] listProjectUsers
  - [x] createProjectUser
  - [x] retrieveProjectUser
  - [x] modifyProjectUser
  - [x] deleteProjectUser
- [x] Update main.go flag actions
  - [x] Update output flag action to use new context pattern
  - [x] Update api-key flag to use Sources instead of EnvVars
  - [x] Fix ValidOutputFormats access
  - [x] Fix App initialization (changed from App to Command)
  - [x] Fix Run function usage (updated to use app.Run with context)

## Migration Pattern

For each function that needs to be updated:

1. Change signature from `func name(c *cli.Context) error` to `func name(ctx context.Context, cmd *cli.Command) error`
2. Update all `c.String()` calls to `cmd.String()`
3. Update all `c.Int()` calls to `int(cmd.Int())` for int parameters, or `int64(cmd.Int())` for int64 parameters
4. Update all `c.Bool()` calls to `cmd.Bool()`
5. Update client creation to use `newClient(ctx, cmd)`
6. Update Action functions to use `func(ctx context.Context, cmd *cli.Command) error`
7. Use common flag definitions from utils.go where possible (e.g., projectIDFlag, limitFlag)
8. When updating API calls, make sure to use the correct parameter types and structures
9. For API calls that modify data but don't return the updated object, follow up with a retrieve call to show the changes
10. Pay attention to API response types and handle them correctly (e.g., []Type vs ListResponse[Type])
11. Be mindful of field name differences between APIs (e.g., CreatedAt vs AddedAt)
12. Flag actions in v3 have a different signature: `func(ctx context.Context, cmd *cli.Command, value T) error`
13. Use `Sources` instead of `EnvVars` for environment variable configuration in v3
14. In v3, the main app is a Command instead of an App, and needs to be run with app.Run(context.Background(), os.Args)

## Progress

### 2024-02-22

- ✅ Completed initial setup and import path updates
- ✅ Updated utils.go with new context pattern
- ✅ Updated projects.go as first example of full migration
- ✅ Updated project_rate_limits.go with new context pattern and fixed Int64 handling
- ✅ Updated project_api_keys.go with new context pattern
- ✅ Updated project_service_accounts.go with new context pattern and standardized flags
- ✅ Updated audit_logs.go with new context pattern and fixed API parameter handling
- ✅ Updated users.go with new context pattern and improved user feedback
- ✅ Updated org_invites.go with new context pattern and fixed response types
- ✅ Updated project_users.go with new context pattern and fixed field names
- ✅ Completed main.go updates
  - ✅ Updated flag actions to use new context pattern
  - ✅ Updated environment variable configuration
  - ✅ Fixed App initialization to use Command
  - ✅ Fixed Run function usage to use app.Run with context

## Notes

- The `Int()` method in v3 returns `int64`, needs explicit conversion to `int` where required
- All command handlers need both `context.Context` and `*cli.Command` parameters
- Flag actions in main.go need special attention as they use a different signature
- When dealing with int64 fields, use `int64(cmd.Int())` for the conversion
- Be careful with API field names and types, they should match the actual API
- Use common flag definitions from utils.go to maintain consistency
- Pay attention to API parameter structures and use them correctly (e.g., AuditLogListParams)
- For better UX, retrieve and show updated objects after modification operations
- Some APIs return direct slices while others use ListResponse wrapper - handle accordingly
- Field names may vary between APIs (e.g., CreatedAt vs AddedAt) - check API docs
- Flag actions in v3 take an additional value parameter of the flag's type
- Environment variables are configured using `Sources` in v3 instead of `EnvVars`
- Commands in v3 are run using app.Run(context.Background(), os.Args) instead of cli.Run
