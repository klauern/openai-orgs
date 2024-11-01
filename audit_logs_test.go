package openaiorgs

import (
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
	auditLog := &AuditLog{
		Type: "access_policy.created",
		Event: Event{
			Payload: map[string]interface{}{
				"id":   "policy_123",
				"name": "Test Policy",
			},
		},
	}

	payload, err := ParseAuditLogPayload(auditLog)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	accessPolicyCreated, ok := payload.(*AccessPolicyCreated)
	if !ok {
		t.Errorf("Expected payload to be of type *AccessPolicyCreated, got %T", payload)
		return
	}

	if accessPolicyCreated.ID != "policy_123" {
		t.Errorf("Expected ID %s, got %s", "policy_123", accessPolicyCreated.ID)
	}
	if accessPolicyCreated.Name != "Test Policy" {
		t.Errorf("Expected Name %s, got %s", "Test Policy", accessPolicyCreated.Name)
	}
}
