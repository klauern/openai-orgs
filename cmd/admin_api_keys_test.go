package cmd

import (
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Helper to create a mock admin API key
func createMockAdminAPIKey(id, name, redactedValue string, scopes []string) openaiorgs.AdminAPIKey {
	now := time.Now()
	return openaiorgs.AdminAPIKey{
		Object:        "organization.admin_api_key",
		ID:            id,
		Name:          name,
		RedactedValue: redactedValue,
		CreatedAt:     openaiorgs.UnixSeconds(now),
		LastUsedAt:    openaiorgs.UnixSeconds(now),
		Scopes:        scopes,
	}
}

func TestListAdminAPIKeysCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful list with scopes",
			args:       []string{"admin-api-keys", "list"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.AdminAPIKey]{
				Object: "list",
				Data: []openaiorgs.AdminAPIKey{
					createMockAdminAPIKey("key_123", "My API Key", "sk-...abc", []string{"api.read", "api.write"}),
				},
				FirstID: "key_123",
				LastID:  "key_123",
				HasMore: false,
			},
			wantContains: []string{
				"ID | Name | Redacted Value | Created At | Last Used At | Scopes",
				"key_123",
				"api.read, api.write",
			},
		},
		{
			name:       "error from API",
			args:       []string{"admin-api-keys", "list"},
			statusCode: 500,
			response:   map[string]string{"error": "API error"},
			wantErr:    true,
		},
		{
			name:       "empty list",
			args:       []string{"admin-api-keys", "list"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.AdminAPIKey]{
				Object:  "list",
				Data:    []openaiorgs.AdminAPIKey{},
				HasMore: false,
			},
			wantContains: []string{"ID | Name | Redacted Value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/admin_api_keys", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(AdminAPIKeysCommand(), tt.args)
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

func TestCreateAdminAPIKeyCommand(t *testing.T) {
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
			args:       []string{"admin-api-keys", "create", "--name", "My New Key", "--scopes", "api.read", "--scopes", "api.write"},
			statusCode: 200,
			response:   createMockAdminAPIKey("key_123", "My New Key", "sk-...xyz", []string{"api.read", "api.write"}),
			wantContains: []string{
				"API Key created:",
				"ID: key_123",
				"Scopes: api.read, api.write",
			},
		},
		{
			name:       "error from API",
			args:       []string{"admin-api-keys", "create", "--name", "Bad Key", "--scopes", "api.read"},
			statusCode: 500,
			response:   map[string]string{"error": "create failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("POST", "/organization/admin_api_keys", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(AdminAPIKeysCommand(), tt.args)
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

func TestRetrieveAdminAPIKeyCommand(t *testing.T) {
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
			args:       []string{"admin-api-keys", "retrieve", "--id", "key_123"},
			statusCode: 200,
			response:   createMockAdminAPIKey("key_123", "My Key", "sk-...abc", []string{"api.read"}),
			wantContains: []string{
				"API Key details:",
				"ID: key_123",
				"Scopes: api.read",
				"Last Used At:",
			},
		},
		{
			name:       "error from API",
			args:       []string{"admin-api-keys", "retrieve", "--id", "key_999"},
			statusCode: 404,
			response:   map[string]string{"error": "key not found"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			keyID := tt.args[3]
			h.mockResponse("GET", "/organization/admin_api_keys/"+keyID, tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(AdminAPIKeysCommand(), tt.args)
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

func TestDeleteAdminAPIKeyCommand(t *testing.T) {
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
			args:         []string{"admin-api-keys", "delete", "--id", "key_123"},
			statusCode:   200,
			response:     nil,
			wantContains: []string{"API Key key_123 deleted successfully"},
		},
		{
			name:       "error from API",
			args:       []string{"admin-api-keys", "delete", "--id", "key_999"},
			statusCode: 500,
			response:   map[string]string{"error": "delete failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			keyID := tt.args[3]
			h.mockResponse("DELETE", "/organization/admin_api_keys/"+keyID, tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(AdminAPIKeysCommand(), tt.args)
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
