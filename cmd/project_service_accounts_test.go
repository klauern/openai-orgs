package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Mock client interface for testing project service accounts
type mockProjectServiceAccountClient interface {
	ListProjectServiceAccounts(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectServiceAccount], error)
	CreateProjectServiceAccount(projectID, name string) (*openaiorgs.ProjectServiceAccount, error)
	RetrieveProjectServiceAccount(projectID, serviceAccountID string) (*openaiorgs.ProjectServiceAccount, error)
	DeleteProjectServiceAccount(projectID, serviceAccountID string) error
}

// Mock implementation
type mockProjectServiceAccountClientImpl struct {
	ListProjectServiceAccountsFunc    func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectServiceAccount], error)
	CreateProjectServiceAccountFunc   func(projectID, name string) (*openaiorgs.ProjectServiceAccount, error)
	RetrieveProjectServiceAccountFunc func(projectID, serviceAccountID string) (*openaiorgs.ProjectServiceAccount, error)
	DeleteProjectServiceAccountFunc   func(projectID, serviceAccountID string) error
}

func (m *mockProjectServiceAccountClientImpl) ListProjectServiceAccounts(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectServiceAccount], error) {
	if m.ListProjectServiceAccountsFunc != nil {
		return m.ListProjectServiceAccountsFunc(projectID, limit, after)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectServiceAccountClientImpl) CreateProjectServiceAccount(projectID, name string) (*openaiorgs.ProjectServiceAccount, error) {
	if m.CreateProjectServiceAccountFunc != nil {
		return m.CreateProjectServiceAccountFunc(projectID, name)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectServiceAccountClientImpl) RetrieveProjectServiceAccount(projectID, serviceAccountID string) (*openaiorgs.ProjectServiceAccount, error) {
	if m.RetrieveProjectServiceAccountFunc != nil {
		return m.RetrieveProjectServiceAccountFunc(projectID, serviceAccountID)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectServiceAccountClientImpl) DeleteProjectServiceAccount(projectID, serviceAccountID string) error {
	if m.DeleteProjectServiceAccountFunc != nil {
		return m.DeleteProjectServiceAccountFunc(projectID, serviceAccountID)
	}
	return fmt.Errorf("not implemented")
}

// Testable handlers

func listProjectServiceAccountsHandler(client mockProjectServiceAccountClient, projectID string, limit int, after string) error {
	serviceAccounts, err := client.ListProjectServiceAccounts(projectID, limit, after)
	if err != nil {
		return wrapError("list project service accounts", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Created At"},
		Rows:    make([][]string, len(serviceAccounts.Data)),
	}

	for i, account := range serviceAccounts.Data {
		data.Rows[i] = []string{
			account.ID,
			account.Name,
			account.CreatedAt.String(),
		}
	}

	printTableData(data)
	return nil
}

func createProjectServiceAccountHandler(client mockProjectServiceAccountClient, projectID, name string) error {
	serviceAccount, err := client.CreateProjectServiceAccount(projectID, name)
	if err != nil {
		return wrapError("create project service account", err)
	}

	fmt.Printf("Project Service Account created:\n")
	fmt.Printf("ID: %s\nName: %s\nCreated At: %s\n",
		serviceAccount.ID,
		serviceAccount.Name,
		serviceAccount.CreatedAt.String(),
	)
	return nil
}

func retrieveProjectServiceAccountHandler(client mockProjectServiceAccountClient, projectID, serviceAccountID string) error {
	serviceAccount, err := client.RetrieveProjectServiceAccount(projectID, serviceAccountID)
	if err != nil {
		return wrapError("retrieve project service account", err)
	}

	fmt.Printf("Project Service Account details:\n")
	fmt.Printf("ID: %s\nName: %s\nCreated At: %s\n",
		serviceAccount.ID,
		serviceAccount.Name,
		serviceAccount.CreatedAt.String(),
	)
	return nil
}

func deleteProjectServiceAccountHandler(client mockProjectServiceAccountClient, projectID, serviceAccountID string) error {
	err := client.DeleteProjectServiceAccount(projectID, serviceAccountID)
	if err != nil {
		return wrapError("delete project service account", err)
	}

	fmt.Printf("Project Service Account %s deleted successfully\n", serviceAccountID)
	return nil
}

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

// Tests

func TestListProjectServiceAccountsHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		limit        int
		after        string
		mockFn       func(*mockProjectServiceAccountClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful list",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectServiceAccountClientImpl) {
				sa := createMockProjectServiceAccount("sa_1", "My Service Account")
				m.ListProjectServiceAccountsFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectServiceAccount], error) {
					if projectID != "proj_123" {
						t.Errorf("unexpected projectID: %s", projectID)
					}
					return &openaiorgs.ListResponse[openaiorgs.ProjectServiceAccount]{
						Object:  "list",
						Data:    []openaiorgs.ProjectServiceAccount{sa},
						FirstID: "sa_1",
						LastID:  "sa_1",
						HasMore: false,
					}, nil
				}
			},
			wantContains: []string{"ID | Name | Created At", "sa_1", "My Service Account"},
		},
		{
			name:      "empty list",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectServiceAccountClientImpl) {
				m.ListProjectServiceAccountsFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectServiceAccount], error) {
					return &openaiorgs.ListResponse[openaiorgs.ProjectServiceAccount]{
						Object: "list",
						Data:   []openaiorgs.ProjectServiceAccount{},
					}, nil
				}
			},
			wantContains: []string{"ID | Name | Created At"},
		},
		{
			name:      "error from client",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectServiceAccountClientImpl) {
				m.ListProjectServiceAccountsFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectServiceAccount], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectServiceAccountClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listProjectServiceAccountsHandler(mock, tt.projectID, tt.limit, tt.after)
				if (err != nil) != tt.wantErr {
					t.Errorf("listProjectServiceAccountsHandler() error = %v, wantErr %v", err, tt.wantErr)
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

func TestCreateProjectServiceAccountHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		saName       string
		mockFn       func(*mockProjectServiceAccountClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful create",
			projectID: "proj_123",
			saName:    "New SA",
			mockFn: func(m *mockProjectServiceAccountClientImpl) {
				m.CreateProjectServiceAccountFunc = func(projectID, name string) (*openaiorgs.ProjectServiceAccount, error) {
					if projectID != "proj_123" || name != "New SA" {
						t.Errorf("unexpected params: projectID=%s, name=%s", projectID, name)
					}
					sa := createMockProjectServiceAccount("sa_new", "New SA")
					return &sa, nil
				}
			},
			wantContains: []string{"Project Service Account created:", "sa_new", "New SA"},
		},
		{
			name:      "error from client",
			projectID: "proj_123",
			saName:    "Bad SA",
			mockFn: func(m *mockProjectServiceAccountClientImpl) {
				m.CreateProjectServiceAccountFunc = func(projectID, name string) (*openaiorgs.ProjectServiceAccount, error) {
					return nil, fmt.Errorf("creation failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectServiceAccountClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := createProjectServiceAccountHandler(mock, tt.projectID, tt.saName)
				if (err != nil) != tt.wantErr {
					t.Errorf("createProjectServiceAccountHandler() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRetrieveProjectServiceAccountHandler(t *testing.T) {
	tests := []struct {
		name             string
		projectID        string
		serviceAccountID string
		mockFn           func(*mockProjectServiceAccountClientImpl)
		wantErr          bool
		wantContains     []string
	}{
		{
			name:             "successful retrieve",
			projectID:        "proj_123",
			serviceAccountID: "sa_1",
			mockFn: func(m *mockProjectServiceAccountClientImpl) {
				m.RetrieveProjectServiceAccountFunc = func(projectID, serviceAccountID string) (*openaiorgs.ProjectServiceAccount, error) {
					if projectID != "proj_123" || serviceAccountID != "sa_1" {
						t.Errorf("unexpected params: projectID=%s, saID=%s", projectID, serviceAccountID)
					}
					sa := createMockProjectServiceAccount("sa_1", "My SA")
					return &sa, nil
				}
			},
			wantContains: []string{"Project Service Account details:", "sa_1", "My SA"},
		},
		{
			name:             "error from client",
			projectID:        "proj_123",
			serviceAccountID: "sa_1",
			mockFn: func(m *mockProjectServiceAccountClientImpl) {
				m.RetrieveProjectServiceAccountFunc = func(projectID, serviceAccountID string) (*openaiorgs.ProjectServiceAccount, error) {
					return nil, fmt.Errorf("not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectServiceAccountClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := retrieveProjectServiceAccountHandler(mock, tt.projectID, tt.serviceAccountID)
				if (err != nil) != tt.wantErr {
					t.Errorf("retrieveProjectServiceAccountHandler() error = %v, wantErr %v", err, tt.wantErr)
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

func TestDeleteProjectServiceAccountHandler(t *testing.T) {
	tests := []struct {
		name             string
		projectID        string
		serviceAccountID string
		mockFn           func(*mockProjectServiceAccountClientImpl)
		wantErr          bool
		wantContains     []string
	}{
		{
			name:             "successful delete",
			projectID:        "proj_123",
			serviceAccountID: "sa_1",
			mockFn: func(m *mockProjectServiceAccountClientImpl) {
				m.DeleteProjectServiceAccountFunc = func(projectID, serviceAccountID string) error {
					if projectID != "proj_123" || serviceAccountID != "sa_1" {
						t.Errorf("unexpected params: projectID=%s, saID=%s", projectID, serviceAccountID)
					}
					return nil
				}
			},
			wantContains: []string{"Project Service Account sa_1 deleted successfully"},
		},
		{
			name:             "error from client",
			projectID:        "proj_123",
			serviceAccountID: "sa_1",
			mockFn: func(m *mockProjectServiceAccountClientImpl) {
				m.DeleteProjectServiceAccountFunc = func(projectID, serviceAccountID string) error {
					return fmt.Errorf("delete failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectServiceAccountClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := deleteProjectServiceAccountHandler(mock, tt.projectID, tt.serviceAccountID)
				if (err != nil) != tt.wantErr {
					t.Errorf("deleteProjectServiceAccountHandler() error = %v, wantErr %v", err, tt.wantErr)
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
