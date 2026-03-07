package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Mock client interface for testing project API keys
type mockProjectAPIKeyClient interface {
	ListProjectApiKeys(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectApiKey], error)
	RetrieveProjectApiKey(projectID, apiKeyID string) (*openaiorgs.ProjectApiKey, error)
	DeleteProjectApiKey(projectID, apiKeyID string) error
}

// Mock implementation
type mockProjectAPIKeyClientImpl struct {
	ListProjectApiKeysFunc    func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectApiKey], error)
	RetrieveProjectApiKeyFunc func(projectID, apiKeyID string) (*openaiorgs.ProjectApiKey, error)
	DeleteProjectApiKeyFunc   func(projectID, apiKeyID string) error
}

func (m *mockProjectAPIKeyClientImpl) ListProjectApiKeys(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectApiKey], error) {
	if m.ListProjectApiKeysFunc != nil {
		return m.ListProjectApiKeysFunc(projectID, limit, after)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectAPIKeyClientImpl) RetrieveProjectApiKey(projectID, apiKeyID string) (*openaiorgs.ProjectApiKey, error) {
	if m.RetrieveProjectApiKeyFunc != nil {
		return m.RetrieveProjectApiKeyFunc(projectID, apiKeyID)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectAPIKeyClientImpl) DeleteProjectApiKey(projectID, apiKeyID string) error {
	if m.DeleteProjectApiKeyFunc != nil {
		return m.DeleteProjectApiKeyFunc(projectID, apiKeyID)
	}
	return fmt.Errorf("not implemented")
}

// Testable handlers

func listProjectAPIKeysTableHandler(client mockProjectAPIKeyClient, projectID string, limit int, after string) error {
	apiKeys, err := client.ListProjectApiKeys(projectID, limit, after)
	if err != nil {
		return fmt.Errorf("failed to list project API keys: %w", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Created At", "Owner"},
		Rows:    make([][]string, len(apiKeys.Data)),
	}
	for i, key := range apiKeys.Data {
		data.Rows[i] = []string{
			key.ID,
			key.Name,
			key.CreatedAt.String(),
			key.Owner.String(),
		}
	}
	printTableData(data)
	return nil
}

func listProjectAPIKeysJSONHandler(client mockProjectAPIKeyClient, projectID string, limit int, after string) error {
	apiKeys, err := client.ListProjectApiKeys(projectID, limit, after)
	if err != nil {
		return fmt.Errorf("failed to list project API keys: %w", err)
	}

	jsonData, err := json.Marshal(apiKeys)
	if err != nil {
		return fmt.Errorf("failed to marshal API keys: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

func retrieveProjectAPIKeyTableHandler(client mockProjectAPIKeyClient, projectID, apiKeyID string) error {
	apiKey, err := client.RetrieveProjectApiKey(projectID, apiKeyID)
	if err != nil {
		return fmt.Errorf("failed to retrieve project API key: %w", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Created At", "Owner"},
		Rows: [][]string{{
			apiKey.ID,
			apiKey.Name,
			apiKey.CreatedAt.String(),
			apiKey.Owner.String(),
		}},
	}
	printTableData(data)
	return nil
}

func deleteProjectAPIKeyHandler(client mockProjectAPIKeyClient, projectID, apiKeyID string) error {
	err := client.DeleteProjectApiKey(projectID, apiKeyID)
	if err != nil {
		return fmt.Errorf("failed to delete project API key: %w", err)
	}

	fmt.Printf("Successfully deleted project API key %s\n", apiKeyID)
	return nil
}

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

// Tests

func TestListProjectAPIKeysTableHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		limit        int
		after        string
		mockFn       func(*mockProjectAPIKeyClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful list pretty",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectAPIKeyClientImpl) {
				key := createMockProjectApiKey("key_1", "My API Key")
				m.ListProjectApiKeysFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectApiKey], error) {
					if projectID != "proj_123" {
						t.Errorf("unexpected projectID: %s", projectID)
					}
					return &openaiorgs.ListResponse[openaiorgs.ProjectApiKey]{
						Object:  "list",
						Data:    []openaiorgs.ProjectApiKey{key},
						FirstID: "key_1",
						LastID:  "key_1",
						HasMore: false,
					}, nil
				}
			},
			wantContains: []string{"ID | Name | Created At | Owner", "key_1", "My API Key"},
		},
		{
			name:      "empty list",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectAPIKeyClientImpl) {
				m.ListProjectApiKeysFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectApiKey], error) {
					return &openaiorgs.ListResponse[openaiorgs.ProjectApiKey]{
						Object: "list",
						Data:   []openaiorgs.ProjectApiKey{},
					}, nil
				}
			},
			wantContains: []string{"ID | Name | Created At | Owner"},
		},
		{
			name:      "error from client",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectAPIKeyClientImpl) {
				m.ListProjectApiKeysFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectApiKey], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectAPIKeyClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listProjectAPIKeysTableHandler(mock, tt.projectID, tt.limit, tt.after)
				if (err != nil) != tt.wantErr {
					t.Errorf("listProjectAPIKeysTableHandler() error = %v, wantErr %v", err, tt.wantErr)
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

func TestListProjectAPIKeysJSONHandler(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		limit     int
		after     string
		mockFn    func(*mockProjectAPIKeyClientImpl)
		wantErr   bool
	}{
		{
			name:      "successful list json",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectAPIKeyClientImpl) {
				key := createMockProjectApiKey("key_1", "My API Key")
				m.ListProjectApiKeysFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectApiKey], error) {
					return &openaiorgs.ListResponse[openaiorgs.ProjectApiKey]{
						Object:  "list",
						Data:    []openaiorgs.ProjectApiKey{key},
						FirstID: "key_1",
						LastID:  "key_1",
						HasMore: false,
					}, nil
				}
			},
		},
		{
			name:      "error from client json",
			projectID: "proj_123",
			limit:     10,
			after:     "",
			mockFn: func(m *mockProjectAPIKeyClientImpl) {
				m.ListProjectApiKeysFunc = func(projectID string, limit int, after string) (*openaiorgs.ListResponse[openaiorgs.ProjectApiKey], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectAPIKeyClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listProjectAPIKeysJSONHandler(mock, tt.projectID, tt.limit, tt.after)
				if (err != nil) != tt.wantErr {
					t.Errorf("listProjectAPIKeysJSONHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr {
				// Verify it's valid JSON
				var result map[string]interface{}
				if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &result); err != nil {
					t.Errorf("Expected valid JSON output, got: %s", output)
				}
				if !strings.Contains(output, "key_1") {
					t.Errorf("Expected output to contain key_1, got: %s", output)
				}
			}
		})
	}
}

func TestRetrieveProjectAPIKeyTableHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		apiKeyID     string
		mockFn       func(*mockProjectAPIKeyClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful retrieve",
			projectID: "proj_123",
			apiKeyID:  "key_1",
			mockFn: func(m *mockProjectAPIKeyClientImpl) {
				m.RetrieveProjectApiKeyFunc = func(projectID, apiKeyID string) (*openaiorgs.ProjectApiKey, error) {
					if projectID != "proj_123" || apiKeyID != "key_1" {
						t.Errorf("unexpected params: projectID=%s, apiKeyID=%s", projectID, apiKeyID)
					}
					key := createMockProjectApiKey("key_1", "My API Key")
					return &key, nil
				}
			},
			wantContains: []string{"ID | Name | Created At | Owner", "key_1", "My API Key"},
		},
		{
			name:      "error from client",
			projectID: "proj_123",
			apiKeyID:  "key_1",
			mockFn: func(m *mockProjectAPIKeyClientImpl) {
				m.RetrieveProjectApiKeyFunc = func(projectID, apiKeyID string) (*openaiorgs.ProjectApiKey, error) {
					return nil, fmt.Errorf("not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectAPIKeyClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := retrieveProjectAPIKeyTableHandler(mock, tt.projectID, tt.apiKeyID)
				if (err != nil) != tt.wantErr {
					t.Errorf("retrieveProjectAPIKeyTableHandler() error = %v, wantErr %v", err, tt.wantErr)
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

func TestDeleteProjectAPIKeyHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		apiKeyID     string
		mockFn       func(*mockProjectAPIKeyClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful delete",
			projectID: "proj_123",
			apiKeyID:  "key_1",
			mockFn: func(m *mockProjectAPIKeyClientImpl) {
				m.DeleteProjectApiKeyFunc = func(projectID, apiKeyID string) error {
					if projectID != "proj_123" || apiKeyID != "key_1" {
						t.Errorf("unexpected params: projectID=%s, apiKeyID=%s", projectID, apiKeyID)
					}
					return nil
				}
			},
			wantContains: []string{"Successfully deleted project API key key_1"},
		},
		{
			name:      "error from client",
			projectID: "proj_123",
			apiKeyID:  "key_1",
			mockFn: func(m *mockProjectAPIKeyClientImpl) {
				m.DeleteProjectApiKeyFunc = func(projectID, apiKeyID string) error {
					return fmt.Errorf("delete failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectAPIKeyClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := deleteProjectAPIKeyHandler(mock, tt.projectID, tt.apiKeyID)
				if (err != nil) != tt.wantErr {
					t.Errorf("deleteProjectAPIKeyHandler() error = %v, wantErr %v", err, tt.wantErr)
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
