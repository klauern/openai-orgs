package mcp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/mark3labs/mcp-go/mcp"
)

// mockProjects is a helper type for fake project list results
// Implements String() for result formatting
type mockProjects struct {
	val string
}

func (m mockProjects) String() string { return m.val }

// mockClient is a minimal interface for ListProjects
type mockClient interface {
	ListProjects(limit int, after string, activeOnly bool) (fmt.Stringer, error)
}

// handlerForTest is the core logic for list_projects, testable with a mock client
func handlerForTest(ctx context.Context, client mockClient, params map[string]any) (any, error) {
	limit := 100
	if v, ok := params["limit"]; ok {
		limit = int(v.(float64))
	}
	after := ""
	if v, ok := params["after"]; ok {
		after = v.(string)
	}
	activeOnly := false
	if v, ok := params["activeOnly"]; ok {
		activeOnly = v.(bool)
	}
	return client.ListProjects(limit, after, activeOnly)
}

// TestListProjectsHandler tests the handler logic for list_projects
func TestListProjectsHandler(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]any
		want    string
		mockFn  func(*mockClientImpl)
		wantErr bool
	}{
		{
			name: "default params",
			args: map[string]any{},
			want: "projects-default",
			mockFn: func(m *mockClientImpl) {
				m.ListProjectsFunc = func(limit int, after string, activeOnly bool) (fmt.Stringer, error) {
					if limit != 100 || after != "" || activeOnly != false {
						t.Errorf("unexpected params: %v %v %v", limit, after, activeOnly)
					}
					return mockProjects{"projects-default"}, nil
				}
			},
		},
		{
			name: "custom limit",
			args: map[string]any{"limit": float64(5)},
			want: "projects-limit-5",
			mockFn: func(m *mockClientImpl) {
				m.ListProjectsFunc = func(limit int, after string, activeOnly bool) (fmt.Stringer, error) {
					if limit != 5 {
						t.Errorf("expected limit 5, got %d", limit)
					}
					return mockProjects{"projects-limit-5"}, nil
				}
			},
		},
		{
			name: "with after",
			args: map[string]any{"after": "abc"},
			want: "projects-after-abc",
			mockFn: func(m *mockClientImpl) {
				m.ListProjectsFunc = func(limit int, after string, activeOnly bool) (fmt.Stringer, error) {
					if after != "abc" {
						t.Errorf("expected after 'abc', got %q", after)
					}
					return mockProjects{"projects-after-abc"}, nil
				}
			},
		},
		{
			name: "activeOnly true",
			args: map[string]any{"activeOnly": true},
			want: "projects-active-true",
			mockFn: func(m *mockClientImpl) {
				m.ListProjectsFunc = func(limit int, after string, activeOnly bool) (fmt.Stringer, error) {
					if !activeOnly {
						t.Errorf("expected activeOnly true")
					}
					return mockProjects{"projects-active-true"}, nil
				}
			},
		},
		{
			name:    "client error",
			args:    map[string]any{},
			wantErr: true,
			mockFn: func(m *mockClientImpl) {
				m.ListProjectsFunc = func(limit int, after string, activeOnly bool) (fmt.Stringer, error) {
					return nil, errors.New("fail")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockClientImpl{}
			if tc.mockFn != nil {
				tc.mockFn(mock)
			}
			ctx := context.Background()
			res, err := handlerForTest(ctx, mock, tc.args)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil && res != nil && res.(fmt.Stringer).String() != tc.want {
				t.Errorf("got %q, want %q", res.(fmt.Stringer).String(), tc.want)
			}
		})
	}
}

// mockClientImpl is a test double for mockClient
type mockClientImpl struct {
	ListProjectsFunc func(limit int, after string, activeOnly bool) (fmt.Stringer, error)
}

func (m *mockClientImpl) ListProjects(limit int, after string, activeOnly bool) (fmt.Stringer, error) {
	if m.ListProjectsFunc != nil {
		return m.ListProjectsFunc(limit, after, activeOnly)
	}
	return mockProjects{"default"}, nil
}

func TestGenericToolHandler_NoAuthToken(t *testing.T) {
	schema := ParamSchema{}
	handler := GenericToolHandler(
		func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
			return "ok", nil
		},
		schema,
	)
	// Call with a context that has NO auth token set
	ctx := context.Background()
	req := mcp.CallToolRequest{}
	_, err := handler(ctx, req)
	if err == nil {
		t.Fatal("expected error when auth token is missing, got nil")
	}
	if err.Error() != ErrNoAuthToken.Error() {
		t.Errorf("expected ErrNoAuthToken, got: %v", err)
	}
}

func TestParamSchema_ToMCPParameterSchema(t *testing.T) {
	schema := ParamSchema{
		Fields: []ParamField{
			{Name: "name", Required: true, Type: reflect.String, Description: "A name"},
			{Name: "count", Required: false, Type: reflect.Float64, Description: "A count"},
			{Name: "active", Required: false, Type: reflect.Bool, Description: "Active flag"},
		},
	}

	result := schema.ToMCPParameterSchema()

	// Check type
	if result["type"] != "object" {
		t.Errorf("expected type 'object', got %v", result["type"])
	}

	// Check properties exist
	props, ok := result["properties"].(map[string]any)
	if !ok {
		t.Fatal("properties is not a map")
	}

	if len(props) != 3 {
		t.Errorf("expected 3 properties, got %d", len(props))
	}

	// Check required
	required, ok := result["required"].([]string)
	if !ok {
		t.Fatal("required is not a string slice")
	}
	if len(required) != 1 || required[0] != "name" {
		t.Errorf("expected required=[name], got %v", required)
	}
}

func TestParamSchema_ExtractAndValidate(t *testing.T) {
	schema := ParamSchema{
		Fields: []ParamField{
			{Name: "name", Required: true, Type: reflect.String, Description: "Name"},
			{Name: "count", Required: false, Type: reflect.Float64, Description: "Count"},
		},
	}

	tests := []struct {
		name    string
		args    map[string]any
		wantErr bool
	}{
		{"valid with all params", map[string]any{"name": "test", "count": float64(5)}, false},
		{"valid with required only", map[string]any{"name": "test"}, false},
		{"missing required", map[string]any{"count": float64(5)}, true},
		{"wrong type for required", map[string]any{"name": 123}, true},
		{"nil arguments", nil, true},
		{"extra params ignored", map[string]any{"name": "test", "extra": "value"}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := mcp.CallToolRequest{}
			req.Params.Arguments = tc.args
			_, err := schema.ExtractAndValidate(req)
			if (err != nil) != tc.wantErr {
				t.Errorf("ExtractAndValidate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
