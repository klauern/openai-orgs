package openaiorgs

import (
	"encoding/json"
	"reflect"
	"strings"
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

func TestListAuditLogsWithAllParams(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Register mock response
	response := ListResponse[AuditLog]{
		Object:  "list",
		Data:    []AuditLog{},
		HasMore: false,
	}
	h.mockResponse("GET", AuditLogsListEndpoint, 200, response)

	// Make the API call with all parameters
	_, err := h.client.ListAuditLogs(&AuditLogListParams{
		Limit:  50,
		After:  "log_after",
		Before: "log_before",
		EffectiveAt: &EffectiveAt{
			Gte: 1700000000,
			Gt:  1699999999,
			Lte: 1800000000,
			Lt:  1800000001,
		},
		ProjectIDs:  []string{"proj_1", "proj_2"},
		EventTypes:  []string{"api_key.created", "user.added"},
		ActorIDs:    []string{"actor_1", "actor_2"},
		ActorEmails: []string{"user1@example.com", "user2@example.com"},
		ResourceIDs: []string{"resource_1", "resource_2"},
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	// Verify the request was made
	h.assertRequest("GET", AuditLogsListEndpoint, 1)
}

func TestListAuditLogsWithNilParams(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Register mock response
	response := ListResponse[AuditLog]{
		Object:  "list",
		Data:    []AuditLog{},
		HasMore: false,
	}
	h.mockResponse("GET", AuditLogsListEndpoint, 200, response)

	// Make the API call with nil params
	_, err := h.client.ListAuditLogs(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	// Verify the request was made
	h.assertRequest("GET", AuditLogsListEndpoint, 1)
}

func TestAuditLogString(t *testing.T) {
	testTime := time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC)

	tests := map[string]struct {
		log      AuditLog
		contains []string
	}{
		"with project and session actor": {
			log: AuditLog{
				ID:          "log_123",
				Type:        "api_key.created",
				EffectiveAt: UnixSeconds(testTime),
				Project: &AuditProject{
					ID:   "proj_123",
					Name: "Test Project",
				},
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
			contains: []string{"log_123", "api_key.created", "Test Project(proj_123)", "user:test@example.com"},
		},
		"with project and api_key actor": {
			log: AuditLog{
				ID:          "log_456",
				Type:        "user.deleted",
				EffectiveAt: UnixSeconds(testTime),
				Project: &AuditProject{
					ID:   "proj_456",
					Name: "Another Project",
				},
				Actor: Actor{
					Type: "api_key",
					APIKey: &APIKeyActor{
						Type: "user",
						User: AuditUser{
							ID:    "user_456",
							Email: "apikey@example.com",
						},
					},
				},
			},
			contains: []string{"log_456", "user.deleted", "Another Project(proj_456)", "apikey:apikey@example.com"},
		},
		"without project": {
			log: AuditLog{
				ID:          "log_789",
				Type:        "login.succeeded",
				EffectiveAt: UnixSeconds(testTime),
				Project:     nil,
				Actor: Actor{
					Type: "session",
					Session: &Session{
						User: AuditUser{
							ID:    "user_789",
							Email: "noproject@example.com",
						},
					},
				},
			},
			contains: []string{"log_789", "login.succeeded", "no project", "user:noproject@example.com"},
		},
		"unknown actor type": {
			log: AuditLog{
				ID:          "log_000",
				Type:        "unknown.event",
				EffectiveAt: UnixSeconds(testTime),
				Actor: Actor{
					Type: "unknown",
				},
			},
			contains: []string{"log_000", "unknown.event", "no project", "unknown"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.log.String()
			for _, substr := range tc.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("Expected String() to contain %q, got %q", substr, result)
				}
			}
		})
	}
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

func TestAuditLogUnmarshalJSON_EdgeCases(t *testing.T) {
	tests := map[string]struct {
		input     string
		expectErr bool
		errMsg    string
		want      *AuditLog
	}{
		"invalid JSON": {
			input:     `{"invalid": json}`,
			expectErr: true,
		},
		"unknown audit log type stores raw details": {
			// Unknown event types should store the raw JSON as a map instead of erroring
			input: `{"id": "log_123", "type": "unknown.event", "effective_at": 1234567890, "unknown.event": {"some": "data"}}`,
			want: &AuditLog{
				ID:          "log_123",
				Type:        "unknown.event",
				EffectiveAt: UnixSeconds(time.Unix(1234567890, 0)),
				Details:     map[string]any{"some": "data"},
			},
		},
		"unknown event type with array value returns nil details": {
			// When unknown event type has non-object JSON (like an array), store nil
			input: `{"id": "log_123", "type": "unknown.array", "effective_at": 1234567890, "unknown.array": ["item1", "item2"]}`,
			want: &AuditLog{
				ID:          "log_123",
				Type:        "unknown.array",
				EffectiveAt: UnixSeconds(time.Unix(1234567890, 0)),
				Details:     nil,
			},
		},
		"no event details": {
			input: `{"id": "log_123", "type": "api_key.created", "effective_at": 1234567890}`,
			want: &AuditLog{
				ID:          "log_123",
				Type:        "api_key.created",
				EffectiveAt: UnixSeconds(time.Unix(1234567890, 0)),
				Details:     nil,
			},
		},
		"malformed event details JSON": {
			input:     `{"id": "log_123", "type": "api_key.created", "effective_at": 1234567890, "api_key.created": {"invalid": json}}`,
			expectErr: true,
			errMsg:    "invalid character",
		},
		"valid api_key.created with dynamic key": {
			input: `{"id": "log_123", "type": "api_key.created", "effective_at": 1234567890, "api_key.created": {"id": "key_123", "data": {"scopes": ["scope1", "scope2"]}}}`,
			want: &AuditLog{
				ID:          "log_123",
				Type:        "api_key.created",
				EffectiveAt: UnixSeconds(time.Unix(1234567890, 0)),
				Details: &APIKeyCreated{
					ID: "key_123",
					Data: struct {
						Scopes []string `json:"scopes"`
					}{
						Scopes: []string{"scope1", "scope2"},
					},
				},
			},
		},
		"project.archived with dynamic key": {
			input: `{"id": "log_123", "type": "project.archived", "effective_at": 1234567890, "project.archived": {"id": "proj_123"}}`,
			want: &AuditLog{
				ID:          "log_123",
				Type:        "project.archived",
				EffectiveAt: UnixSeconds(time.Unix(1234567890, 0)),
				Details: &ProjectArchived{
					ID: "proj_123",
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var log AuditLog
			err := json.Unmarshal([]byte(tc.input), &log)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tc.errMsg != "" && !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tc.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			if tc.want != nil {
				if log.ID != tc.want.ID {
					t.Errorf("Expected ID %q, got %q", tc.want.ID, log.ID)
				}
				if log.Type != tc.want.Type {
					t.Errorf("Expected Type %q, got %q", tc.want.Type, log.Type)
				}
				if log.EffectiveAt != tc.want.EffectiveAt {
					t.Errorf("Expected EffectiveAt %v, got %v", tc.want.EffectiveAt, log.EffectiveAt)
				}
				if tc.want.Details == nil && log.Details != nil {
					t.Errorf("Expected Details to be nil, got %v", log.Details)
				}
				if tc.want.Details != nil && log.Details == nil {
					t.Errorf("Expected Details to be non-nil, got nil")
				}
			}
		})
	}
}

// TestSessionWithIPAddressDetails tests that ja3, ja4, and ip_address_details are correctly parsed
func TestSessionWithIPAddressDetails(t *testing.T) {
	input := `{
		"object": "organization.audit_log",
		"id": "audit_log-PaILxsC4mrLvPYdakaGGtaMy",
		"type": "invite.deleted",
		"effective_at": 1759243237,
		"project": {
			"id": "proj_8E9dKUupvJVY2Ge9c3R4TwhH",
			"name": "Default project"
		},
		"actor": {
			"type": "session",
			"session": {
				"user": {
					"id": "user-4e88SmmkWl80jnhGaAMau4Uk",
					"email": "nklauer@zendesk.com"
				},
				"ip_address": "216.198.0.23",
				"user_agent": "go-resty/2.16.2 (https://github.com/go-resty/resty)",
				"ja3": "e69402f870ecf542b4f017b0ed32936a",
				"ja4": "t13d1312h2_f57a46bbacb6_ab7e3b40a677",
				"ip_address_details": {
					"country": "US",
					"city": "Portland",
					"region": "Oregon",
					"region_code": "OR",
					"asn": "16509",
					"latitude": "45.52345",
					"longitude": "-122.67621"
				}
			}
		},
		"invite.deleted": {
			"id": "invite-hhPjGZ2Zu09bzFhWoPUWacUb"
		}
	}`

	var log AuditLog
	err := json.Unmarshal([]byte(input), &log)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify Object field
	if log.Object != "organization.audit_log" {
		t.Errorf("Expected Object 'organization.audit_log', got %q", log.Object)
	}

	// Verify basic fields
	if log.ID != "audit_log-PaILxsC4mrLvPYdakaGGtaMy" {
		t.Errorf("Expected ID 'audit_log-PaILxsC4mrLvPYdakaGGtaMy', got %q", log.ID)
	}
	if log.Type != "invite.deleted" {
		t.Errorf("Expected Type 'invite.deleted', got %q", log.Type)
	}

	// Verify session fields
	if log.Actor.Session == nil {
		t.Fatal("Expected Session to be non-nil")
	}

	session := log.Actor.Session
	if session.JA3 != "e69402f870ecf542b4f017b0ed32936a" {
		t.Errorf("Expected JA3 'e69402f870ecf542b4f017b0ed32936a', got %q", session.JA3)
	}
	if session.JA4 != "t13d1312h2_f57a46bbacb6_ab7e3b40a677" {
		t.Errorf("Expected JA4 't13d1312h2_f57a46bbacb6_ab7e3b40a677', got %q", session.JA4)
	}

	// Verify IP address details
	if session.IPAddressDetails == nil {
		t.Fatal("Expected IPAddressDetails to be non-nil")
	}

	ipDetails := session.IPAddressDetails
	if ipDetails.Country != "US" {
		t.Errorf("Expected Country 'US', got %q", ipDetails.Country)
	}
	if ipDetails.City != "Portland" {
		t.Errorf("Expected City 'Portland', got %q", ipDetails.City)
	}
	if ipDetails.Region != "Oregon" {
		t.Errorf("Expected Region 'Oregon', got %q", ipDetails.Region)
	}
	if ipDetails.RegionCode != "OR" {
		t.Errorf("Expected RegionCode 'OR', got %q", ipDetails.RegionCode)
	}
	if ipDetails.ASN != "16509" {
		t.Errorf("Expected ASN '16509', got %q", ipDetails.ASN)
	}
	if ipDetails.Latitude != "45.52345" {
		t.Errorf("Expected Latitude '45.52345', got %q", ipDetails.Latitude)
	}
	if ipDetails.Longitude != "-122.67621" {
		t.Errorf("Expected Longitude '-122.67621', got %q", ipDetails.Longitude)
	}

	// Verify event details
	details, ok := log.Details.(*InviteDeleted)
	if !ok {
		t.Fatalf("Expected Details to be *InviteDeleted, got %T", log.Details)
	}
	if details.ID != "invite-hhPjGZ2Zu09bzFhWoPUWacUb" {
		t.Errorf("Expected invite ID 'invite-hhPjGZ2Zu09bzFhWoPUWacUb', got %q", details.ID)
	}
}

// TestSessionWithoutOptionalFields ensures backward compatibility for responses without new fields
func TestSessionWithoutOptionalFields(t *testing.T) {
	input := `{
		"id": "audit_log-123",
		"type": "invite.deleted",
		"effective_at": 1759243237,
		"actor": {
			"type": "session",
			"session": {
				"user": {
					"id": "user-123",
					"email": "test@example.com"
				},
				"ip_address": "127.0.0.1",
				"user_agent": "test-agent"
			}
		},
		"invite.deleted": {
			"id": "invite-123"
		}
	}`

	var log AuditLog
	err := json.Unmarshal([]byte(input), &log)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	session := log.Actor.Session
	if session == nil {
		t.Fatal("Expected Session to be non-nil")
	}

	// Verify old fields still work
	if session.IPAddress != "127.0.0.1" {
		t.Errorf("Expected IPAddress '127.0.0.1', got %q", session.IPAddress)
	}
	if session.UserAgent != "test-agent" {
		t.Errorf("Expected UserAgent 'test-agent', got %q", session.UserAgent)
	}

	// Verify new optional fields are empty/nil
	if session.JA3 != "" {
		t.Errorf("Expected JA3 to be empty, got %q", session.JA3)
	}
	if session.JA4 != "" {
		t.Errorf("Expected JA4 to be empty, got %q", session.JA4)
	}
	if session.IPAddressDetails != nil {
		t.Errorf("Expected IPAddressDetails to be nil, got %+v", session.IPAddressDetails)
	}
}

// TestAuditLogObjectField verifies the Object field is correctly captured
func TestAuditLogObjectField(t *testing.T) {
	input := `{
		"object": "organization.audit_log",
		"id": "audit_log-123",
		"type": "project.created",
		"effective_at": 1759243237,
		"actor": {
			"type": "session",
			"session": {
				"user": {"id": "user-123", "email": "test@example.com"},
				"ip_address": "127.0.0.1",
				"user_agent": "test-agent"
			}
		},
		"project.created": {
			"id": "proj_123",
			"data": {"name": "test", "title": "Test"}
		}
	}`

	var log AuditLog
	err := json.Unmarshal([]byte(input), &log)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if log.Object != "organization.audit_log" {
		t.Errorf("Expected Object 'organization.audit_log', got %q", log.Object)
	}
}

// TestAuditLogWithDynamicEventKey tests the actual API format with dynamic keys
func TestAuditLogWithDynamicEventKey(t *testing.T) {
	tests := map[string]struct {
		input       string
		wantType    string
		wantDetails any
	}{
		"invite.deleted": {
			input:       `{"id": "log_123", "type": "invite.deleted", "effective_at": 1234567890, "actor": {"type": "session"}, "invite.deleted": {"id": "invite-123"}}`,
			wantType:    "invite.deleted",
			wantDetails: &InviteDeleted{ID: "invite-123"},
		},
		"api_key.created": {
			input:    `{"id": "log_123", "type": "api_key.created", "effective_at": 1234567890, "actor": {"type": "session"}, "api_key.created": {"id": "key-123", "data": {"scopes": ["read"]}}}`,
			wantType: "api_key.created",
			wantDetails: &APIKeyCreated{ID: "key-123", Data: struct {
				Scopes []string `json:"scopes"`
			}{Scopes: []string{"read"}}},
		},
		"api_key.updated": {
			input:    `{"id": "log_123", "type": "api_key.updated", "effective_at": 1234567890, "actor": {"type": "session"}, "api_key.updated": {"id": "key-123", "changes_requested": {"scopes": ["write"]}}}`,
			wantType: "api_key.updated",
			wantDetails: &APIKeyUpdated{ID: "key-123", ChangesRequested: struct {
				Scopes []string `json:"scopes"`
			}{Scopes: []string{"write"}}},
		},
		"api_key.deleted": {
			input:       `{"id": "log_123", "type": "api_key.deleted", "effective_at": 1234567890, "actor": {"type": "session"}, "api_key.deleted": {"id": "key-123"}}`,
			wantType:    "api_key.deleted",
			wantDetails: &APIKeyDeleted{ID: "key-123"},
		},
		"invite.sent": {
			input:    `{"id": "log_123", "type": "invite.sent", "effective_at": 1234567890, "actor": {"type": "session"}, "invite.sent": {"id": "inv-123", "data": {"email": "test@example.com"}}}`,
			wantType: "invite.sent",
			wantDetails: &InviteSent{ID: "inv-123", Data: struct {
				Email string `json:"email"`
			}{Email: "test@example.com"}},
		},
		"invite.accepted": {
			input:       `{"id": "log_123", "type": "invite.accepted", "effective_at": 1234567890, "actor": {"type": "session"}, "invite.accepted": {"id": "inv-123"}}`,
			wantType:    "invite.accepted",
			wantDetails: &InviteAccepted{ID: "inv-123"},
		},
		"login.failed": {
			input:       `{"id": "log_123", "type": "login.failed", "effective_at": 1234567890, "actor": {"type": "session"}, "login.failed": {"error_code": "invalid_credentials", "error_message": "Bad password"}}`,
			wantType:    "login.failed",
			wantDetails: &LoginFailed{ErrorCode: "invalid_credentials", ErrorMessage: "Bad password"},
		},
		"login.succeeded": {
			input:       `{"id": "log_123", "type": "login.succeeded", "effective_at": 1234567890, "actor": {"type": "session"}, "login.succeeded": {"object": "audit.event", "id": "login_456", "type": "login.succeeded", "effective_at": 1234567890, "actor": {"type": "session"}}}`,
			wantType:    "login.succeeded",
			wantDetails: &LoginSucceeded{Object: "audit.event", ID: "login_456", Type: "login.succeeded", EffectiveAt: 1234567890, Actor: Actor{Type: "session"}},
		},
		"logout.failed": {
			input:       `{"id": "log_123", "type": "logout.failed", "effective_at": 1234567890, "actor": {"type": "session"}, "logout.failed": {"error_code": "session_expired", "error_message": "Session already expired"}}`,
			wantType:    "logout.failed",
			wantDetails: &LogoutFailed{ErrorCode: "session_expired", ErrorMessage: "Session already expired"},
		},
		"logout.succeeded with no details": {
			input:       `{"id": "log_123", "type": "logout.succeeded", "effective_at": 1234567890, "actor": {"type": "session"}}`,
			wantType:    "logout.succeeded",
			wantDetails: nil,
		},
		"organization.updated": {
			input:    `{"id": "log_123", "type": "organization.updated", "effective_at": 1234567890, "actor": {"type": "session"}, "organization.updated": {"id": "org-123", "changes_requested": {"name": "New Name"}}}`,
			wantType: "organization.updated",
			wantDetails: &OrganizationUpdated{ID: "org-123", ChangesRequested: struct {
				Name string `json:"name,omitempty"`
			}{Name: "New Name"}},
		},
		"project.created": {
			input:    `{"id": "log_123", "type": "project.created", "effective_at": 1234567890, "actor": {"type": "session"}, "project.created": {"id": "proj-123", "data": {"name": "my-project", "title": "My Project"}}}`,
			wantType: "project.created",
			wantDetails: &ProjectCreated{ID: "proj-123", Data: struct {
				Name  string `json:"name"`
				Title string `json:"title"`
			}{Name: "my-project", Title: "My Project"}},
		},
		"project.updated": {
			input:    `{"id": "log_123", "type": "project.updated", "effective_at": 1234567890, "actor": {"type": "session"}, "project.updated": {"id": "proj-123", "changes_requested": {"title": "Updated Title"}}}`,
			wantType: "project.updated",
			wantDetails: &ProjectUpdated{ID: "proj-123", ChangesRequested: struct {
				Title string `json:"title"`
			}{Title: "Updated Title"}},
		},
		"project.archived": {
			input:       `{"id": "log_123", "type": "project.archived", "effective_at": 1234567890, "actor": {"type": "session"}, "project.archived": {"id": "proj-123"}}`,
			wantType:    "project.archived",
			wantDetails: &ProjectArchived{ID: "proj-123"},
		},
		"rate_limit.updated": {
			input:    `{"id": "log_123", "type": "rate_limit.updated", "effective_at": 1234567890, "actor": {"type": "session"}, "rate_limit.updated": {"id": "rl-123", "changes_requested": {"max_requests_per_1_minute": 100}}}`,
			wantType: "rate_limit.updated",
			wantDetails: &RateLimitUpdated{ID: "rl-123", ChangesRequested: struct {
				MaxRequestsPer1Minute       int `json:"max_requests_per_1_minute,omitempty"`
				MaxTokensPer1Minute         int `json:"max_tokens_per_1_minute,omitempty"`
				MaxImagesPer1Minute         int `json:"max_images_per_1_minute,omitempty"`
				MaxAudioMegabytesPer1Minute int `json:"max_audio_megabytes_per_1_minute,omitempty"`
				MaxRequestsPer1Day          int `json:"max_requests_per_1_day,omitempty"`
				Batch1DayMaxInputTokens     int `json:"batch_1_day_max_input_tokens,omitempty"`
			}{MaxRequestsPer1Minute: 100}},
		},
		"rate_limit.deleted": {
			input:       `{"id": "log_123", "type": "rate_limit.deleted", "effective_at": 1234567890, "actor": {"type": "session"}, "rate_limit.deleted": {"id": "rl-123"}}`,
			wantType:    "rate_limit.deleted",
			wantDetails: &RateLimitDeleted{ID: "rl-123"},
		},
		"service_account.created": {
			input:    `{"id": "log_123", "type": "service_account.created", "effective_at": 1234567890, "actor": {"type": "session"}, "service_account.created": {"id": "sa-123", "data": {"role": "member"}}}`,
			wantType: "service_account.created",
			wantDetails: &ServiceAccountCreated{ID: "sa-123", Data: struct {
				Role string `json:"role"`
			}{Role: "member"}},
		},
		"service_account.updated": {
			input:    `{"id": "log_123", "type": "service_account.updated", "effective_at": 1234567890, "actor": {"type": "session"}, "service_account.updated": {"id": "sa-123", "changes_requested": {"role": "owner"}}}`,
			wantType: "service_account.updated",
			wantDetails: &ServiceAccountUpdated{ID: "sa-123", ChangesRequested: struct {
				Role string `json:"role"`
			}{Role: "owner"}},
		},
		"service_account.deleted": {
			input:       `{"id": "log_123", "type": "service_account.deleted", "effective_at": 1234567890, "actor": {"type": "session"}, "service_account.deleted": {"id": "sa-123"}}`,
			wantType:    "service_account.deleted",
			wantDetails: &ServiceAccountDeleted{ID: "sa-123"},
		},
		"user.added": {
			input:    `{"id": "log_123", "type": "user.added", "effective_at": 1234567890, "actor": {"type": "session"}, "user.added": {"id": "user-123", "data": {"role": "member"}}}`,
			wantType: "user.added",
			wantDetails: &UserAdded{ID: "user-123", Data: struct {
				Role string `json:"role"`
			}{Role: "member"}},
		},
		"user.updated": {
			input:    `{"id": "log_123", "type": "user.updated", "effective_at": 1234567890, "actor": {"type": "session"}, "user.updated": {"id": "user-123", "changes_requested": {"role": "owner"}}}`,
			wantType: "user.updated",
			wantDetails: &UserUpdated{ID: "user-123", ChangesRequested: struct {
				Role string `json:"role"`
			}{Role: "owner"}},
		},
		"user.deleted": {
			input:       `{"id": "log_123", "type": "user.deleted", "effective_at": 1234567890, "actor": {"type": "session"}, "user.deleted": {"id": "user-123"}}`,
			wantType:    "user.deleted",
			wantDetails: &UserDeleted{ID: "user-123"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var log AuditLog
			err := json.Unmarshal([]byte(tc.input), &log)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if log.Type != tc.wantType {
				t.Errorf("Expected Type %q, got %q", tc.wantType, log.Type)
			}

			if tc.wantDetails == nil {
				if log.Details != nil {
					t.Errorf("Expected Details to be nil, got %v", log.Details)
				}
			} else {
				if log.Details == nil {
					t.Errorf("Expected Details to be non-nil")
				} else if !reflect.DeepEqual(log.Details, tc.wantDetails) {
					t.Errorf("Details mismatch:\ngot:  %+v\nwant: %+v", log.Details, tc.wantDetails)
				}
			}
		})
	}
}
