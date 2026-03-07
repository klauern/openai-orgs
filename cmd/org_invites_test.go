package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Mock client interface for testing invites
type mockInviteClient interface {
	ListInvites(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.Invite], error)
	CreateInvite(email, role string) (*openaiorgs.Invite, error)
	RetrieveInvite(id string) (*openaiorgs.Invite, error)
	DeleteInvite(id string) error
}

// Mock implementation
type mockInviteClientImpl struct {
	ListInvitesFunc    func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.Invite], error)
	CreateInviteFunc   func(email, role string) (*openaiorgs.Invite, error)
	RetrieveInviteFunc func(id string) (*openaiorgs.Invite, error)
	DeleteInviteFunc   func(id string) error
}

func (m *mockInviteClientImpl) ListInvites(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.Invite], error) {
	if m.ListInvitesFunc != nil {
		return m.ListInvitesFunc(limit, after)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockInviteClientImpl) CreateInvite(email, role string) (*openaiorgs.Invite, error) {
	if m.CreateInviteFunc != nil {
		return m.CreateInviteFunc(email, role)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockInviteClientImpl) RetrieveInvite(id string) (*openaiorgs.Invite, error) {
	if m.RetrieveInviteFunc != nil {
		return m.RetrieveInviteFunc(id)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockInviteClientImpl) DeleteInvite(id string) error {
	if m.DeleteInviteFunc != nil {
		return m.DeleteInviteFunc(id)
	}
	return fmt.Errorf("not implemented")
}

// Testable handler functions

func listInvitesHandler(client mockInviteClient, limit int, after string) error {
	resp, err := client.ListInvites(limit, after)
	if err != nil {
		return wrapError("list invites", err)
	}

	data := TableData{
		Headers: []string{"ID", "Email", "Role", "Status", "Created At", "Expires At", "Accepted At"},
		Rows:    make([][]string, len(resp.Data)),
	}

	for i, invite := range resp.Data {
		acceptedAt := "N/A"
		if invite.AcceptedAt != nil {
			acceptedAt = invite.AcceptedAt.String()
		}
		data.Rows[i] = []string{
			invite.ID,
			invite.Email,
			invite.Role,
			invite.Status,
			invite.CreatedAt.String(),
			invite.ExpiresAt.String(),
			acceptedAt,
		}
	}

	printTableData(data)
	return nil
}

func createInviteHandler(client mockInviteClient, email, role string) error {
	invite, err := client.CreateInvite(email, role)
	if err != nil {
		return wrapError("create invite", err)
	}

	fmt.Printf("Invite created:\nID: %s\nEmail: %s\nRole: %s\nCreated At: %s\nExpires At: %s\n",
		invite.ID, invite.Email, invite.Role, invite.CreatedAt.String(), invite.ExpiresAt.String())
	return nil
}

func retrieveInviteHandler(client mockInviteClient, id string) error {
	invite, err := client.RetrieveInvite(id)
	if err != nil {
		return wrapError("retrieve invite", err)
	}

	fmt.Printf("Invite details:\nID: %s\nEmail: %s\nRole: %s\nCreated At: %s\nExpires At: %s\n",
		invite.ID, invite.Email, invite.Role, invite.CreatedAt.String(), invite.ExpiresAt.String())
	return nil
}

func deleteInviteHandler(client mockInviteClient, id string) error {
	err := client.DeleteInvite(id)
	if err != nil {
		return wrapError("delete invite", err)
	}

	fmt.Printf("Invite %s deleted successfully\n", id)
	return nil
}

// Helper to create a mock invite
func createMockInvite(id, email, role, status string, acceptedAt *openaiorgs.UnixSeconds) openaiorgs.Invite {
	now := time.Now()
	return openaiorgs.Invite{
		ObjectType: "organization.invite",
		ID:         id,
		Email:      email,
		Role:       role,
		Status:     status,
		ExpiresAt:  openaiorgs.UnixSeconds(now.Add(7 * 24 * time.Hour)),
		AcceptedAt: acceptedAt,
		CreatedAt:  openaiorgs.UnixSeconds(now),
	}
}

// Tests

func TestListInvitesHandler(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		after   string
		mockFn  func(*mockInviteClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:  "successful list with nil AcceptedAt",
			limit: 10,
			after: "",
			mockFn: func(m *mockInviteClientImpl) {
				m.ListInvitesFunc = func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.Invite], error) {
					return &openaiorgs.ListResponse[openaiorgs.Invite]{
						Object: "list",
						Data: []openaiorgs.Invite{
							createMockInvite("inv_123", "alice@example.com", "member", "pending", nil),
						},
						FirstID: "inv_123",
						LastID:  "inv_123",
						HasMore: false,
					}, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "ID | Email | Role | Status | Created At | Expires At | Accepted At") {
					t.Errorf("Expected table headers in output, got: %s", output)
				}
				if !strings.Contains(output, "N/A") {
					t.Errorf("Expected N/A for nil AcceptedAt, got: %s", output)
				}
				if !strings.Contains(output, "inv_123") {
					t.Errorf("Expected inv_123 in output, got: %s", output)
				}
			},
		},
		{
			name:  "successful list with AcceptedAt set",
			limit: 10,
			after: "",
			mockFn: func(m *mockInviteClientImpl) {
				m.ListInvitesFunc = func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.Invite], error) {
					accepted := openaiorgs.UnixSeconds(time.Now())
					return &openaiorgs.ListResponse[openaiorgs.Invite]{
						Object: "list",
						Data: []openaiorgs.Invite{
							createMockInvite("inv_456", "bob@example.com", "owner", "accepted", &accepted),
						},
						FirstID: "inv_456",
						LastID:  "inv_456",
						HasMore: false,
					}, nil
				}
			},
			check: func(t *testing.T, output string) {
				if strings.Contains(output, "N/A") {
					t.Errorf("Expected no N/A when AcceptedAt is set, got: %s", output)
				}
			},
		},
		{
			name:  "error from client",
			limit: 10,
			after: "",
			mockFn: func(m *mockInviteClientImpl) {
				m.ListInvitesFunc = func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.Invite], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
		{
			name:  "empty list",
			limit: 10,
			after: "",
			mockFn: func(m *mockInviteClientImpl) {
				m.ListInvitesFunc = func(limit int, after string) (*openaiorgs.ListResponse[openaiorgs.Invite], error) {
					return &openaiorgs.ListResponse[openaiorgs.Invite]{
						Object:  "list",
						Data:    []openaiorgs.Invite{},
						HasMore: false,
					}, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "ID | Email | Role | Status") {
					t.Errorf("Expected table headers even for empty list, got: %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockInviteClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listInvitesHandler(mock, tt.limit, tt.after)
				if (err != nil) != tt.wantErr {
					t.Errorf("listInvitesHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestCreateInviteHandler(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		role    string
		mockFn  func(*mockInviteClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name:  "successful create",
			email: "alice@example.com",
			role:  "member",
			mockFn: func(m *mockInviteClientImpl) {
				m.CreateInviteFunc = func(email, role string) (*openaiorgs.Invite, error) {
					if email != "alice@example.com" || role != "member" {
						t.Errorf("unexpected params: email=%s, role=%s", email, role)
					}
					invite := createMockInvite("inv_123", "alice@example.com", "member", "pending", nil)
					return &invite, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "Invite created:") {
					t.Errorf("Expected 'Invite created:' in output, got: %s", output)
				}
				if !strings.Contains(output, "ID: inv_123") {
					t.Errorf("Expected invite ID in output, got: %s", output)
				}
				if !strings.Contains(output, "Email: alice@example.com") {
					t.Errorf("Expected email in output, got: %s", output)
				}
			},
		},
		{
			name:  "error from client",
			email: "bad@example.com",
			role:  "member",
			mockFn: func(m *mockInviteClientImpl) {
				m.CreateInviteFunc = func(email, role string) (*openaiorgs.Invite, error) {
					return nil, fmt.Errorf("create failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockInviteClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := createInviteHandler(mock, tt.email, tt.role)
				if (err != nil) != tt.wantErr {
					t.Errorf("createInviteHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestRetrieveInviteHandler(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*mockInviteClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name: "successful retrieve",
			id:   "inv_123",
			mockFn: func(m *mockInviteClientImpl) {
				m.RetrieveInviteFunc = func(id string) (*openaiorgs.Invite, error) {
					if id != "inv_123" {
						t.Errorf("unexpected id: %s", id)
					}
					invite := createMockInvite("inv_123", "alice@example.com", "member", "pending", nil)
					return &invite, nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "Invite details:") {
					t.Errorf("Expected 'Invite details:' in output, got: %s", output)
				}
				if !strings.Contains(output, "ID: inv_123") {
					t.Errorf("Expected invite ID in output, got: %s", output)
				}
			},
		},
		{
			name: "error from client",
			id:   "inv_999",
			mockFn: func(m *mockInviteClientImpl) {
				m.RetrieveInviteFunc = func(id string) (*openaiorgs.Invite, error) {
					return nil, fmt.Errorf("invite not found")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockInviteClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := retrieveInviteHandler(mock, tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("retrieveInviteHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}

func TestDeleteInviteHandler(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(*mockInviteClientImpl)
		wantErr bool
		check   func(t *testing.T, output string)
	}{
		{
			name: "successful delete",
			id:   "inv_123",
			mockFn: func(m *mockInviteClientImpl) {
				m.DeleteInviteFunc = func(id string) error {
					if id != "inv_123" {
						t.Errorf("unexpected id: %s", id)
					}
					return nil
				}
			},
			check: func(t *testing.T, output string) {
				if !strings.Contains(output, "Invite inv_123 deleted successfully") {
					t.Errorf("Expected delete success message, got: %s", output)
				}
			},
		},
		{
			name: "error from client",
			id:   "inv_999",
			mockFn: func(m *mockInviteClientImpl) {
				m.DeleteInviteFunc = func(id string) error {
					return fmt.Errorf("delete failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockInviteClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := deleteInviteHandler(mock, tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("deleteInviteHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && tt.check != nil {
				tt.check(t, output)
			}
		})
	}
}
