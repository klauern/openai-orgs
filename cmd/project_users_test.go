package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Mock client interface for testing project users
type mockProjectUserClient interface {
	ListProjectUsers(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectUser], error)
	CreateProjectUser(projectID, userID, role string) (*openaiorgs.ProjectUser, error)
	RetrieveProjectUser(projectID, userID string) (*openaiorgs.ProjectUser, error)
	ModifyProjectUser(projectID, userID, role string) (*openaiorgs.ProjectUser, error)
	DeleteProjectUser(projectID, userID string) error
}

// Mock implementation
type mockProjectUserClientImpl struct {
	ListProjectUsersFunc    func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectUser], error)
	CreateProjectUserFunc   func(projectID, userID, role string) (*openaiorgs.ProjectUser, error)
	RetrieveProjectUserFunc func(projectID, userID string) (*openaiorgs.ProjectUser, error)
	ModifyProjectUserFunc   func(projectID, userID, role string) (*openaiorgs.ProjectUser, error)
	DeleteProjectUserFunc   func(projectID, userID string) error
}

func (m *mockProjectUserClientImpl) ListProjectUsers(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectUser], error) {
	if m.ListProjectUsersFunc != nil {
		return m.ListProjectUsersFunc(projectID, limit, after)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectUserClientImpl) CreateProjectUser(projectID, userID, role string) (*openaiorgs.ProjectUser, error) {
	if m.CreateProjectUserFunc != nil {
		return m.CreateProjectUserFunc(projectID, userID, role)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectUserClientImpl) RetrieveProjectUser(projectID, userID string) (*openaiorgs.ProjectUser, error) {
	if m.RetrieveProjectUserFunc != nil {
		return m.RetrieveProjectUserFunc(projectID, userID)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectUserClientImpl) ModifyProjectUser(projectID, userID, role string) (*openaiorgs.ProjectUser, error) {
	if m.ModifyProjectUserFunc != nil {
		return m.ModifyProjectUserFunc(projectID, userID, role)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectUserClientImpl) DeleteProjectUser(projectID, userID string) error {
	if m.DeleteProjectUserFunc != nil {
		return m.DeleteProjectUserFunc(projectID, userID)
	}
	return fmt.Errorf("not implemented")
}

// Testable handlers

func listProjectUsersHandler(client mockProjectUserClient, projectID string, limit int, after string) error {
	projectUsers, err := client.ListProjectUsers(projectID, limit, after)
	if err != nil {
		return wrapError("list project users", err)
	}

	data := TableData{
		Headers: []string{"ID", "Email", "Name", "Role", "Added At"},
		Rows:    make([][]string, len(projectUsers.Data)),
	}

	for i, user := range projectUsers.Data {
		data.Rows[i] = []string{
			user.ID,
			user.Email,
			user.Name,
			user.Role,
			user.AddedAt.String(),
		}
	}

	printTableData(data)
	return nil
}

func createProjectUserHandler(client mockProjectUserClient, projectID, userID, role string) error {
	projectUser, err := client.CreateProjectUser(projectID, userID, role)
	if err != nil {
		return wrapError("create project user", err)
	}

	fmt.Printf("Project User created:\n")
	fmt.Printf("ID: %s\nEmail: %s\nName: %s\nRole: %s\nAdded At: %s\n",
		projectUser.ID,
		projectUser.Email,
		projectUser.Name,
		projectUser.Role,
		projectUser.AddedAt.String(),
	)
	return nil
}

func retrieveProjectUserHandler(client mockProjectUserClient, projectID, userID string) error {
	projectUser, err := client.RetrieveProjectUser(projectID, userID)
	if err != nil {
		return wrapError("retrieve project user", err)
	}

	fmt.Printf("Project User details:\n")
	fmt.Printf("ID: %s\nEmail: %s\nName: %s\nRole: %s\nAdded At: %s\n",
		projectUser.ID,
		projectUser.Email,
		projectUser.Name,
		projectUser.Role,
		projectUser.AddedAt.String(),
	)
	return nil
}

func modifyProjectUserHandler(client mockProjectUserClient, projectID, userID, role string) error {
	projectUser, err := client.ModifyProjectUser(projectID, userID, role)
	if err != nil {
		return wrapError("modify project user", err)
	}

	fmt.Printf("Project User modified:\n")
	fmt.Printf("ID: %s\nEmail: %s\nName: %s\nNew Role: %s\nAdded At: %s\n",
		projectUser.ID,
		projectUser.Email,
		projectUser.Name,
		projectUser.Role,
		projectUser.AddedAt.String(),
	)
	return nil
}

func deleteProjectUserHandler(client mockProjectUserClient, projectID, userID string) error {
	err := client.DeleteProjectUser(projectID, userID)
	if err != nil {
		return wrapError("delete project user", err)
	}

	fmt.Printf("Project User %s deleted successfully\n", userID)
	return nil
}

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

// Tests

func TestListProjectUsersHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		limit        int
		after        string
		mockFn       func(*mockProjectUserClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful list",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectUserClientImpl) {
				user := createMockProjectUser("user_1", "Alice", "alice@example.com", "member")
				m.ListProjectUsersFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectUser], error) {
					if projectID != "proj_123" {
						t.Errorf("unexpected projectID: %s", projectID)
					}
					return &openaiorgs.ListResponse[openaiorgs.ProjectUser]{
						Object:  "list",
						Data:    []openaiorgs.ProjectUser{user},
						FirstID: "user_1",
						LastID:  "user_1",
						HasMore: false,
					}, nil
				}
			},
			wantContains: []string{"ID | Email | Name | Role | Added At", "user_1", "alice@example.com", "Alice", "member"},
		},
		{
			name:      "empty list",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectUserClientImpl) {
				m.ListProjectUsersFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectUser], error) {
					return &openaiorgs.ListResponse[openaiorgs.ProjectUser]{
						Object: "list",
						Data:   []openaiorgs.ProjectUser{},
					}, nil
				}
			},
			wantContains: []string{"ID | Email | Name | Role | Added At"},
		},
		{
			name:      "error from client",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectUserClientImpl) {
				m.ListProjectUsersFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectUser], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectUserClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listProjectUsersHandler(mock, tt.projectID, tt.limit, tt.after)
				if (err != nil) != tt.wantErr {
					t.Errorf("listProjectUsersHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestCreateProjectUserHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		userID       string
		role         string
		mockFn       func(*mockProjectUserClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful create",
			projectID: "proj_123",
			userID:    "user_1",
			role:      "member",
			mockFn: func(m *mockProjectUserClientImpl) {
				m.CreateProjectUserFunc = func(projectID, userID, role string) (*openaiorgs.ProjectUser, error) {
					if projectID != "proj_123" || userID != "user_1" || role != "member" {
						t.Errorf("unexpected params: projectID=%s, userID=%s, role=%s", projectID, userID, role)
					}
					user := createMockProjectUser("user_1", "Alice", "alice@example.com", "member")
					return &user, nil
				}
			},
			wantContains: []string{"Project User created:", "user_1", "alice@example.com", "Alice", "member"},
		},
		{
			name:      "error from client",
			projectID: "proj_123",
			userID:    "user_1",
			role:      "member",
			mockFn: func(m *mockProjectUserClientImpl) {
				m.CreateProjectUserFunc = func(projectID, userID, role string) (*openaiorgs.ProjectUser, error) {
					return nil, fmt.Errorf("creation failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectUserClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := createProjectUserHandler(mock, tt.projectID, tt.userID, tt.role)
				if (err != nil) != tt.wantErr {
					t.Errorf("createProjectUserHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestRetrieveProjectUserHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		userID       string
		mockFn       func(*mockProjectUserClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful retrieve",
			projectID: "proj_123",
			userID:    "user_1",
			mockFn: func(m *mockProjectUserClientImpl) {
				m.RetrieveProjectUserFunc = func(projectID, userID string) (*openaiorgs.ProjectUser, error) {
					if projectID != "proj_123" || userID != "user_1" {
						t.Errorf("unexpected params: projectID=%s, userID=%s", projectID, userID)
					}
					user := createMockProjectUser("user_1", "Alice", "alice@example.com", "owner")
					return &user, nil
				}
			},
			wantContains: []string{"Project User details:", "user_1", "alice@example.com", "Alice", "owner"},
		},
		{
			name:      "error from client",
			projectID: "proj_123",
			userID:    "user_1",
			mockFn: func(m *mockProjectUserClientImpl) {
				m.RetrieveProjectUserFunc = func(projectID, userID string) (*openaiorgs.ProjectUser, error) {
					return nil, fmt.Errorf("not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectUserClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := retrieveProjectUserHandler(mock, tt.projectID, tt.userID)
				if (err != nil) != tt.wantErr {
					t.Errorf("retrieveProjectUserHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestModifyProjectUserHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		userID       string
		role         string
		mockFn       func(*mockProjectUserClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful modify",
			projectID: "proj_123",
			userID:    "user_1",
			role:      "owner",
			mockFn: func(m *mockProjectUserClientImpl) {
				m.ModifyProjectUserFunc = func(projectID, userID, role string) (*openaiorgs.ProjectUser, error) {
					if projectID != "proj_123" || userID != "user_1" || role != "owner" {
						t.Errorf("unexpected params: projectID=%s, userID=%s, role=%s", projectID, userID, role)
					}
					user := createMockProjectUser("user_1", "Alice", "alice@example.com", "owner")
					return &user, nil
				}
			},
			wantContains: []string{"Project User modified:", "user_1", "owner"},
		},
		{
			name:      "error from client",
			projectID: "proj_123",
			userID:    "user_1",
			role:      "owner",
			mockFn: func(m *mockProjectUserClientImpl) {
				m.ModifyProjectUserFunc = func(projectID, userID, role string) (*openaiorgs.ProjectUser, error) {
					return nil, fmt.Errorf("modify failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectUserClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := modifyProjectUserHandler(mock, tt.projectID, tt.userID, tt.role)
				if (err != nil) != tt.wantErr {
					t.Errorf("modifyProjectUserHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestDeleteProjectUserHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		userID       string
		mockFn       func(*mockProjectUserClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful delete",
			projectID: "proj_123",
			userID:    "user_1",
			mockFn: func(m *mockProjectUserClientImpl) {
				m.DeleteProjectUserFunc = func(projectID, userID string) error {
					if projectID != "proj_123" || userID != "user_1" {
						t.Errorf("unexpected params: projectID=%s, userID=%s", projectID, userID)
					}
					return nil
				}
			},
			wantContains: []string{"Project User user_1 deleted successfully"},
		},
		{
			name:      "error from client",
			projectID: "proj_123",
			userID:    "user_1",
			mockFn: func(m *mockProjectUserClientImpl) {
				m.DeleteProjectUserFunc = func(projectID, userID string) error {
					return fmt.Errorf("delete failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectUserClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := deleteProjectUserHandler(mock, tt.projectID, tt.userID)
				if (err != nil) != tt.wantErr {
					t.Errorf("deleteProjectUserHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}
