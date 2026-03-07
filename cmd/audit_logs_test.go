package cmd

import (
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

func createTestAuditLog(id, logType string, details interface{}) openaiorgs.AuditLog {
	return openaiorgs.AuditLog{
		ID:          id,
		Type:        logType,
		EffectiveAt: openaiorgs.UnixSeconds(time.Now()),
		Actor: openaiorgs.Actor{
			Type: "session",
			Session: &openaiorgs.Session{
				User:      openaiorgs.AuditUser{ID: "user_123", Email: "test@example.com"},
				IPAddress: "1.2.3.4",
				UserAgent: "test-agent",
			},
		},
		Details: details,
	}
}

func createTestResponse(logs ...openaiorgs.AuditLog) *openaiorgs.ListResponse[openaiorgs.AuditLog] {
	return &openaiorgs.ListResponse[openaiorgs.AuditLog]{
		Object:  "list",
		Data:    logs,
		FirstID: "log_first",
		LastID:  "log_last",
		HasMore: false,
	}
}

func TestOutputJSON(t *testing.T) {
	log1 := createTestAuditLog("log_1", "api_key.created", &openaiorgs.APIKeyCreated{
		ID: "key_abc",
	})
	response := createTestResponse(log1)

	output := captureOutput(func() {
		err := outputJSON(response, false)
		if err != nil {
			t.Errorf("outputJSON() error = %v", err)
		}
	})

	if !strings.Contains(output, `"id"`) {
		t.Errorf("Expected JSON output to contain 'id' field, got: %s", output)
	}
	if !strings.Contains(output, `"log_1"`) {
		t.Errorf("Expected JSON output to contain log ID, got: %s", output)
	}
	if !strings.Contains(output, `"object"`) {
		t.Errorf("Expected JSON output to contain 'object' field, got: %s", output)
	}
	if !strings.Contains(output, `"list"`) {
		t.Errorf("Expected JSON output to contain 'list' object type, got: %s", output)
	}
}

func TestOutputJSONL(t *testing.T) {
	log1 := createTestAuditLog("log_1", "api_key.created", nil)
	log2 := createTestAuditLog("log_2", "invite.sent", nil)
	response := createTestResponse(log1, log2)

	t.Run("non-verbose", func(t *testing.T) {
		output := captureOutput(func() {
			err := outputJSONL(response, false)
			if err != nil {
				t.Errorf("outputJSONL() error = %v", err)
			}
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines (one per log entry), got %d: %s", len(lines), output)
		}
		if !strings.Contains(lines[0], "log_1") {
			t.Errorf("Expected first line to contain log_1, got: %s", lines[0])
		}
		if !strings.Contains(lines[1], "log_2") {
			t.Errorf("Expected second line to contain log_2, got: %s", lines[1])
		}
	})

	t.Run("verbose", func(t *testing.T) {
		output := captureOutput(func() {
			err := outputJSONL(response, true)
			if err != nil {
				t.Errorf("outputJSONL() error = %v", err)
			}
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")
		// verbose adds metadata line first, then 2 log entries
		if len(lines) != 3 {
			t.Errorf("Expected 3 lines (metadata + 2 entries), got %d: %s", len(lines), output)
		}
		if !strings.Contains(lines[0], `"total"`) {
			t.Errorf("Expected metadata line to contain 'total', got: %s", lines[0])
		}
		if !strings.Contains(lines[0], `"first_id"`) {
			t.Errorf("Expected metadata line to contain 'first_id', got: %s", lines[0])
		}
		if !strings.Contains(lines[0], `"last_id"`) {
			t.Errorf("Expected metadata line to contain 'last_id', got: %s", lines[0])
		}
	})
}

func TestOutputPretty(t *testing.T) {
	log1 := createTestAuditLog("log_1", "api_key.created", &openaiorgs.APIKeyCreated{
		ID: "key_abc",
	})
	response := createTestResponse(log1)

	output := captureOutput(func() {
		err := outputPretty(response, false)
		if err != nil {
			t.Errorf("outputPretty() error = %v", err)
		}
	})

	if !strings.Contains(output, "=== Audit Log Entry ===") {
		t.Errorf("Expected log entry marker in output, got: %s", output)
	}
	if !strings.Contains(output, "log_1") {
		t.Errorf("Expected log ID in output, got: %s", output)
	}
	if !strings.Contains(output, "api_key.created") {
		t.Errorf("Expected log type in output, got: %s", output)
	}
	if !strings.Contains(output, "API Key created with ID: key_abc") {
		t.Errorf("Expected API key detail in output, got: %s", output)
	}
}

func TestOutputPrettyVerbose(t *testing.T) {
	log1 := createTestAuditLog("log_1", "api_key.created", nil)
	response := createTestResponse(log1)

	output := captureOutput(func() {
		err := outputPretty(response, true)
		if err != nil {
			t.Errorf("outputPretty() error = %v", err)
		}
	})

	if !strings.Contains(output, "=== Audit Log Summary ===") {
		t.Errorf("Expected summary header in verbose output, got: %s", output)
	}
	if !strings.Contains(output, "Total logs: 1") {
		t.Errorf("Expected total logs count in verbose output, got: %s", output)
	}
	if !strings.Contains(output, "Actor Details:") {
		t.Errorf("Expected actor details in verbose output, got: %s", output)
	}
	if !strings.Contains(output, "User Email: test@example.com") {
		t.Errorf("Expected user email in verbose output, got: %s", output)
	}
	if !strings.Contains(output, "IP:         1.2.3.4") {
		t.Errorf("Expected IP address in verbose output, got: %s", output)
	}
	if !strings.Contains(output, "User Agent: test-agent") {
		t.Errorf("Expected user agent in verbose output, got: %s", output)
	}
}

func TestOutputPrettyWithNilSession(t *testing.T) {
	log1 := openaiorgs.AuditLog{
		ID:          "log_apikey",
		Type:        "api_key.created",
		EffectiveAt: openaiorgs.UnixSeconds(time.Now()),
		Actor: openaiorgs.Actor{
			Type:    "api_key",
			Session: nil,
			APIKey: &openaiorgs.APIKeyActor{
				Type: "user",
				User: openaiorgs.AuditUser{ID: "user_456", Email: "apikey@example.com"},
			},
		},
		Details: nil,
	}
	response := createTestResponse(log1)

	output := captureOutput(func() {
		err := outputPretty(response, true)
		if err != nil {
			t.Errorf("outputPretty() error = %v", err)
		}
	})

	if !strings.Contains(output, "API Key Type: user") {
		t.Errorf("Expected API key type in output, got: %s", output)
	}
	if !strings.Contains(output, "User Email:   apikey@example.com") {
		t.Errorf("Expected API key user email in output, got: %s", output)
	}
}

func TestOutputPrettyWithNilDetails(t *testing.T) {
	log1 := createTestAuditLog("log_nil_details", "unknown.event", nil)
	response := createTestResponse(log1)

	output := captureOutput(func() {
		err := outputPretty(response, false)
		if err != nil {
			t.Errorf("outputPretty() error = %v", err)
		}
	})

	if !strings.Contains(output, "log_nil_details") {
		t.Errorf("Expected log ID in output, got: %s", output)
	}
	if strings.Contains(output, "Payload Details:") {
		t.Errorf("Expected no payload details section for nil details, got: %s", output)
	}
}

func TestOutputPrettyDetailTypes(t *testing.T) {
	tests := []struct {
		name            string
		logType         string
		details         interface{}
		expectedStrings []string
	}{
		{
			name:    "APIKeyCreated",
			logType: "api_key.created",
			details: &openaiorgs.APIKeyCreated{
				ID: "key_123",
			},
			expectedStrings: []string{"API Key created with ID: key_123"},
		},
		{
			name:    "InviteSent",
			logType: "invite.sent",
			details: func() *openaiorgs.InviteSent {
				inv := &openaiorgs.InviteSent{ID: "inv_123"}
				inv.Data.Email = "invited@example.com"
				return inv
			}(),
			expectedStrings: []string{"Invite sent with ID: inv_123", "Email: invited@example.com"},
		},
		{
			name:    "ProjectCreated",
			logType: "project.created",
			details: func() *openaiorgs.ProjectCreated {
				pc := &openaiorgs.ProjectCreated{ID: "proj_123"}
				pc.Data.Name = "my-project"
				pc.Data.Title = "My Project"
				return pc
			}(),
			expectedStrings: []string{"Project created with ID: proj_123", "Name: my-project", "Title: My Project"},
		},
		{
			name:    "LoginFailed",
			logType: "login.failed",
			details: &openaiorgs.LoginFailed{
				ErrorCode:    "invalid_credentials",
				ErrorMessage: "Wrong password",
			},
			expectedStrings: []string{"Login failed", "Error code: invalid_credentials", "Error message: Wrong password"},
		},
		{
			name:    "UserAdded",
			logType: "user.added",
			details: func() *openaiorgs.UserAdded {
				ua := &openaiorgs.UserAdded{ID: "user_789"}
				ua.Data.Role = "member"
				return ua
			}(),
			expectedStrings: []string{"User added with ID: user_789", "Role: member"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := createTestAuditLog("log_"+tt.name, tt.logType, tt.details)
			response := createTestResponse(log)

			output := captureOutput(func() {
				err := outputPretty(response, false)
				if err != nil {
					t.Errorf("outputPretty() error = %v", err)
				}
			})

			for _, expected := range tt.expectedStrings {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, got: %s", expected, output)
				}
			}
		})
	}
}

func TestOutputResponse_UnknownFormat(t *testing.T) {
	response := createTestResponse()
	err := outputResponse(response, "xml", false)
	if err == nil {
		t.Error("Expected error for unknown format, got nil")
	}
	if !strings.Contains(err.Error(), "unknown output format: xml") {
		t.Errorf("Expected descriptive error message, got: %v", err)
	}
}

func TestOutputResponse_Routing(t *testing.T) {
	log1 := createTestAuditLog("log_route", "api_key.created", nil)
	response := createTestResponse(log1)

	tests := []struct {
		name          string
		format        string
		expectedInOut string
	}{
		{
			name:          "json format",
			format:        "json",
			expectedInOut: `"id"`,
		},
		{
			name:          "jsonl format",
			format:        "jsonl",
			expectedInOut: "log_route",
		},
		{
			name:          "pretty format",
			format:        "pretty",
			expectedInOut: "=== Audit Log Entry ===",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				err := outputResponse(response, tt.format, false)
				if err != nil {
					t.Errorf("outputResponse(%q) error = %v", tt.format, err)
				}
			})
			if !strings.Contains(output, tt.expectedInOut) {
				t.Errorf("outputResponse(%q) expected output to contain %q, got: %s", tt.format, tt.expectedInOut, output)
			}
		})
	}
}

