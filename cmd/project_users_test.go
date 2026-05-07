package cmd

import (
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Helper
func createMockProjectUser(id, name, email, role string) openaiorgs.ProjectUser {
	return openaiorgs.ProjectUser{
		Object:  "organization.project.user",
		ID:      id,
		Name:    name,
		Email:   email,
		Role:    role,
		AddedAt: openaiorgs.UnixSeconds(time.Now()),
	}
}

func TestListProjectUsersCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful list",
			args:       []string{"project-users", "list", "--project-id", "proj_123"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.ProjectUser]{
				Object:  "list",
				Data:    []openaiorgs.ProjectUser{createMockProjectUser("user_1", "Alice", "alice@example.com", "member")},
				FirstID: "user_1",
				LastID:  "user_1",
				HasMore: false,
			},
			wantContains: []string{"ID | Email | Name | Role | Added At", "user_1", "alice@example.com", "Alice", "member"},
		},
		{
			name:       "empty list",
			args:       []string{"project-users", "list", "--project-id", "proj_123"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.ProjectUser]{
				Object: "list",
				Data:   []openaiorgs.ProjectUser{},
			},
			wantContains: []string{"ID | Email | Name | Role | Added At"},
		},
		{
			name:       "error from API",
			args:       []string{"project-users", "list", "--project-id", "proj_123"},
			statusCode: 500,
			response:   map[string]string{"error": "API error"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/projects/proj_123/users", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectUsersCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestCreateProjectUserCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:         "successful create",
			args:         []string{"project-users", "create", "--project-id", "proj_123", "--id", "user_1", "--role", "member"},
			statusCode:   200,
			response:     createMockProjectUser("user_1", "Alice", "alice@example.com", "member"),
			wantContains: []string{"Project User created:", "user_1", "alice@example.com", "Alice", "member"},
		},
		{
			name:       "error from API",
			args:       []string{"project-users", "create", "--project-id", "proj_123", "--id", "user_1", "--role", "member"},
			statusCode: 500,
			response:   map[string]string{"error": "creation failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("POST", "/organization/projects/proj_123/users", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectUsersCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestRetrieveProjectUserCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:         "successful retrieve",
			args:         []string{"project-users", "retrieve", "--project-id", "proj_123", "--id", "user_1"},
			statusCode:   200,
			response:     createMockProjectUser("user_1", "Alice", "alice@example.com", "owner"),
			wantContains: []string{"Project User details:", "user_1", "alice@example.com", "Alice", "owner"},
		},
		{
			name:       "error from API",
			args:       []string{"project-users", "retrieve", "--project-id", "proj_123", "--id", "user_1"},
			statusCode: 404,
			response:   map[string]string{"error": "not found"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/projects/proj_123/users/user_1", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectUsersCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestModifyProjectUserCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:         "successful modify",
			args:         []string{"project-users", "modify", "--project-id", "proj_123", "--id", "user_1", "--role", "owner"},
			statusCode:   200,
			response:     createMockProjectUser("user_1", "Alice", "alice@example.com", "owner"),
			wantContains: []string{"Project User modified:", "user_1", "owner"},
		},
		{
			name:       "error from API",
			args:       []string{"project-users", "modify", "--project-id", "proj_123", "--id", "user_1", "--role", "owner"},
			statusCode: 500,
			response:   map[string]string{"error": "modify failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("POST", "/organization/projects/proj_123/users/user_1", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectUsersCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestDeleteProjectUserCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:         "successful delete",
			args:         []string{"project-users", "delete", "--project-id", "proj_123", "--id", "user_1"},
			statusCode:   200,
			response:     nil,
			wantContains: []string{"Project User user_1 deleted successfully"},
		},
		{
			name:       "error from API",
			args:       []string{"project-users", "delete", "--project-id", "proj_123", "--id", "user_1"},
			statusCode: 500,
			response:   map[string]string{"error": "delete failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("DELETE", "/organization/projects/proj_123/users/user_1", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectUsersCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}
