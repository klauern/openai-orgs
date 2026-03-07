package cmd

import (
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Helper to create a mock user
func createMockUser(id, email, name, role string) openaiorgs.User {
	return openaiorgs.User{
		Object:  "organization.user",
		ID:      id,
		Email:   email,
		Name:    name,
		Role:    role,
		AddedAt: openaiorgs.UnixSeconds(time.Now()),
	}
}

func TestListUsersCommand(t *testing.T) {
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
			args:       []string{"users", "list"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.User]{
				Object: "list",
				Data: []openaiorgs.User{
					createMockUser("user_123", "alice@example.com", "Alice", "owner"),
					createMockUser("user_456", "bob@example.com", "Bob", "member"),
				},
				FirstID: "user_123",
				LastID:  "user_456",
				HasMore: false,
			},
			wantContains: []string{
				"ID | Email | Name | Role",
				"user_123",
				"alice@example.com",
			},
		},
		{
			name:       "error from API",
			args:       []string{"users", "list"},
			statusCode: 500,
			response:   map[string]string{"error": "API error"},
			wantErr:    true,
		},
		{
			name:       "empty list",
			args:       []string{"users", "list"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.User]{
				Object:  "list",
				Data:    []openaiorgs.User{},
				HasMore: false,
			},
			wantContains: []string{"ID | Email | Name | Role"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/users", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(UsersCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}

			if !tt.wantErr {
				h.assertRequest("GET", "/organization/users", 1)
			}
		})
	}
}

func TestRetrieveUserCommand(t *testing.T) {
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
			args:       []string{"users", "retrieve", "--id", "user_123"},
			statusCode: 200,
			response:   createMockUser("user_123", "alice@example.com", "Alice", "owner"),
			wantContains: []string{
				"User details:",
				"ID: user_123",
				"Email: alice@example.com",
			},
		},
		{
			name:       "error from API",
			args:       []string{"users", "retrieve", "--id", "user_999"},
			statusCode: 404,
			response:   map[string]string{"error": "user not found"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			userID := tt.args[3]
			h.mockResponse("GET", "/organization/users/"+userID, tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(UsersCommand(), tt.args)
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

func TestDeleteUserCommand(t *testing.T) {
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
			args:         []string{"users", "delete", "--id", "user_123"},
			statusCode:   200,
			response:     nil,
			wantContains: []string{"User user_123 deleted successfully"},
		},
		{
			name:       "error from API",
			args:       []string{"users", "delete", "--id", "user_999"},
			statusCode: 500,
			response:   map[string]string{"error": "delete failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			userID := tt.args[3]
			h.mockResponse("DELETE", "/organization/users/"+userID, tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(UsersCommand(), tt.args)
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

func TestModifyUserRoleCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		mockSetup    func(h *cmdTestHelper)
		wantErr      bool
		wantContains []string
	}{
		{
			name: "successful modify",
			args: []string{"users", "modify-role", "--id", "user_123", "--role", "owner"},
			mockSetup: func(h *cmdTestHelper) {
				// ModifyUserRole uses Post[User] which expects a JSON User response
				h.mockResponse("POST", "/organization/users/user_123", 200,
					createMockUser("user_123", "alice@example.com", "Alice", "owner"))
				// Then RetrieveUser is called to show updated user
				h.mockResponse("GET", "/organization/users/user_123", 200,
					createMockUser("user_123", "alice@example.com", "Alice", "owner"))
			},
			wantContains: []string{
				"User role modified:",
				"New Role: owner",
			},
		},
		{
			name: "modify error",
			args: []string{"users", "modify-role", "--id", "user_123", "--role", "owner"},
			mockSetup: func(h *cmdTestHelper) {
				h.mockResponse("POST", "/organization/users/user_123", 500,
					map[string]string{"error": "modify failed"})
			},
			wantErr: true,
		},
		{
			name: "retrieve after modify error",
			args: []string{"users", "modify-role", "--id", "user_123", "--role", "owner"},
			mockSetup: func(h *cmdTestHelper) {
				h.mockResponse("POST", "/organization/users/user_123", 200,
					createMockUser("user_123", "alice@example.com", "Alice", "owner"))
				h.mockResponse("GET", "/organization/users/user_123", 500,
					map[string]string{"error": "retrieve failed"})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			tt.mockSetup(h)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(UsersCommand(), tt.args)
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