func TestEmptyAuditLogs(t *testing.T) {
	response := createTestResponse()

	t.Run("json empty", func(t *testing.T) {
		output := captureOutput(func() {
			err := outputJSON(response, false)
			if err != nil {
				t.Errorf("outputJSON() error = %v", err)
			}
		})
		if !strings.Contains(output, `"data"`) {
			t.Errorf("Expected 'data' field in JSON output, got: %s", output)
		}
	})

	t.Run("jsonl empty", func(t *testing.T) {
		output := captureOutput(func() {
			err := outputJSONL(response, false)
			if err != nil {
				t.Errorf("outputJSONL() error = %v", err)
			}
		})
		if strings.TrimSpace(output) != "" {
			t.Errorf("Expected empty JSONL output for no entries, got: %s", output)
		}
	})

	t.Run("jsonl empty verbose", func(t *testing.T) {
		output := captureOutput(func() {
			err := outputJSONL(response, true)
			if err != nil {
				t.Errorf("outputJSONL() error = %v", err)
			}
		})
		if !strings.Contains(output, `"total":0`) {
			t.Errorf("Expected metadata with total:0, got: %s", output)
		}
	})

	t.Run("pretty empty", func(t *testing.T) {
		output := captureOutput(func() {
			err := outputPretty(response, false)
			if err != nil {
				t.Errorf("outputPretty() error = %v", err)
			}
		})
		if strings.Contains(output, "=== Audit Log Entry ===") {
			t.Errorf("Expected no log entry markers for empty data, got: %s", output)
		}
	})
}

func TestListAuditLogsCommand(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		log1 := createTestAuditLog("log_1", "api_key.created", nil)
		h.mockResponse("GET", "/organization/audit_logs", 200, createTestResponse(log1))

		output := captureOutput(func() {
			err := h.runCmd(AuditLogsCommand(), []string{"audit-logs"})
			if err != nil {
				t.Errorf("runCmd() error = %v", err)
			}
		})

		if !strings.Contains(output, "=== Audit Log Entry ===") {
			t.Errorf("Expected log entry in output, got: %s", output)
		}
		h.assertRequest("GET", "/organization/audit_logs", 1)
	})

	t.Run("error", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		h.mockResponse("GET", "/organization/audit_logs", 500, map[string]string{"error": "API error: rate limited"})

		err := h.runCmd(AuditLogsCommand(), []string{"audit-logs"})
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !strings.Contains(err.Error(), "list audit logs") {
			t.Errorf("Expected wrapped error message, got: %v", err)
		}
	})
}
