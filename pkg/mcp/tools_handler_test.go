package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/jarcoal/httpmock"
	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// newToolTestClient sets up an httpmock-backed client and overrides the
// package-level newToolClient factory so that GenericToolHandler (and therefore
// every lambda registered via AddTools) uses the mocked client.
//
// Callers MUST defer the returned cleanup function.
// Tests using this helper must NOT call t.Parallel().
func newToolTestClient(t *testing.T) (*openaiorgs.Client, func()) {
	t.Helper()
	client := openaiorgs.NewClient("https://api.openai.com/v1", "test-token")
	client.SetRetryCount(0)
	httpmock.ActivateNonDefault(client.GetHTTPClient())

	origFactory := newToolClient
	newToolClient = func(token string) *openaiorgs.Client {
		return client
	}

	cleanup := func() {
		newToolClient = origFactory
		httpmock.DeactivateAndReset()
	}
	return client, cleanup
}

// callTool builds a JSON-RPC "tools/call" message and sends it through the
// MCPServer. It returns the parsed response or error response.
func callTool(t *testing.T, s *server.MCPServer, ctx context.Context, toolName string, args map[string]any) mcp.JSONRPCMessage {
	t.Helper()

	reqMap := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]any{
			"name":      toolName,
			"arguments": args,
		},
	}
	raw, err := json.Marshal(reqMap)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	return s.HandleMessage(ctx, raw)
}

// assertToolSuccess checks that the response is a JSONRPCResponse (not an error).
func assertToolSuccess(t *testing.T, resp mcp.JSONRPCMessage) {
	t.Helper()
	if errResp, ok := resp.(mcp.JSONRPCError); ok {
		t.Fatalf("expected success response, got error: code=%d message=%s", errResp.Error.Code, errResp.Error.Message)
	}
	if _, ok := resp.(mcp.JSONRPCResponse); !ok {
		t.Fatalf("expected JSONRPCResponse, got %T", resp)
	}
}

// setupToolServer creates an MCPServer with all tools registered and returns
// a context with an auth token set.
func setupToolServer(t *testing.T) (*server.MCPServer, context.Context) {
	t.Helper()
	s := server.NewMCPServer("test", "1.0", server.WithToolCapabilities(false))
	AddTools(s)
	ctx := context.WithValue(context.Background(), authToken{}, "test-token")
	return s, ctx
}

// emptyListResponse is the JSON body for a generic empty paginated list.
var emptyListResponse = map[string]any{
	"object":   "list",
	"data":     []any{},
	"first_id": "",
	"last_id":  "",
	"has_more": false,
}

func TestGenericToolHandler_Success(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	handler := GenericToolHandler(
		func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
			return "success-value", nil
		},
		ParamSchema{},
	)

	ctx := context.WithValue(context.Background(), authToken{}, "test-token")
	req := mcp.CallToolRequest{}
	result, err := handler(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestGenericToolHandler_ValidationError(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	handler := GenericToolHandler(
		func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
			return "ok", nil
		},
		ParamSchema{
			Fields: []ParamField{
				{Name: "required_field", Required: true, Type: 0, Description: "A required field"},
			},
		},
	)

	ctx := context.WithValue(context.Background(), authToken{}, "test-token")
	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]any{} // missing required field
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

// --- Tool handler tests via MCPServer.HandleMessage ---
// Each test exercises the actual lambda registered in AddTools, covering
// the tool handler lines in tools.go.

func TestToolHandler_ListProjects(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/projects.*",
		httpmock.NewJsonResponderOrPanic(200, emptyListResponse))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "list_projects", map[string]any{})
	assertToolSuccess(t, resp)
}

