package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Mock client interface for testing admin API keys
type mockAdminAPIKeyClient interface {
	ListAdminAPIKeys(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.AdminAPIKey], error)
	CreateAdminAPIKey(name string, scopes []string) (*openaiorgs.AdminAPIKey, error)
	RetrieveAdminAPIKey(apiKeyID string) (*openaiorgs.AdminAPIKey, error)
	DeleteAdminAPIKey(apiKeyID string) error
}

// Mock implementation
type mockAdminAPIKeyClientImpl struct {
	ListAdminAPIKeysFunc    func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.AdminAPIKey], error)
	CreateAdminAPIKeyFunc   func(name string, scopes []string) (*openaiorgs.AdminAPIKey, error)
	RetrieveAdminAPIKeyFunc func(apiKeyID string) (*openaiorgs.AdminAPIKey, error)
	DeleteAdminAPIKeyFunc   func(apiKeyID string) error
}

func (m *mockAdminAPIKeyClientImpl) ListAdminAPIKeys(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.AdminAPIKey], error) {
	if m.ListAdminAPIKeysFunc != nil {
		return m.ListAdminAPIKeysFunc(limit, after)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockAdminAPIKeyClientImpl) CreateAdminAPIKey(name string, scopes []string) (*openaiorgs.AdminAPIKey, error) {
	if m.CreateAdminAPIKeyFunc != nil {
		return m.CreateAdminAPIKeyFunc(name, scopes)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockAdminAPIKeyClientImpl) RetrieveAdminAPIKey(apiKeyID string) (*openaiorgs.AdminAPIKey, error) {
	if m.RetrieveAdminAPIKeyFunc != nil {
		return m.RetrieveAdminAPIKeyFunc(apiKeyID)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockAdminAPIKeyClientImpl) DeleteAdminAPIKey(apiKeyID string) error {
	if m.DeleteAdminAPIKeyFunc != nil {
		return m.DeleteAdminAPIKeyFunc(apiKeyID)
	}
	return fmt.Errorf("not implemented")
}

// Testable handler functions

func listAdminAPIKeysHandler(client mockAdminAPIKeyClient, limit int, after string) error {
	apiKeys, err := client.ListAdminAPIKeys(limit, after)
	if err != nil {
		return wrapError("list admin API keys", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Redacted Value", "Created At", "Last Used At", "Scopes"},
		Rows:    make([][]string, len(apiKeys.Data)),
	}

	for i, key := range apiKeys.Data {
		data.Rows[i] = []string{
			key.ID,
			key.Name,
			key.RedactedValue,
			key.CreatedAt.String(),
			key.LastUsedAt.String(),
			strings.Join(key.Scopes, ", "),
		}
	}

	printTableData(data)
	return nil
}

func createAdminAPIKeyHandler(client mockAdminAPIKeyClient, name string, scopes []string) error {
	apiKey, err := client.CreateAdminAPIKey(name, scopes)
	if err != nil {
		return wrapError("create admin API key", err)
	}

	fmt.Printf("API Key created:\nID: %s\nName: %s\nRedacted Value: %s\nCreated At: %s\n",
		apiKey.ID, apiKey.Name, apiKey.RedactedValue, apiKey.CreatedAt.String())
	fmt.Printf("Scopes: %s\n", strings.Join(apiKey.Scopes, ", "))
	return nil
}

func retrieveAdminAPIKeyHandler(client mockAdminAPIKeyClient, id string) error {
	apiKey, err := client.RetrieveAdminAPIKey(id)
	if err != nil {
		return wrapError("retrieve admin API key", err)
	}

	fmt.Printf("API Key details:\nID: %s\nName: %s\nRedacted Value: %s\nCreated At: %s\n",
		apiKey.ID, apiKey.Name, apiKey.RedactedValue, apiKey.CreatedAt.String())
	fmt.Printf("Last Used At: %s\n", apiKey.LastUsedAt.String())
	fmt.Printf("Scopes: %s\n", strings.Join(apiKey.Scopes, ", "))
	return nil
}

func deleteAdminAPIKeyHandler(client mockAdminAPIKeyClient, id string) error {
	err := client.DeleteAdminAPIKey(id)
	if err != nil {
		return wrapError("delete admin API key", err)
	}

	fmt.Printf("API Key %s deleted successfully\n", id)
	return nil
}

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

// Tests

func TestListAdminAPIKeysHandler(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		after   string
		mockFn  func(*mockAdminAPIKeyClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:  "successful list with scopes",
			limit: 10,
			after: "",
			mockFn: func(m *mockAdminAPIKeyClientImpl) {
				m.ListAdminAPIKeysFunc = func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.AdminAPIKey], error) {
					return &openaiorgs.ListResponse[openaiorgs.AdminAPIKey]{
						Object: "list",
						Data: []openaiorgs.AdminAPIKey{
							createMockAdminAPIKey("key_123", "My API Key", "sk-...abc", []string{"api.read", "api.write"}),
						},
						FirstID: "key_123",
						LastID:  "key_123",
						HasMore: false,
					}, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "ID | Name | Redacted Value | Created At | Last Used At | Scopes") {
					t.Errorf("Expected table headers in output, got: %s", output)
				}
				if !strings.Contains(output, "key_123") {
					t.Errorf("Expected key_123 in output, got: %s", output)
				}
				if !strings.Contains(output, "api.read, api.write") {
					t.Errorf("Expected scopes in output, got: %s", output)
				}
			},
		},
		{
			name:  "error from client",
			limit: 10,
			after: "",
			mockFn: func(m *mockAdminAPIKeyClientImpl) {
				m.ListAdminAPIKeysFunc = func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.AdminAPIKey], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
		{
			name:  "empty list",
			limit: 10,
			after: "",
			mockFn: func(m *mockAdminAPIKeyClientImpl) {
				m.ListAdminAPIKeysFunc = func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.AdminAPIKey], error) {
					return &openaiorgs.ListResponse[openaiorgs.AdminAPIKey]{
						Object:  "list",
						Data:    []openaiorgs.AdminAPIKey{},
						HasMore: false,
					}, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "ID | Name | Redacted Value") {
					t.Errorf("Expected table headers even for empty list, got: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAdminAPIKeyClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listAdminAPIKeysHandler(mock, tt.limit, tt.after)
				if (err != nil) != tt.wantErr {
					t.Errorf("listAdminAPIKeysHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestCreateAdminAPIKeyHandler(t *testing.T) {
	tests := []struct {
		name    string
		keyName string
		scopes  []string
		mockFn  func(*mockAdminAPIKeyClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:    "successful create",
			keyName: "My New Key",
			scopes:  []string{"api.read", "api.write"},
			mockFn: func(m *mockAdminAPIKeyClientImpl) {
				m.CreateAdminAPIKeyFunc = func(name string, scopes []string) (*openaiorgs.AdminAPIKey, error) {
					if name != "My New Key" {
						t.Errorf("unexpected name: %s", name)
					}
					if len(scopes) != 2 || scopes[0] != "api.read" || scopes[1] != "api.write" {
						t.Errorf("unexpected scopes: %v", scopes)
					}
					key := createMockAdminAPIKey("key_123", "My New Key", "sk-...xyz", scopes)
					return &key, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "API Key created:") {
					t.Errorf("Expected 'API Key created:' in output, got: %s", output)
				}
				if !strings.Contains(output, "ID: key_123") {
					t.Errorf("Expected key ID in output, got: %s", output)
				}
				if !strings.Contains(output, "Scopes: api.read, api.write") {
					t.Errorf("Expected scopes in output, got: %s", output)
				}
			},
		},
		{
			name:    "error from client",
			keyName: "Bad Key",
			scopes:  []string{"api.read"},
			mockFn: func(m *mockAdminAPIKeyClientImpl) {
				m.CreateAdminAPIKeyFunc = func(name string, scopes []string) (*openaiorgs.AdminAPIKey, error) {
					return nil, fmt.Errorf("create failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAdminAPIKeyClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := createAdminAPIKeyHandler(mock, tt.keyName, tt.scopes)
				if (err != nil) != tt.wantErr {
					t.Errorf("createAdminAPIKeyHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestRetrieveAdminAPIKeyHandler(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*mockAdminAPIKeyClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name: "successful retrieve",
			id:   "key_123",
			mockFn: func(m *mockAdminAPIKeyClientImpl) {
				m.RetrieveAdminAPIKeyFunc = func(apiKeyID string) (*openaiorgs.AdminAPIKey, error) {
					if apiKeyID != "key_123" {
						t.Errorf("unexpected id: %s", apiKeyID)
					}
					key := createMockAdminAPIKey("key_123", "My Key", "sk-...abc", []string{"api.read"})
					return &key, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "API Key details:") {
					t.Errorf("Expected 'API Key details:' in output, got: %s", output)
				}
				if !strings.Contains(output, "ID: key_123") {
					t.Errorf("Expected key ID in output, got: %s", output)
				}
				if !strings.Contains(output, "Scopes: api.read") {
					t.Errorf("Expected scopes in output, got: %s", output)
				}
				if !strings.Contains(output, "Last Used At:") {
					t.Errorf("Expected 'Last Used At:' in output, got: %s", output)
				}
			},
		},
		{
			name: "error from client",
			id:   "key_999",
			mockFn: func(m *mockAdminAPIKeyClientImpl) {
				m.RetrieveAdminAPIKeyFunc = func(apiKeyID string) (*openaiorgs.AdminAPIKey, error) {
					return nil, fmt.Errorf("key not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAdminAPIKeyClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := retrieveAdminAPIKeyHandler(mock, tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("retrieveAdminAPIKeyHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestDeleteAdminAPIKeyHandler(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*mockAdminAPIKeyClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name: "successful delete",
			id:   "key_123",
			mockFn: func(m *mockAdminAPIKeyClientImpl) {
				m.DeleteAdminAPIKeyFunc = func(apiKeyID string) error {
					if apiKeyID != "key_123" {
						t.Errorf("unexpected id: %s", apiKeyID)
					}
					return nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "API Key key_123 deleted successfully") {
					t.Errorf("Expected delete success message, got: %s", output)
				}
			},
		},
		{
			name: "error from client",
			id:   "key_999",
			mockFn: func(m *mockAdminAPIKeyClientImpl) {
				m.DeleteAdminAPIKeyFunc = func(apiKeyID string) error {
					return fmt.Errorf("delete failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAdminAPIKeyClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := deleteAdminAPIKeyHandler(mock, tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("deleteAdminAPIKeyHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}
