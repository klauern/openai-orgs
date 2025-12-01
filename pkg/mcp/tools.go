package mcp

import (
	"context"
	"fmt"
	"os"
	"reflect"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// TODO: Refactor GenericToolHandler to accept a client factory or ClientProvider interface for dependency injection
// TODO: Extract an OpenAIOrgsClient interface for all client methods used by tools (e.g., ListProjects, etc.)
// TODO: Update all tool handlers to use the interface, not the concrete client
// TODO: Update AddTools to allow injecting a mock client for tests
// TODO: Write integration-style tests for GenericToolHandler with a mock client
// These changes will enable full-stack testing with mocks and easier future maintenance.

// ParamField defines a single parameter for a tool
// Name: parameter name
// Required: whether the parameter is mandatory
// Type: expected reflect.Kind (e.g., reflect.String)
// Description: human-readable description
// Enum: optional set of allowed values
type ParamField struct {
	Name        string
	Required    bool
	Type        reflect.Kind
	Description string
	Enum        []any
}

// ParamSchema defines the schema for tool parameters
// Fields: list of parameter fields
type ParamSchema struct {
	Fields []ParamField
}

// ToMCPParameterSchema converts ParamSchema to a JSON Schema-like map for MCP
func (ps *ParamSchema) ToMCPParameterSchema() map[string]any {
	props := map[string]any{}
	var required []string
	for _, f := range ps.Fields {
		fieldType := "string"
		switch f.Type {
		case reflect.String:
			fieldType = "string"
		case reflect.Float64, reflect.Int, reflect.Int64:
			fieldType = "number"
		case reflect.Bool:
			fieldType = "boolean"
		}
		prop := map[string]any{
			"type":        fieldType,
			"description": f.Description,
		}
		if len(f.Enum) > 0 {
			prop["enum"] = f.Enum
		}
		props[f.Name] = prop
		if f.Required {
			required = append(required, f.Name)
		}
	}
	schema := map[string]any{
		"type":       "object",
		"properties": props,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

// ExtractAndValidate extracts and validates parameters from a CallToolRequest
// Returns a map of validated parameters or an error if validation fails
func (ps *ParamSchema) ExtractAndValidate(req mcp.CallToolRequest) (map[string]any, error) {
	params := make(map[string]any)
	args, ok := req.Params.Arguments.(map[string]any)
	if !ok {
		// If Arguments is nil or not a map, treat as empty
		args = make(map[string]any)
	}
	for _, field := range ps.Fields {
		val, ok := args[field.Name]
		if !ok {
			if field.Required {
				return nil, fmt.Errorf("missing required parameter: %s", field.Name)
			}
			continue
		}
		if field.Type == reflect.Bool {
			switch v := val.(type) {
			case bool:
				params[field.Name] = v
			case string:
				switch v {
				case "true", "1":
					params[field.Name] = true
				case "false", "0":
					params[field.Name] = false
				default:
					return nil, fmt.Errorf("parameter '%s' must be a boolean (true/false/1/0)", field.Name)
				}
				continue // skip the type check below, already handled
			default:
				return nil, fmt.Errorf("parameter '%s' must be a boolean", field.Name)
			}
			continue
		}
		if field.Type != reflect.Invalid && reflect.TypeOf(val).Kind() != field.Type {
			return nil, fmt.Errorf("parameter '%s' must be of type %s", field.Name, field.Type.String())
		}
		if len(field.Enum) > 0 {
			found := false
			for _, allowed := range field.Enum {
				if val == allowed {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("parameter '%s' must be one of %v", field.Name, field.Enum)
			}
		}
		params[field.Name] = val
	}
	return params, nil
}

// ToolHandlerFunc is a generic handler signature for MCP tools
// ctx: context, client: OpenAI orgs client, params: validated parameters
// Returns: result (any) and error
type ToolHandlerFunc func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error)

// GenericToolHandler wraps a ToolHandlerFunc for MCP
// Handles parameter extraction/validation, client instantiation, error handling, and result formatting
func GenericToolHandler(handler ToolHandlerFunc, paramSchema ParamSchema) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params, err := paramSchema.ExtractAndValidate(req)
		if err != nil {
			return nil, err
		}
		authToken := ctx.Value(authToken{}).(string)
		client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
		result, err := handler(ctx, client, params)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(fmt.Sprintf("%v", result)), nil
	}
}

func AddTools(s *server.MCPServer) {
	logTool := func(name string, schema ParamSchema) {
		fmt.Fprintf(os.Stderr, "Registering tool: %s\n", name)
		for _, field := range schema.Fields {
			fmt.Fprintf(os.Stderr, "  Param: %s (required: %v, type: %v, desc: %s)\n", field.Name, field.Required, field.Type, field.Description)
		}
	}

	// Project Management
	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of projects to return (default 100)"},
				{Name: "after", Required: false, Type: reflect.String, Description: "Project ID to start after (for pagination)"},
				{Name: "activeOnly", Required: false, Type: reflect.Bool, Description: "If true, only return active projects"},
			},
		}
		s.AddTool(mcp.NewTool(
			"list_projects",
			mcp.WithDescription("Lists all projects for the authenticated user"),
			mcp.WithNumber("limit", mcp.Description("Maximum number of projects to return (default 100)")),
			mcp.WithString("after", mcp.Description("Project ID to start after (for pagination)")),
			mcp.WithBoolean("activeOnly", mcp.Description("If true, only return active projects")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					limit := 100
					if v, ok := params["limit"]; ok {
						limit = int(v.(float64)) // MCP sends numbers as float64
					}
					after := ""
					if v, ok := params["after"]; ok {
						after = v.(string)
					}
					activeOnly := false
					if v, ok := params["activeOnly"]; ok {
						activeOnly = v.(bool)
					}
					projects, err := client.ListProjects(limit, after, activeOnly)
					if err != nil {
						return nil, err
					}
					return projects.String(), nil
				},
				schema,
			),
		)
		logTool("list_projects", schema)
	}

	// --- Project Management ---
	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "name", Required: true, Type: reflect.String, Description: "Project name"},
			},
		}
		s.AddTool(mcp.NewTool(
			"create_project",
			mcp.WithDescription("Creates a new project in the organization"),
			mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					name := params["name"].(string)
					project, err := client.CreateProject(name)
					if err != nil {
						return nil, fmt.Errorf("failed to create project: %w", err)
					}
					return project.String(), nil
				},
				schema,
			),
		)
		logTool("create_project", schema)
	}

	s.AddTool(mcp.NewTool(
		"retrieve_project",
		mcp.WithDescription("Retrieves details of a specific project by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Project ID")),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				id := params["id"].(string)
				project, err := client.RetrieveProject(id)
				if err != nil {
					return nil, fmt.Errorf("failed to retrieve project: %w", err)
				}
				return project.String(), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "id", Required: true, Type: reflect.String, Description: "Project ID"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"modify_project",
		mcp.WithDescription("Modifies the name of an existing project"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Project ID")),
		mcp.WithString("name", mcp.Required(), mcp.Description("New project name")),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				id := params["id"].(string)
				name := params["name"].(string)
				project, err := client.ModifyProject(id, name)
				if err != nil {
					return nil, fmt.Errorf("failed to modify project: %w", err)
				}
				return project.String(), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "id", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "name", Required: true, Type: reflect.String, Description: "New project name"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"archive_project",
		mcp.WithDescription("Archives a project by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Project ID")),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				id := params["id"].(string)
				project, err := client.ArchiveProject(id)
				if err != nil {
					return nil, fmt.Errorf("failed to archive project: %w", err)
				}
				return project.String(), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "id", Required: true, Type: reflect.String, Description: "Project ID"},
				},
			},
		),
	)

	// --- Project User Management ---
	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
				{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of users to return"},
				{Name: "after", Required: false, Type: reflect.String, Description: "User ID to start after (for pagination)"},
			},
		}
		s.AddTool(mcp.NewTool(
			"list_project_users",
			mcp.WithDescription("Lists all users for a given project"),
			mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithNumber("limit", mcp.Description("Maximum number of users to return")),
			mcp.WithString("after", mcp.Description("User ID to start after (for pagination)")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					projectID := params["projectId"].(string)
					limit := 0
					if v, ok := params["limit"]; ok {
						limit = int(v.(float64))
					}
					after := ""
					if v, ok := params["after"]; ok {
						after = v.(string)
					}
					users, err := client.ListProjectUsers(projectID, limit, after)
					if err != nil {
						return nil, fmt.Errorf("failed to list project users: %w", err)
					}
					return users.String(), nil
				},
				schema,
			),
		)
	}

	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
				{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
				{Name: "role", Required: true, Type: reflect.String, Description: "Role (e.g., owner, member)"},
			},
		}
		s.AddTool(mcp.NewTool(
			"add_project_user",
			mcp.WithDescription("Adds a user to a project with a specific role"),
			mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("userId", mcp.Required(), mcp.Description("User ID")),
			mcp.WithString("role", mcp.Required(), mcp.Description("Role (e.g., owner, member)")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					projectID := params["projectId"].(string)
					userID := params["userId"].(string)
					role := params["role"].(string)
					user, err := client.CreateProjectUser(projectID, userID, role)
					if err != nil {
						return nil, fmt.Errorf("failed to add project user: %w", err)
					}
					return user.String(), nil
				},
				schema,
			),
		)
	}

	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
				{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
			},
		}
		s.AddTool(mcp.NewTool(
			"remove_project_user",
			mcp.WithDescription("Removes a user from a project"),
			mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("userId", mcp.Required(), mcp.Description("User ID")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					projectID := params["projectId"].(string)
					userID := params["userId"].(string)
					err := client.DeleteProjectUser(projectID, userID)
					if err != nil {
						return nil, fmt.Errorf("failed to remove project user: %w", err)
					}
					return fmt.Sprintf("User %s removed from project %s", userID, projectID), nil
				},
				schema,
			),
		)
	}

	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
				{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
			},
		}
		s.AddTool(mcp.NewTool(
			"retrieve_project_user",
			mcp.WithDescription("Retrieves a specific user from a project"),
			mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("userId", mcp.Required(), mcp.Description("User ID")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					projectID := params["projectId"].(string)
					userID := params["userId"].(string)
					user, err := client.RetrieveProjectUser(projectID, userID)
					if err != nil {
						return nil, fmt.Errorf("failed to retrieve project user: %w", err)
					}
					return user.String(), nil
				},
				schema,
			),
		)
	}

	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
				{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
				{Name: "role", Required: true, Type: reflect.String, Description: "New role (e.g., owner, member)"},
			},
		}
		s.AddTool(mcp.NewTool(
			"modify_project_user",
			mcp.WithDescription("Modifies a user's role in a project"),
			mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("userId", mcp.Required(), mcp.Description("User ID")),
			mcp.WithString("role", mcp.Required(), mcp.Description("New role (e.g., owner, member)")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					projectID := params["projectId"].(string)
					userID := params["userId"].(string)
					role := params["role"].(string)
					user, err := client.ModifyProjectUser(projectID, userID, role)
					if err != nil {
						return nil, fmt.Errorf("failed to modify project user: %w", err)
					}
					return user.String(), nil
				},
				schema,
			),
		)
	}

	// --- Project API Keys ---
	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
				{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of API keys to return"},
				{Name: "after", Required: false, Type: reflect.String, Description: "API key ID to start after (for pagination)"},
			},
		}
		s.AddTool(mcp.NewTool(
			"list_project_api_keys",
			mcp.WithDescription("Lists all API keys for a given project"),
			mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithNumber("limit", mcp.Description("Maximum number of API keys to return")),
			mcp.WithString("after", mcp.Description("API key ID to start after (for pagination)")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					projectID := params["projectId"].(string)
					limit := 0
					if v, ok := params["limit"]; ok {
						limit = int(v.(float64))
					}
					after := ""
					if v, ok := params["after"]; ok {
						after = v.(string)
					}
					keys, err := client.ListProjectApiKeys(projectID, limit, after)
					if err != nil {
						return nil, fmt.Errorf("failed to list project API keys: %w", err)
					}
					return keys.String(), nil
				},
				schema,
			),
		)
	}

	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
				{Name: "apiKeyId", Required: true, Type: reflect.String, Description: "API Key ID"},
			},
		}
		s.AddTool(mcp.NewTool(
			"delete_project_api_key",
			mcp.WithDescription("Deletes a specific API key from a project"),
			mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("apiKeyId", mcp.Required(), mcp.Description("API Key ID")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					projectID := params["projectId"].(string)
					apiKeyID := params["apiKeyId"].(string)
					err := client.DeleteProjectApiKey(projectID, apiKeyID)
					if err != nil {
						return nil, fmt.Errorf("failed to delete project API key: %w", err)
					}
					return fmt.Sprintf("API key %s deleted from project %s", apiKeyID, projectID), nil
				},
				schema,
			),
		)
	}

	// --- Project Service Accounts ---
	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
				{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of service accounts to return"},
				{Name: "after", Required: false, Type: reflect.String, Description: "Service account ID to start after (for pagination)"},
			},
		}
		s.AddTool(mcp.NewTool(
			"list_project_service_accounts",
			mcp.WithDescription("Lists all service accounts for a given project"),
			mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithNumber("limit", mcp.Description("Maximum number of service accounts to return")),
			mcp.WithString("after", mcp.Description("Service account ID to start after (for pagination)")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					projectID := params["projectId"].(string)
					limit := 0
					if v, ok := params["limit"]; ok {
						limit = int(v.(float64))
					}
					after := ""
					if v, ok := params["after"]; ok {
						after = v.(string)
					}
					accounts, err := client.ListProjectServiceAccounts(projectID, limit, after)
					if err != nil {
						return nil, fmt.Errorf("failed to list project service accounts: %w", err)
					}
					return accounts.String(), nil
				},
				schema,
			),
		)
	}

	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
				{Name: "name", Required: true, Type: reflect.String, Description: "Service account name"},
			},
		}
		s.AddTool(mcp.NewTool(
			"create_project_service_account",
			mcp.WithDescription("Creates a new service account for a project"),
			mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Service account name")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					projectID := params["projectId"].(string)
					name := params["name"].(string)
					account, err := client.CreateProjectServiceAccount(projectID, name)
					if err != nil {
						return nil, fmt.Errorf("failed to create project service account: %w", err)
					}
					return account.String(), nil
				},
				schema,
			),
		)
	}

	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
				{Name: "serviceAccountId", Required: true, Type: reflect.String, Description: "Service Account ID"},
			},
		}
		s.AddTool(mcp.NewTool(
			"delete_project_service_account",
			mcp.WithDescription("Deletes a service account from a project"),
			mcp.WithString("projectId", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("serviceAccountId", mcp.Required(), mcp.Description("Service Account ID")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					projectID := params["projectId"].(string)
					serviceAccountID := params["serviceAccountId"].(string)
					err := client.DeleteProjectServiceAccount(projectID, serviceAccountID)
					if err != nil {
						return nil, fmt.Errorf("failed to delete project service account: %w", err)
					}
					return fmt.Sprintf("Service account %s deleted from project %s", serviceAccountID, projectID), nil
				},
				schema,
			),
		)
	}

	// --- User Management ---
	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of users to return"},
				{Name: "after", Required: false, Type: reflect.String, Description: "User ID to start after (for pagination)"},
			},
		}
		s.AddTool(mcp.NewTool(
			"list_users",
			mcp.WithDescription("Lists all users in the organization"),
			mcp.WithNumber("limit", mcp.Description("Maximum number of users to return")),
			mcp.WithString("after", mcp.Description("User ID to start after (for pagination)")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					limit := 0
					if v, ok := params["limit"]; ok {
						limit = int(v.(float64))
					}
					after := ""
					if v, ok := params["after"]; ok {
						after = v.(string)
					}
					users, err := client.ListUsers(limit, after)
					if err != nil {
						return nil, fmt.Errorf("failed to list users: %w", err)
					}
					return users.String(), nil
				},
				schema,
			),
		)
	}

	s.AddTool(mcp.NewTool(
		"retrieve_user",
		mcp.WithDescription("Retrieves details of a specific user by ID"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				userID := params["userId"].(string)
				user, err := client.RetrieveUser(userID)
				if err != nil {
					return nil, fmt.Errorf("failed to retrieve user: %w", err)
				}
				return user.String(), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"delete_user",
		mcp.WithDescription("Deletes a user from the organization by ID"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				userID := params["userId"].(string)
				err := client.DeleteUser(userID)
				if err != nil {
					return nil, fmt.Errorf("failed to delete user: %w", err)
				}
				return fmt.Sprintf("User %s deleted", userID), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"modify_user_role",
		mcp.WithDescription("Modifies a user's role in the organization"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				userID := params["userId"].(string)
				role := params["role"].(string)
				err := client.ModifyUserRole(userID, role)
				if err != nil {
					return nil, fmt.Errorf("failed to modify user role: %w", err)
				}
				return fmt.Sprintf("User %s role updated to %s", userID, role), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
					{Name: "role", Required: true, Type: reflect.String, Description: "New role (e.g., owner, member)"},
				},
			},
		),
	)

	// --- Invites ---
	{
		schema := ParamSchema{
			Fields: []ParamField{},
		}
		s.AddTool(mcp.NewTool(
			"list_invites",
			mcp.WithDescription("Lists all pending invites in the organization"),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					invites, err := client.ListInvites()
					if err != nil {
						return nil, fmt.Errorf("failed to list invites: %w", err)
					}
					var result string
					for _, invite := range invites {
						result += invite.String() + "\n"
					}
					return result, nil
				},
				schema,
			),
		)
	}

	s.AddTool(mcp.NewTool(
		"create_invite",
		mcp.WithDescription("Creates a new invite for a user to join the organization"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				email := params["email"].(string)
				role := params["role"].(string)
				invite, err := client.CreateInvite(email, role)
				if err != nil {
					return nil, fmt.Errorf("failed to create invite: %w", err)
				}
				return invite.String(), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "email", Required: true, Type: reflect.String, Description: "Email address to invite"},
					{Name: "role", Required: true, Type: reflect.String, Description: "Role for the invited user (e.g., owner, member)"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"retrieve_invite",
		mcp.WithDescription("Retrieves details of a specific invite by ID"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				inviteID := params["inviteId"].(string)
				invite, err := client.RetrieveInvite(inviteID)
				if err != nil {
					return nil, fmt.Errorf("failed to retrieve invite: %w", err)
				}
				return invite.String(), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "inviteId", Required: true, Type: reflect.String, Description: "Invite ID"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"delete_invite",
		mcp.WithDescription("Deletes a pending invite by ID"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				inviteID := params["inviteId"].(string)
				err := client.DeleteInvite(inviteID)
				if err != nil {
					return nil, fmt.Errorf("failed to delete invite: %w", err)
				}
				return fmt.Sprintf("Invite %s deleted", inviteID), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "inviteId", Required: true, Type: reflect.String, Description: "Invite ID"},
				},
			},
		),
	)

	// --- Usage/Billing ---
	{
		schema := ParamSchema{
			Fields: []ParamField{
				{Name: "type", Required: true, Type: reflect.String, Description: "Usage type (completions, embeddings, moderations, images, audio_speeches, audio_transcriptions, vector_stores, code_interpreter, costs)"},
				{Name: "startTime", Required: false, Type: reflect.String, Description: "Start time (RFC3339)"},
				{Name: "endTime", Required: false, Type: reflect.String, Description: "End time (RFC3339)"},
			},
		}
		s.AddTool(mcp.NewTool(
			"get_usage",
			mcp.WithDescription("Retrieves usage statistics for a given type (completions, embeddings, etc.)"),
			mcp.WithString("type", mcp.Required(), mcp.Description("Usage type (completions, embeddings, moderations, images, audio_speeches, audio_transcriptions, vector_stores, code_interpreter, costs")),
			mcp.WithString("startTime", mcp.Description("Start time (RFC3339)")),
			mcp.WithString("endTime", mcp.Description("End time (RFC3339)")),
		),
			GenericToolHandler(
				func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
					typeStr := params["type"].(string)
					startTime := ""
					if v, ok := params["startTime"]; ok {
						startTime = v.(string)
					}
					endTime := ""
					if v, ok := params["endTime"]; ok {
						endTime = v.(string)
					}
					queryParams := map[string]string{}
					if startTime != "" {
						queryParams["start_time"] = startTime
					}
					if endTime != "" {
						queryParams["end_time"] = endTime
					}
					switch typeStr {
					case "completions":
						usage, err := client.GetCompletionsUsage(queryParams)
						if err != nil {
							return nil, fmt.Errorf("failed to get completions usage: %w", err)
						}
						return fmt.Sprintf("%+v", usage), nil
					case "embeddings":
						usage, err := client.GetEmbeddingsUsage(queryParams)
						if err != nil {
							return nil, fmt.Errorf("failed to get embeddings usage: %w", err)
						}
						return fmt.Sprintf("%+v", usage), nil
					case "moderations":
						usage, err := client.GetModerationsUsage(queryParams)
						if err != nil {
							return nil, fmt.Errorf("failed to get moderations usage: %w", err)
						}
						return fmt.Sprintf("%+v", usage), nil
					case "images":
						usage, err := client.GetImagesUsage(queryParams)
						if err != nil {
							return nil, fmt.Errorf("failed to get images usage: %w", err)
						}
						return fmt.Sprintf("%+v", usage), nil
					case "audio_speeches":
						usage, err := client.GetAudioSpeechesUsage(queryParams)
						if err != nil {
							return nil, fmt.Errorf("failed to get audio speeches usage: %w", err)
						}
						return fmt.Sprintf("%+v", usage), nil
					case "audio_transcriptions":
						usage, err := client.GetAudioTranscriptionsUsage(queryParams)
						if err != nil {
							return nil, fmt.Errorf("failed to get audio transcriptions usage: %w", err)
						}
						return fmt.Sprintf("%+v", usage), nil
					case "vector_stores":
						usage, err := client.GetVectorStoresUsage(queryParams)
						if err != nil {
							return nil, fmt.Errorf("failed to get vector stores usage: %w", err)
						}
						return fmt.Sprintf("%+v", usage), nil
					case "code_interpreter":
						usage, err := client.GetCodeInterpreterUsage(queryParams)
						if err != nil {
							return nil, fmt.Errorf("failed to get code interpreter usage: %w", err)
						}
						return fmt.Sprintf("%+v", usage), nil
					case "costs":
						usage, err := client.GetCostsUsage(queryParams)
						if err != nil {
							return nil, fmt.Errorf("failed to get costs usage: %w", err)
						}
						return fmt.Sprintf("%+v", usage), nil
					default:
						return nil, fmt.Errorf("unsupported usage type: %s", typeStr)
					}
				},
				schema,
			),
		)
	}

	s.AddTool(mcp.NewTool(
		"retrieve_project_api_key",
		mcp.WithDescription("Retrieves a specific API key from a project by ID"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				projectID := params["projectId"].(string)
				apiKeyID := params["apiKeyId"].(string)
				key, err := client.RetrieveProjectApiKey(projectID, apiKeyID)
				if err != nil {
					return nil, fmt.Errorf("failed to retrieve project API key: %w", err)
				}
				return key.String(), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "apiKeyId", Required: true, Type: reflect.String, Description: "API Key ID"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"retrieve_project_service_account",
		mcp.WithDescription("Retrieves a specific service account from a project by ID"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				projectID := params["projectId"].(string)
				serviceAccountID := params["serviceAccountId"].(string)
				account, err := client.RetrieveProjectServiceAccount(projectID, serviceAccountID)
				if err != nil {
					return nil, fmt.Errorf("failed to retrieve project service account: %w", err)
				}
				return account.String(), nil
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "serviceAccountId", Required: true, Type: reflect.String, Description: "Service Account ID"},
				},
			},
		),
	)
}
