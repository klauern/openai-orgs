package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Mock client interface for testing users
type mockUserClient interface {
	ListUsers(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.User], error)
	RetrieveUser(id string) (*openaiorgs.User, error)
	DeleteUser(id string) error
	ModifyUserRole(id string, role string) error
}

// Mock implementation
type mockUserClientImpl struct {
	ListUsersFunc      func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.User], error)
	RetrieveUserFunc   func(id string) (*openaiorgs.User, error)
	DeleteUserFunc     func(id string) error
	ModifyUserRoleFunc func(id string, role string) error
}

func (m *mockUserClientImpl) ListUsers(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.User], error) {
	if m.ListUsersFunc != nil {
		return m.ListUsersFunc(limit, after)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUserClientImpl) RetrieveUser(id string) (*openaiorgs.User, error) {
	if m.RetrieveUserFunc != nil {
		return m.RetrieveUserFunc(id)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUserClientImpl) DeleteUser(id string) error {
	if m.DeleteUserFunc != nil {
		return m.DeleteUserFunc(id)
	}
	return fmt.Errorf("not implemented")
}

func (m *mockUserClientImpl) ModifyUserRole(id string, role string) error {
	if m.ModifyUserRoleFunc != nil {
		return m.ModifyUserRoleFunc(id, role)
	}
	return fmt.Errorf("not implemented")
}

// Testable handler functions

func listUsersHandler(client mockUserClient, limit int, after string) error {
	users, err := client.ListUsers(limit, after)
	if err != nil {
		return wrapError("list users", err)
	}

	data := TableData{
		Headers: []string{"ID", "Email", "Name", "Role"},
		Rows:    make([][]string, len(users.Data)),
	}

	for i, user := range users.Data {
		data.Rows[i] = []string{user.ID, user.Email, user.Name, user.Role}
	}

	printTableData(data)
	return nil
}

func retrieveUserHandler(client mockUserClient, id string) error {
	user, err := client.RetrieveUser(id)
	if err != nil {
		return wrapError("retrieve user", err)
	}

	fmt.Printf("User details:\nID: %s\nEmail: %s\nName: %s\nRole: %s\n",
		user.ID, user.Email, user.Name, user.Role)
	return nil
}

func deleteUserHandler(client mockUserClient, id string) error {
	err := client.DeleteUser(id)
	if err != nil {
		return wrapError("delete user", err)
	}

	fmt.Printf("User %s deleted successfully\n", id)
	return nil
}

func modifyUserRoleHandler(client mockUserClient, id, role string) error {
	err := client.ModifyUserRole(id, role)
	if err != nil {
		return wrapError("modify user role", err)
	}

	user, err := client.RetrieveUser(id)
	if err != nil {
		return wrapError("retrieve updated user", err)
	}

	fmt.Printf("User role modified:\nID: %s\nEmail: %s\nName: %s\nNew Role: %s\n",
		user.ID, user.Email, user.Name, user.Role)
	return nil
}

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

// Tests

func TestListUsersHandler(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		after   string
		mockFn  func(*mockUserClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:  "successful list",
			limit: 10,
			after: "",
			mockFn: func(m *mockUserClientImpl) {
				m.ListUsersFunc = func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.User], error) {
					if limit != 10 || after != "" {
						t.Errorf("unexpected params: limit=%d, after=%s", limit, after)
					}
					return &openaiorgs.ListResponse[openaiorgs.User]{
						Object: "list",
						Data: []openaiorgs.User{
							createMockUser("user_123", "alice@example.com", "Alice", "owner"),
							createMockUser("user_456", "bob@example.com", "Bob", "member"),
						},
						FirstID: "user_123",
						LastID:  "user_456",
						HasMore: false,
					}, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "ID | Email | Name | Role") {
					t.Errorf("Expected table headers in output, got: %s", output)
				}
				if !strings.Contains(output, "user_123") {
					t.Errorf("Expected user_123 in output, got: %s", output)
				}
				if !strings.Contains(output, "alice@example.com") {
					t.Errorf("Expected alice@example.com in output, got: %s", output)
				}
			},
		},
		{
			name:  "error from client",
			limit: 10,
			after: "",
			mockFn: func(m *mockUserClientImpl) {
				m.ListUsersFunc = func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.User], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
		{
			name:  "empty list",
			limit: 10,
			after: "",
			mockFn: func(m *mockUserClientImpl) {
				m.ListUsersFunc = func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.User], error) {
					return &openaiorgs.ListResponse[openaiorgs.User]{
						Object:  "list",
						Data:    []openaiorgs.User{},
						HasMore: false,
					}, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "ID | Email | Name | Role") {
					t.Errorf("Expected table headers even for empty list, got: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listUsersHandler(mock, tt.limit, tt.after)
				if (err != nil) != tt.wantErr {
					t.Errorf("listUsersHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestRetrieveUserHandler(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*mockUserClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name: "successful retrieve",
			id:   "user_123",
			mockFn: func(m *mockUserClientImpl) {
				m.RetrieveUserFunc = func(id string) (*openaiorgs.User, error) {
					if id != "user_123" {
						t.Errorf("unexpected id: %s", id)
					}
					user := createMockUser("user_123", "alice@example.com", "Alice", "owner")
					return &user, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "User details:") {
					t.Errorf("Expected 'User details:' in output, got: %s", output)
				}
				if !strings.Contains(output, "ID: user_123") {
					t.Errorf("Expected 'ID: user_123' in output, got: %s", output)
				}
				if !strings.Contains(output, "Email: alice@example.com") {
					t.Errorf("Expected email in output, got: %s", output)
				}
			},
		},
		{
			name: "error from client",
			id:   "user_999",
			mockFn: func(m *mockUserClientImpl) {
				m.RetrieveUserFunc = func(id string) (*openaiorgs.User, error) {
					return nil, fmt.Errorf("user not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := retrieveUserHandler(mock, tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("retrieveUserHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*mockUserClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name: "successful delete",
			id:   "user_123",
			mockFn: func(m *mockUserClientImpl) {
				m.DeleteUserFunc = func(id string) error {
					if id != "user_123" {
						t.Errorf("unexpected id: %s", id)
					}
					return nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "User user_123 deleted successfully") {
					t.Errorf("Expected delete success message, got: %s", output)
				}
			},
		},
		{
			name: "error from client",
			id:   "user_999",
			mockFn: func(m *mockUserClientImpl) {
				m.DeleteUserFunc = func(id string) error {
					return fmt.Errorf("delete failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := deleteUserHandler(mock, tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("deleteUserHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestModifyUserRoleHandler(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		role    string
		mockFn  func(*mockUserClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name: "successful modify",
			id:   "user_123",
			role: "owner",
			mockFn: func(m *mockUserClientImpl) {
				m.ModifyUserRoleFunc = func(id string, role string) error {
					if id != "user_123" || role != "owner" {
						t.Errorf("unexpected params: id=%s, role=%s", id, role)
					}
					return nil
				}
				m.RetrieveUserFunc = func(id string) (*openaiorgs.User, error) {
					user := createMockUser("user_123", "alice@example.com", "Alice", "owner")
					return &user, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "User role modified:") {
					t.Errorf("Expected 'User role modified:' in output, got: %s", output)
				}
				if !strings.Contains(output, "New Role: owner") {
					t.Errorf("Expected 'New Role: owner' in output, got: %s", output)
				}
			},
		},
		{
			name: "modify error",
			id:   "user_123",
			role: "owner",
			mockFn: func(m *mockUserClientImpl) {
				m.ModifyUserRoleFunc = func(id string, role string) error {
					return fmt.Errorf("modify failed")
				}
			},
			wantErr: true,
		},
		{
			name: "retrieve after modify error",
			id:   "user_123",
			role: "owner",
			mockFn: func(m *mockUserClientImpl) {
				m.ModifyUserRoleFunc = func(id string, role string) error {
					return nil
				}
				m.RetrieveUserFunc = func(id string) (*openaiorgs.User, error) {
					return nil, fmt.Errorf("retrieve failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := modifyUserRoleHandler(mock, tt.id, tt.role)
				if (err != nil) != tt.wantErr {
					t.Errorf("modifyUserRoleHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}
