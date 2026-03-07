package cmd

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Helper
func createMockProjectApiKey(id, name string) openaiorgs.ProjectApiKey {
	return openaiorgs.ProjectApiKey{
		Object:        "organization.project.api_key",
		ID:            id,
		Name:          name,
		RedactedValue: "sk-...abc",
		CreatedAt:     openaiorgs.UnixSeconds(time.Now()),
		Owner: openaiorgs.Owner{
			Object: "owner",
			ID:     "owner_1",
			Name:   "Test Owner",
			Type:   openaiorgs.OwnerTypeUser,
			User: &openaiorgs.User{
				Object: "user",
				ID:     "user_1",
				Name:   "Test User",
				Email:  "test@example.com",
			},
		},
	}
}

func TestListProjectAPIKeysCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful list pretty",
			args:       []string{"project-api-keys", "list", "--project-id", "proj_123"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.ProjectApiKey]{
				Object:  "list",
				Data:    []openaiorgs.ProjectApiKey{createMockProjectApiKey("key_1", "My API Key")},
				FirstID: "key_1",
				LastID:  "key_1",
				HasMore: false,
			},
			wantContains: []string{"ID | Name | Created At | Owner", "key_1", "My API Key"},
		},
		{
			name:       "empty list",
			args:       []string{"project-api-keys", "list", "--project-id", "proj_123"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.ProjectApiKey]{
				Object: "list",
				Data:   []openaiorgs.ProjectApiKey{},
			},
			wantContains: []string{"ID | Name | Created At | Owner"},
		},
		{
			name:       "error from API",
			args:       []string{"project-api-keys", "list", "--project-id", "proj_123"},
			statusCode: 500,
			response:   map[string]string{"error": "API error"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/projects/proj_123/api_keys", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectAPIKeysCommand(), tt.args)
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

func TestListProjectAPIKeysJSONCommand(t *testing.T) {
	h := newCmdTestHelper(t)
	defer h.cleanup()

	key := createMockProjectApiKey("key_1", "My API Key")
	h.mockResponse("GET", "/organization/projects/proj_123/api_keys", 200,
		openaiorgs.ListResponse[openaiorgs.ProjectApiKey]{
			Object:  "list",
			Data:    []openaiorgs.ProjectApiKey{key},
			FirstID: "key_1",
			LastID:  "key_1",
			HasMore: false,
		})

	output := captureOutput(func() {
		err := h.runCmd(ProjectAPIKeysCommand(), []string{"--output", "json", "project-api-keys", "list", "--project-id", "proj_123"})
		if err != nil {
			t.Errorf("runCmd() error = %v", err)
		}
	})

	// Verify valid JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &result); err != nil {
		t.Errorf("Expected valid JSON output, got: %s", output)
	}
	if !strings.Contains(output, "key_1") {
		t.Errorf("Expected output to contain key_1, got: %s", output)
	}
}

func TestRetrieveProjectAPIKeyCommand(t *testing.T) {
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
			args:         []string{"project-api-keys", "retrieve", "--project-id", "proj_123", "--id", "key_1"},
			statusCode:   200,
			response:     createMockProjectApiKey("key_1", "My API Key"),
			wantContains: []string{"ID | Name | Created At | Owner", "key_1", "My API Key"},
		},
		{
			name:       "error from API",
			args:       []string{"project-api-keys", "retrieve", "--project-id", "proj_123", "--id", "key_1"},
			statusCode: 404,
			response:   map[string]string{"error": "not found"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/projects/proj_123/api_keys/key_1", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectAPIKeysCommand(), tt.args)
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

func TestDeleteProjectAPIKeyCommand(t *testing.T) {
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
			args:         []string{"project-api-keys", "delete", "--project-id", "proj_123", "--id", "key_1"},
			statusCode:   200,
			response:     nil,
			wantContains: []string{"Successfully deleted project API key key_1"},
		},
		{
			name:       "error from API",
			args:       []string{"project-api-keys", "delete", "--project-id", "proj_123", "--id", "key_1"},
			statusCode: 500,
			response:   map[string]string{"error": "delete failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("DELETE", "/organization/projects/proj_123/api_keys/key_1", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectAPIKeysCommand(), tt.args)
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