func TestToolHandler_CreateProject(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("POST", "=~.*/organization/projects$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"id": "proj-1", "object": "organization.project", "name": "test",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "create_project", map[string]any{"name": "test"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_RetrieveProject(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/projects/proj-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"id": "proj-1", "object": "organization.project", "name": "test",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "retrieve_project", map[string]any{"id": "proj-1"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_ModifyProject(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("POST", "=~.*/organization/projects/proj-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"id": "proj-1", "object": "organization.project", "name": "renamed",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "modify_project", map[string]any{"id": "proj-1", "name": "renamed"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_ArchiveProject(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("POST", "=~.*/organization/projects/proj-1/archive$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"id": "proj-1", "object": "organization.project", "name": "test", "status": "archived",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "archive_project", map[string]any{"id": "proj-1"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_ListProjectUsers(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/projects/proj-1/users.*",
		httpmock.NewJsonResponderOrPanic(200, emptyListResponse))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "list_project_users", map[string]any{"projectId": "proj-1"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_AddProjectUser(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("POST", "=~.*/organization/projects/proj-1/users$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "organization.project.user", "id": "user-1", "role": "member",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "add_project_user", map[string]any{
		"projectId": "proj-1", "userId": "user-1", "role": "member",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_RemoveProjectUser(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("DELETE", "=~.*/organization/projects/proj-1/users/user-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "remove_project_user", map[string]any{
		"projectId": "proj-1", "userId": "user-1",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_RetrieveProjectUser(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/projects/proj-1/users/user-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "organization.project.user", "id": "user-1", "role": "member",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "retrieve_project_user", map[string]any{
		"projectId": "proj-1", "userId": "user-1",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_ModifyProjectUser(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("POST", "=~.*/organization/projects/proj-1/users/user-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "organization.project.user", "id": "user-1", "role": "owner",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "modify_project_user", map[string]any{
		"projectId": "proj-1", "userId": "user-1", "role": "owner",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_ListProjectApiKeys(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/projects/proj-1/api_keys.*",
		httpmock.NewJsonResponderOrPanic(200, emptyListResponse))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "list_project_api_keys", map[string]any{"projectId": "proj-1"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_DeleteProjectApiKey(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("DELETE", "=~.*/organization/projects/proj-1/api_keys/key-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "delete_project_api_key", map[string]any{
		"projectId": "proj-1", "apiKeyId": "key-1",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_RetrieveProjectApiKey(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/projects/proj-1/api_keys/key-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "organization.project.api_key", "id": "key-1",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "retrieve_project_api_key", map[string]any{
		"projectId": "proj-1", "apiKeyId": "key-1",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_ListProjectServiceAccounts(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/projects/proj-1/service_accounts.*",
		httpmock.NewJsonResponderOrPanic(200, emptyListResponse))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "list_project_service_accounts", map[string]any{"projectId": "proj-1"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_CreateProjectServiceAccount(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("POST", "=~.*/organization/projects/proj-1/service_accounts$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "organization.project.service_account", "id": "sa-1", "name": "test-sa",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "create_project_service_account", map[string]any{
		"projectId": "proj-1", "name": "test-sa",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_DeleteProjectServiceAccount(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("DELETE", "=~.*/organization/projects/proj-1/service_accounts/sa-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "delete_project_service_account", map[string]any{
		"projectId": "proj-1", "serviceAccountId": "sa-1",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_RetrieveProjectServiceAccount(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/projects/proj-1/service_accounts/sa-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "organization.project.service_account", "id": "sa-1", "name": "test-sa",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "retrieve_project_service_account", map[string]any{
		"projectId": "proj-1", "serviceAccountId": "sa-1",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_ListUsers(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/users.*",
		httpmock.NewJsonResponderOrPanic(200, emptyListResponse))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "list_users", map[string]any{})
	assertToolSuccess(t, resp)
}

func TestToolHandler_RetrieveUser(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/users/user-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "organization.user", "id": "user-1", "name": "Test User",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "retrieve_user", map[string]any{"userId": "user-1"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_DeleteUser(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("DELETE", "=~.*/organization/users/user-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "delete_user", map[string]any{"userId": "user-1"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_ModifyUserRole(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("POST", "=~.*/organization/users/user-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "organization.user", "id": "user-1", "role": "owner",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "modify_user_role", map[string]any{
		"userId": "user-1", "role": "owner",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_ListInvites(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/invites.*",
		httpmock.NewJsonResponderOrPanic(200, emptyListResponse))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "list_invites", map[string]any{})
	assertToolSuccess(t, resp)
}

func TestToolHandler_CreateInvite(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("POST", "=~.*/organization/invites$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "organization.invite", "id": "inv-1", "email": "test@example.com", "role": "member",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "create_invite", map[string]any{
		"email": "test@example.com", "role": "member",
	})
	assertToolSuccess(t, resp)
}

func TestToolHandler_RetrieveInvite(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/invites/inv-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "organization.invite", "id": "inv-1", "email": "test@example.com",
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "retrieve_invite", map[string]any{"inviteId": "inv-1"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_DeleteInvite(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("DELETE", "=~.*/organization/invites/inv-1$",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "delete_invite", map[string]any{"inviteId": "inv-1"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_GetUsage_Completions(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/usage/completions.*",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "page", "data": []any{},
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{"type": "completions"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_GetUsage_Embeddings(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/usage/embeddings.*",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "page", "data": []any{},
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{"type": "embeddings"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_GetUsage_Moderations(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/usage/moderations.*",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "page", "data": []any{},
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{"type": "moderations"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_GetUsage_Images(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/usage/images.*",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "page", "data": []any{},
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{"type": "images"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_GetUsage_AudioSpeeches(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/usage/audio_speeches.*",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "page", "data": []any{},
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{"type": "audio_speeches"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_GetUsage_AudioTranscriptions(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/usage/audio_transcriptions.*",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "page", "data": []any{},
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{"type": "audio_transcriptions"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_GetUsage_VectorStores(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/usage/vector_stores.*",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "page", "data": []any{},
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{"type": "vector_stores"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_GetUsage_CodeInterpreter(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/usage/code_interpreter.*",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "page", "data": []any{},
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{"type": "code_interpreter"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_GetUsage_Costs(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/costs.*",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "page", "data": []any{},
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{"type": "costs"})
	assertToolSuccess(t, resp)
}

func TestToolHandler_GetUsage_UnsupportedType(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{"type": "invalid_type"})
	// This should return an error response
	if _, ok := resp.(mcp.JSONRPCError); !ok {
		t.Fatalf("expected error response for unsupported usage type, got %T", resp)
	}
}

func TestToolHandler_GetUsage_WithTimeParams(t *testing.T) {
	_, cleanup := newToolTestClient(t)
	defer cleanup()

	httpmock.RegisterResponder("GET", "=~.*/organization/usage/completions.*",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"object": "page", "data": []any{},
		}))

	s, ctx := setupToolServer(t)
	resp := callTool(t, s, ctx, "get_usage", map[string]any{
		"type":      "completions",
		"startTime": "2024-01-01T00:00:00Z",
		"endTime":   "2024-01-31T23:59:59Z",
	})
	assertToolSuccess(t, resp)
}
