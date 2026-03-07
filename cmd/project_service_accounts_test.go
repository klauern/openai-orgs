package cmd

import (
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Helper
func createMockProjectServiceAccount(id, name string) openaiorgs.ProjectServiceAccount {
	return openaiorgs.ProjectServiceAccount{
		Object:    "organization.project.service_account",
		ID:        id,
		Name:      name,
		Role:      "member",
		CreatedAt: openaiorgs.UnixSeconds(time.Now()),
	}
}

func TestListProjectServiceAccountsCommand(t *testing.T) {
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
			args:       []string{"project-service-accounts", "list", "--project-id", "proj_123"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.ProjectServiceAccount]{
				Object:  "list",
				Data:    []openaiorgs.ProjectServiceAccount{createMockProjectServiceAccount("sa_1", "My Service Account")},
				FirstID: "sa_1",
				LastID:  "sa_1",
				HasMore: false,
			},
			wantContains: []string{"ID | Name | Created At", "sa_1", "My Service Account"},
		},
		{
			name:       "empty list",
			args:       []string{"project-service-accounts", "list", "--project-id", "proj_123"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.ProjectServiceAccount]{
				Object: "list",
				Data:   []openaiorgs.ProjectServiceAccount{},
			},
			wantContains: []string{"ID | Name | Created At"},
		},
		{
			name:       "error from API",
			args:       []string{"project-service-accounts", "list", "--project-id", "proj_123"},
			statusCode: 500,
			response:   map[string]string{"error": "API error"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/projects/proj_123/service_accounts", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectServiceAccountsCommand(), tt.args)
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

func TestCreateProjectServiceAccountCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful create",
			args:       []string{"project-service-accounts", "create", "--project-id", "proj_123", "--name", "New SA"},
			statusCode: 200,
			response:   createMockProjectServiceAccount("sa_new", "New SA"),
			wantContains: []string{"Project Service Account created:", "sa_new", "New SA"},
		},
		{
			name:       "error from API",
			args:       []string{"project-service-accounts", "create", "--project-id", "proj_123", "--name", "Bad SA"},
			statusCode: 500,
			response:   map[string]string{"error": "creation failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("POST", "/organization/projects/proj_123/service_accounts", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectServiceAccountsCommand(), tt.args)
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

func TestRetrieveProjectServiceAccountCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful retrieve",
			args:       []string{"project-service-accounts", "retrieve", "--project-id", "proj_123", "--id", "sa_1"},
			statusCode: 200,
			response:   createMockProjectServiceAccount("sa_1", "My SA"),
			wantContains: []string{"Project Service Account details:", "sa_1", "My SA"},
		},
		{
			name:       "error from API",
			args:       []string{"project-service-accounts", "retrieve", "--project-id", "proj_123", "--id", "sa_1"},
			statusCode: 404,
			response:   map[string]string{"error": "not found"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/projects/proj_123/service_accounts/sa_1", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectServiceAccountsCommand(), tt.args)
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

func TestDeleteProjectServiceAccountCommand(t *testing.T) {
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
			args:         []string{"project-service-accounts", "delete", "--project-id", "proj_123", "--id", "sa_1"},
			statusCode:   200,
			response:     nil,
			wantContains: []string{"Project Service Account sa_1 deleted successfully"},
		},
		{
			name:       "error from API",
			args:       []string{"project-service-accounts", "delete", "--project-id", "proj_123", "--id", "sa_1"},
			statusCode: 500,
			response:   map[string]string{"error": "delete failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("DELETE", "/organization/projects/proj_123/service_accounts/sa_1", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectServiceAccountsCommand(), tt.args)
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
