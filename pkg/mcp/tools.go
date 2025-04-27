package mcp

import (
	"context"
	"fmt"
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

// ExtractAndValidate extracts and validates parameters from a CallToolRequest
// Returns a map of validated parameters or an error if validation fails
func (ps *ParamSchema) ExtractAndValidate(req mcp.CallToolRequest) (map[string]any, error) {
	params := make(map[string]any)
	args := req.Params.Arguments
	for _, field := range ps.Fields {
		val, ok := args[field.Name]
		if !ok {
			if field.Required {
				return nil, fmt.Errorf("missing required parameter: %s", field.Name)
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
	// Project Management
	s.AddTool(mcp.NewTool(
		"list_projects",
		mcp.WithDescription("Lists all projects for the authenticated user"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// Extract parameters with defaults
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
			ParamSchema{
				Fields: []ParamField{
					{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of projects to return (default 100)"},
					{Name: "after", Required: false, Type: reflect.String, Description: "Project ID to start after (for pagination)"},
					{Name: "activeOnly", Required: false, Type: reflect.Bool, Description: "If true, only return active projects"},
				},
			},
		),
	)

	// --- Project Management ---
	s.AddTool(mcp.NewTool(
		"create_project",
		mcp.WithDescription("Creates a new project in the organization"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement create_project logic
				return nil, fmt.Errorf("create_project not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "name", Required: true, Type: reflect.String, Description: "Project name"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"retrieve_project",
		mcp.WithDescription("Retrieves details of a specific project by ID"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement retrieve_project logic
				return nil, fmt.Errorf("retrieve_project not implemented")
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
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement modify_project logic
				return nil, fmt.Errorf("modify_project not implemented")
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
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement archive_project logic
				return nil, fmt.Errorf("archive_project not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "id", Required: true, Type: reflect.String, Description: "Project ID"},
				},
			},
		),
	)

	// --- Project User Management ---
	s.AddTool(mcp.NewTool(
		"list_project_users",
		mcp.WithDescription("Lists all users for a given project"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement list_project_users logic
				return nil, fmt.Errorf("list_project_users not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of users to return"},
					{Name: "after", Required: false, Type: reflect.String, Description: "User ID to start after (for pagination)"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"add_project_user",
		mcp.WithDescription("Adds a user to a project with a specific role"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement add_project_user logic
				return nil, fmt.Errorf("add_project_user not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
					{Name: "role", Required: true, Type: reflect.String, Description: "Role (e.g., owner, member)"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"remove_project_user",
		mcp.WithDescription("Removes a user from a project"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement remove_project_user logic
				return nil, fmt.Errorf("remove_project_user not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"retrieve_project_user",
		mcp.WithDescription("Retrieves a specific user from a project"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement retrieve_project_user logic
				return nil, fmt.Errorf("retrieve_project_user not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"modify_project_user",
		mcp.WithDescription("Modifies a user's role in a project"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement modify_project_user logic
				return nil, fmt.Errorf("modify_project_user not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "userId", Required: true, Type: reflect.String, Description: "User ID"},
					{Name: "role", Required: true, Type: reflect.String, Description: "New role (e.g., owner, member)"},
				},
			},
		),
	)

	// --- Project API Keys ---
	s.AddTool(mcp.NewTool(
		"list_project_api_keys",
		mcp.WithDescription("Lists all API keys for a given project"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement list_project_api_keys logic
				return nil, fmt.Errorf("list_project_api_keys not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of API keys to return"},
					{Name: "after", Required: false, Type: reflect.String, Description: "API key ID to start after (for pagination)"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"delete_project_api_key",
		mcp.WithDescription("Deletes a specific API key from a project"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement delete_project_api_key logic
				return nil, fmt.Errorf("delete_project_api_key not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "apiKeyId", Required: true, Type: reflect.String, Description: "API Key ID"},
				},
			},
		),
	)

	// --- Project Service Accounts ---
	s.AddTool(mcp.NewTool(
		"list_project_service_accounts",
		mcp.WithDescription("Lists all service accounts for a given project"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement list_project_service_accounts logic
				return nil, fmt.Errorf("list_project_service_accounts not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of service accounts to return"},
					{Name: "after", Required: false, Type: reflect.String, Description: "Service account ID to start after (for pagination)"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"create_project_service_account",
		mcp.WithDescription("Creates a new service account for a project"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement create_project_service_account logic
				return nil, fmt.Errorf("create_project_service_account not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "name", Required: true, Type: reflect.String, Description: "Service account name"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"delete_project_service_account",
		mcp.WithDescription("Deletes a service account from a project"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement delete_project_service_account logic
				return nil, fmt.Errorf("delete_project_service_account not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "projectId", Required: true, Type: reflect.String, Description: "Project ID"},
					{Name: "serviceAccountId", Required: true, Type: reflect.String, Description: "Service Account ID"},
				},
			},
		),
	)

	// --- User Management ---
	s.AddTool(mcp.NewTool(
		"list_users",
		mcp.WithDescription("Lists all users in the organization"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement list_users logic
				return nil, fmt.Errorf("list_users not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of users to return"},
					{Name: "after", Required: false, Type: reflect.String, Description: "User ID to start after (for pagination)"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"retrieve_user",
		mcp.WithDescription("Retrieves details of a specific user by ID"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement retrieve_user logic
				return nil, fmt.Errorf("retrieve_user not implemented")
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
				// TODO: Implement delete_user logic
				return nil, fmt.Errorf("delete_user not implemented")
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
				// TODO: Implement modify_user_role logic
				return nil, fmt.Errorf("modify_user_role not implemented")
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
	s.AddTool(mcp.NewTool(
		"list_invites",
		mcp.WithDescription("Lists all pending invites in the organization"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement list_invites logic
				return nil, fmt.Errorf("list_invites not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "limit", Required: false, Type: reflect.Float64, Description: "Maximum number of invites to return"},
					{Name: "after", Required: false, Type: reflect.String, Description: "Invite ID to start after (for pagination)"},
				},
			},
		),
	)

	s.AddTool(mcp.NewTool(
		"create_invite",
		mcp.WithDescription("Creates a new invite for a user to join the organization"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement create_invite logic
				return nil, fmt.Errorf("create_invite not implemented")
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
		"delete_invite",
		mcp.WithDescription("Deletes a pending invite by ID"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement delete_invite logic
				return nil, fmt.Errorf("delete_invite not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "inviteId", Required: true, Type: reflect.String, Description: "Invite ID"},
				},
			},
		),
	)

	// --- Usage/Billing ---
	s.AddTool(mcp.NewTool(
		"get_usage",
		mcp.WithDescription("Retrieves usage statistics for a given type (completions, embeddings, etc.)"),
	),
		GenericToolHandler(
			func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
				// TODO: Implement get_usage logic
				return nil, fmt.Errorf("get_usage not implemented")
			},
			ParamSchema{
				Fields: []ParamField{
					{Name: "type", Required: true, Type: reflect.String, Description: "Usage type (completions, embeddings, moderations, images, audio_speeches, audio_transcriptions, vector_stores, code_interpreter, costs)"},
					{Name: "startTime", Required: false, Type: reflect.String, Description: "Start time (RFC3339)"},
					{Name: "endTime", Required: false, Type: reflect.String, Description: "End time (RFC3339)"},
				},
			},
		),
	)
}
