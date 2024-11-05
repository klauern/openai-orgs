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
			ID:        "log_123",
			Type:      "access_policy.created",
			Timestamp: now,
			Version:   "1.0",
			Actor:     Actor{ID: "actor_123", Name: "Test Actor", Type: "user"},
			Event:     Event{ID: "event_123", Type: "access_policy", Action: "created", Auth: Auth{Type: "token", Transport: "http"}},
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

	// Helper function to create a basic change struct
	changes := func(old, new string) map[string]interface{} {
		return map[string]interface{}{
			"changes": map[string]interface{}{
				"name": map[string]interface{}{
					"old": old,
					"new": new,
				},
			},
		}
	}

	tests := map[string]struct {
		logType string
		payload map[string]interface{}
		want    interface{}
		wantErr bool
	}{
		// Access Policy events
		"access_policy.created": {
			logType: "access_policy.created",
			payload: map[string]interface{}{"id": "pol_123", "name": "Test Policy"},
			want:    &AccessPolicyCreated{ID: "pol_123", Name: "Test Policy"},
		},
		"access_policy.deleted": {
			logType: "access_policy.deleted",
			payload: map[string]interface{}{"id": "pol_123"},
			want:    &AccessPolicyDeleted{ID: "pol_123"},
		},
		"access_policy.updated": {
			logType: "access_policy.updated",
			payload: map[string]interface{}{"id": "pol_123", "changes": changes("Old", "New")["changes"]},
			want: &AccessPolicyUpdated{
				ID: "pol_123",
				Changes: struct {
					Name struct {
						Old string `json:"old"`
						New string `json:"new"`
					} `json:"name"`
				}{
					Name: struct {
						Old string `json:"old"`
						New string `json:"new"`
					}{
						Old: "Old",
						New: "New",
					},
				},
			},
		},

		// API Key events
		"api_key.created": {
			logType: "api_key.created",
			payload: map[string]interface{}{"id": "key_123", "name": "Test Key"},
			want:    &APIKeyCreated{ID: "key_123", Name: "Test Key"},
		},
		"api_key.deleted": {
			logType: "api_key.deleted",
			payload: map[string]interface{}{"id": "key_123"},
			want:    &APIKeyDeleted{ID: "key_123"},
		},

		// Assistant events
		"assistant.created": {
			logType: "assistant.created",
			payload: map[string]interface{}{"id": "asst_123", "name": "Test Assistant"},
			want:    &AssistantCreated{ID: "asst_123", Name: "Test Assistant"},
		},
		"assistant.deleted": {
			logType: "assistant.deleted",
			payload: map[string]interface{}{"id": "asst_123"},
			want:    &AssistantDeleted{ID: "asst_123"},
		},
		"assistant.modified": {
			logType: "assistant.modified",
			payload: map[string]interface{}{"id": "asst_123", "changes": changes("Old", "New")["changes"]},
			want: &AssistantModified{ID: "asst_123", Changes: struct {
				Name struct {
					Old string `json:"old"`
					New string `json:"new"`
				} `json:"name"`
			}{Name: struct {
				Old string `json:"old"`
				New string `json:"new"`
			}{"Old", "New"}}},
		},

		// File events
		"file.created": {
			logType: "file.created",
			payload: map[string]interface{}{"id": "file_123", "name": "test.txt"},
			want:    &FileCreated{ID: "file_123", Name: "test.txt"},
		},
		"file.deleted": {
			logType: "file.deleted",
			payload: map[string]interface{}{"id": "file_123"},
			want:    &FileDeleted{ID: "file_123"},
		},

		// Fine-tune events
		"fine_tune.created": {
			logType: "fine_tune.created",
			payload: map[string]interface{}{"id": "ft_123"},
			want:    &FineTuneCreated{ID: "ft_123"},
		},
		"fine_tune.deleted": {
			logType: "fine_tune.deleted",
			payload: map[string]interface{}{"id": "ft_123"},
			want:    &FineTuneDeleted{ID: "ft_123"},
		},
		"fine_tune.event.created": {
			logType: "fine_tune.event.created",
			payload: map[string]interface{}{
				"id": "fte_123", "fine_tune_id": "ft_123",
				"level": "info", "message": "Test message",
				"created_at":    testTime.Unix(),
				"serialized_at": testTime.Unix(),
			},
			want: &FineTuneEventCreated{
				ID: "fte_123", FineTuneID: "ft_123",
				Level: "info", Message: "Test message",
				CreatedAt: testTime, SerializedAt: testTime,
			},
		},

		// Model events
		"model.created": {
			logType: "model.created",
			payload: map[string]interface{}{"id": "model_123", "name": "Test Model"},
			want:    &ModelCreated{ID: "model_123", Name: "Test Model"},
		},
		"model.deleted": {
			logType: "model.deleted",
			payload: map[string]interface{}{"id": "model_123"},
			want:    &ModelDeleted{ID: "model_123"},
		},

		// Run events
		"run.created": {
			logType: "run.created",
			payload: map[string]interface{}{
				"id": "run_123", "thread_id": "thread_123",
				"assistant_id": "asst_123", "status": "completed",
				"started_at": testTime.Unix(),
			},
			want: &RunCreated{
				ID: "run_123", ThreadID: "thread_123",
				AssistantID: "asst_123", Status: "completed",
				StartedAt: testTime,
			},
		},
		"run.modified": {
			logType: "run.modified",
			payload: map[string]interface{}{
				"id": "run_123",
				"changes": map[string]interface{}{
					"status": map[string]interface{}{
						"old": "running",
						"new": "completed",
					},
				},
			},
			want: &RunModified{ID: "run_123", Changes: struct {
				Status struct {
					Old string `json:"old"`
					New string `json:"new"`
				} `json:"status"`
			}{
				Status: struct {
					Old string `json:"old"`
					New string `json:"new"`
				}{"running", "completed"},
			}},
		},

		// Thread events
		"thread.created": {
			logType: "thread.created",
			payload: map[string]interface{}{"id": "thread_123"},
			want:    &ThreadCreated{ID: "thread_123"},
		},
		"thread.deleted": {
			logType: "thread.deleted",
			payload: map[string]interface{}{"id": "thread_123"},
			want:    &ThreadDeleted{ID: "thread_123"},
		},
		"thread.modified": {
			logType: "thread.modified",
			payload: map[string]interface{}{
				"id": "thread_123",
				"changes": map[string]interface{}{
					"metadata": map[string]interface{}{
						"old": map[string]interface{}{"key": "old_value"},
						"new": map[string]interface{}{"key": "new_value"},
					},
				},
			},
			want: &ThreadModified{ID: "thread_123", Changes: struct {
				Metadata struct {
					Old map[string]interface{} `json:"old"`
					New map[string]interface{} `json:"new"`
				} `json:"metadata"`
			}{
				Metadata: struct {
					Old map[string]interface{} `json:"old"`
					New map[string]interface{} `json:"new"`
				}{
					Old: map[string]interface{}{"key": "old_value"},
					New: map[string]interface{}{"key": "new_value"},
				},
			}},
		},

		// Organization events
		"invite.sent": {
			logType: "invite.sent",
			payload: map[string]interface{}{"email": "test@example.com"},
			want:    &InviteSent{Email: "test@example.com"},
		},
		"login.succeeded": {
			logType: "login.succeeded",
			payload: map[string]interface{}{},
			want:    &LoginSucceeded{},
		},
		"logout.succeeded": {
			logType: "logout.succeeded",
			payload: map[string]interface{}{},
			want:    &LogoutSucceeded{},
		},
		"organization.updated": {
			logType: "organization.updated",
			payload: map[string]interface{}{"changes": changes("Old Org", "New Org")["changes"]},
			want: &OrganizationUpdated{Changes: struct {
				Name struct {
					Old string `json:"old"`
					New string `json:"new"`
				} `json:"name,omitempty"`
			}{
				Name: struct {
					Old string `json:"old"`
					New string `json:"new"`
				}{"Old Org", "New Org"},
			}},
		},

		// Error cases
		"unknown type": {
			logType: "unknown.type",
			payload: map[string]interface{}{},
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			auditLog := &AuditLog{
				Type:  tt.logType,
				Event: Event{Payload: tt.payload},
			}

			got, err := ParseAuditLogPayload(auditLog)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAuditLogPayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAuditLogPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
