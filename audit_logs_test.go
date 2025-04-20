package openaiorgs

import (
	"reflect"
	"testing"
	"time"
)

func TestListAuditLogs(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockAuditLogs := []AuditLog{
		{
			ID:          "log_123",
			Type:        "api_key.created",
			EffectiveAt: UnixSeconds(now),
			Actor: Actor{
				Type: "session",
				Session: &Session{
					User: AuditUser{
						ID:    "user_123",
						Email: "test@example.com",
					},
					IPAddress: "127.0.0.1",
					UserAgent: "test-agent",
				},
			},
			Details: &APIKeyCreated{
				ID: "key_123",
				Data: struct {
					Scopes []string `json:"scopes"`
				}{
					Scopes: []string{"read", "write"},
				},
			},
		},
	}

	// Register mock response
	response := ListResponse[AuditLog]{
		Object:  "list",
		Data:    mockAuditLogs,
		FirstID: "log_123",
		LastID:  "log_123",
		HasMore: false,
	}
	h.mockResponse("GET", AuditLogsListEndpoint, 200, response)

	// Make the API call
	auditLogs, err := h.client.ListAuditLogs(&AuditLogListParams{})
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(auditLogs.Data) != 1 {
		t.Errorf("Expected 1 audit log, got %d", len(auditLogs.Data))
		return
	}
	if mockAuditLogs[0].ID != auditLogs.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockAuditLogs[0].ID, auditLogs.Data[0].ID)
	}
	if mockAuditLogs[0].Type != auditLogs.Data[0].Type {
		t.Errorf("Expected Type %s, got %s", mockAuditLogs[0].Type, auditLogs.Data[0].Type)
	}

	// Verify the request was made
	h.assertRequest("GET", AuditLogsListEndpoint, 1)
}

func TestParseAuditLogPayload(t *testing.T) {
	testTime := time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		logType string
		payload map[string]any
		want    any
		wantErr bool
	}{
		// API Key events
		"api_key.created": {
			logType: "api_key.created",
			payload: map[string]any{
				"id": "key_123",
				"data": map[string]any{
					"scopes": []string{"read", "write"},
				},
			},
			want: &APIKeyCreated{
				ID: "key_123",
				Data: struct {
					Scopes []string `json:"scopes"`
				}{
					Scopes: []string{"read", "write"},
				},
			},
		},
		"api_key.updated": {
			logType: "api_key.updated",
			payload: map[string]any{
				"id": "key_123",
				"changes_requested": map[string]any{
					"scopes": []string{"read"},
				},
			},
			want: &APIKeyUpdated{
				ID: "key_123",
				ChangesRequested: struct {
					Scopes []string `json:"scopes"`
				}{
					Scopes: []string{"read"},
				},
			},
		},
		"api_key.deleted": {
			logType: "api_key.deleted",
			payload: map[string]any{"id": "key_123"},
			want:    &APIKeyDeleted{ID: "key_123"},
		},

		// Invite events
		"invite.sent": {
			logType: "invite.sent",
			payload: map[string]any{
				"id": "inv_123",
				"data": map[string]any{
					"email": "test@example.com",
				},
			},
			want: &InviteSent{
				ID: "inv_123",
				Data: struct {
					Email string `json:"email"`
				}{
					Email: "test@example.com",
				},
			},
		},
		"invite.accepted": {
			logType: "invite.accepted",
			payload: map[string]any{"id": "inv_123"},
			want:    &InviteAccepted{ID: "inv_123"},
		},
		"invite.deleted": {
			logType: "invite.deleted",
			payload: map[string]any{"id": "inv_123"},
			want:    &InviteDeleted{ID: "inv_123"},
		},

		// Login/Logout events
		"login.failed": {
			logType: "login.failed",
			payload: map[string]any{
				"error_code":    "invalid_credentials",
				"error_message": "Invalid email or password",
			},
			want: &LoginFailed{
				ErrorCode:    "invalid_credentials",
				ErrorMessage: "Invalid email or password",
			},
		},
		"login.succeeded": {
			logType: "login.succeeded",
			payload: map[string]any{
				"object":       "audit.event",
				"id":           "login_123",
				"type":         "login.succeeded",
				"effective_at": testTime.Unix(),
				"actor": map[string]any{
					"type": "session",
					"session": map[string]any{
						"user": map[string]any{
							"id":    "user_123",
							"email": "test@example.com",
						},
					},
				},
			},
			want: &LoginSucceeded{
				Object:      "audit.event",
				ID:          "login_123",
				Type:        "login.succeeded",
				EffectiveAt: testTime.Unix(),
				Actor: Actor{
					Type: "session",
					Session: &Session{
						User: AuditUser{
							ID:    "user_123",
							Email: "test@example.com",
						},
					},
				},
			},
		},

		// Organization events
		"organization.updated": {
			logType: "organization.updated",
			payload: map[string]any{
				"id": "org_123",
				"changes_requested": map[string]any{
					"name": "New Org Name",
				},
			},
			want: &OrganizationUpdated{
				ID: "org_123",
				ChangesRequested: struct {
					Name string `json:"name,omitempty"`
				}{
					Name: "New Org Name",
				},
			},
		},

		// Project events
		"project.created": {
			logType: "project.created",
			payload: map[string]any{
				"id": "proj_123",
				"data": map[string]any{
					"name":  "test-project",
					"title": "Test Project",
				},
			},
			want: &ProjectCreated{
				ID: "proj_123",
				Data: struct {
					Name  string `json:"name"`
					Title string `json:"title"`
				}{
					Name:  "test-project",
					Title: "Test Project",
				},
			},
		},
		"project.updated": {
			logType: "project.updated",
			payload: map[string]any{
				"id": "proj_123",
				"changes_requested": map[string]any{
					"title": "Updated Project",
				},
			},
			want: &ProjectUpdated{
				ID: "proj_123",
				ChangesRequested: struct {
					Title string `json:"title"`
				}{
					Title: "Updated Project",
				},
			},
		},
		"project.archived": {
			logType: "project.archived",
			payload: map[string]any{"id": "proj_123"},
			want:    &ProjectArchived{ID: "proj_123"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.want
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("parseAuditLogPayload() = %v, want %v", got, tc.want)
			}
		})
	}
}
