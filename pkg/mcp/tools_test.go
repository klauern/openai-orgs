package mcp

import (
	"context"
	"errors"
	"fmt"
	"testing"
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
