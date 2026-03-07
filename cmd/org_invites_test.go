package cmd

import (
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

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

func TestListInvitesCommand(t *testing.T) {
	t.Run("successful list with nil AcceptedAt", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		h.mockResponse("GET", "/organization/invites", 200, openaiorgs.ListResponse[openaiorgs.Invite]{
			Object: "list",
			Data: []openaiorgs.Invite{
				createMockInvite("inv_123", "alice@example.com", "member", "pending", nil),
			},
			FirstID: "inv_123",
			LastID:  "inv_123",
			HasMore: false,
		})

		output := captureOutput(func() {
			err := h.runCmd(InvitesCommand(), []string{"invites", "list"})
			if err != nil {
				t.Errorf("runCmd() error = %v", err)
			}
		})

		for _, want := range []string{
			"ID | Email | Role | Status | Created At | Expires At | Accepted At",
			"N/A",
			"inv_123",
		} {
			if !strings.Contains(output, want) {
				t.Errorf("Expected output to contain %q, got: %s", want, output)
			}
		}
		h.assertRequest("GET", "/organization/invites", 1)
	})

	t.Run("successful list with AcceptedAt set", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		accepted := openaiorgs.UnixSeconds(time.Now())
		h.mockResponse("GET", "/organization/invites", 200, openaiorgs.ListResponse[openaiorgs.Invite]{
			Object: "list",
			Data: []openaiorgs.Invite{
				createMockInvite("inv_456", "bob@example.com", "owner", "accepted", &accepted),
			},
			FirstID: "inv_456",
			LastID:  "inv_456",
			HasMore: false,
		})

		output := captureOutput(func() {
			err := h.runCmd(InvitesCommand(), []string{"invites", "list"})
			if err != nil {
				t.Errorf("runCmd() error = %v", err)
			}
		})

		if strings.Contains(output, "N/A") {
			t.Errorf("Expected no N/A when AcceptedAt is set, got: %s", output)
		}
	})

	t.Run("error from API", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		h.mockResponse("GET", "/organization/invites", 500, map[string]string{"error": "API error"})

		err := h.runCmd(InvitesCommand(), []string{"invites", "list"})
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		h.mockResponse("GET", "/organization/invites", 200, openaiorgs.ListResponse[openaiorgs.Invite]{
			Object:  "list",
			Data:    []openaiorgs.Invite{},
			HasMore: false,
		})

		output := captureOutput(func() {
			err := h.runCmd(InvitesCommand(), []string{"invites", "list"})
			if err != nil {
				t.Errorf("runCmd() error = %v", err)
			}
		})

		if !strings.Contains(output, "ID | Email | Role | Status") {
			t.Errorf("Expected table headers even for empty list, got: %s", output)
		}
	})
}

func TestCreateInviteCommand(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		h.mockResponse("POST", "/organization/invites", 200,
			createMockInvite("inv_123", "alice@example.com", "member", "pending", nil))

		output := captureOutput(func() {
			err := h.runCmd(InvitesCommand(), []string{"invites", "create", "--email", "alice@example.com", "--role", "member"})
			if err != nil {
				t.Errorf("runCmd() error = %v", err)
			}
		})

		for _, want := range []string{"Invite created:", "ID: inv_123", "Email: alice@example.com"} {
			if !strings.Contains(output, want) {
				t.Errorf("Expected output to contain %q, got: %s", want, output)
			}
		}
		h.assertRequest("POST", "/organization/invites", 1)
	})

	t.Run("error from API", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		h.mockResponse("POST", "/organization/invites", 500, map[string]string{"error": "create failed"})

		err := h.runCmd(InvitesCommand(), []string{"invites", "create", "--email", "bad@example.com", "--role", "member"})
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestRetrieveInviteCommand(t *testing.T) {
	t.Run("successful retrieve", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		h.mockResponse("GET", "/organization/invites/inv_123", 200,
			createMockInvite("inv_123", "alice@example.com", "member", "pending", nil))

		output := captureOutput(func() {
			err := h.runCmd(InvitesCommand(), []string{"invites", "retrieve", "--id", "inv_123"})
			if err != nil {
				t.Errorf("runCmd() error = %v", err)
			}
		})

		for _, want := range []string{"Invite details:", "ID: inv_123"} {
			if !strings.Contains(output, want) {
				t.Errorf("Expected output to contain %q, got: %s", want, output)
			}
		}
		h.assertRequest("GET", "/organization/invites/inv_123", 1)
	})

	t.Run("error from API", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		h.mockResponse("GET", "/organization/invites/inv_999", 404, map[string]string{"error": "invite not found"})

		err := h.runCmd(InvitesCommand(), []string{"invites", "retrieve", "--id", "inv_999"})
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestDeleteInviteCommand(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		h.mockResponse("DELETE", "/organization/invites/inv_123", 204, nil)

		output := captureOutput(func() {
			err := h.runCmd(InvitesCommand(), []string{"invites", "delete", "--id", "inv_123"})
			if err != nil {
				t.Errorf("runCmd() error = %v", err)
			}
		})

		if !strings.Contains(output, "Invite inv_123 deleted successfully") {
			t.Errorf("Expected delete success message, got: %s", output)
		}
		h.assertRequest("DELETE", "/organization/invites/inv_123", 1)
	})

	t.Run("error from API", func(t *testing.T) {
		h := newCmdTestHelper(t)
		defer h.cleanup()

		h.mockResponse("DELETE", "/organization/invites/inv_999", 500, map[string]string{"error": "delete failed"})

		err := h.runCmd(InvitesCommand(), []string{"invites", "delete", "--id", "inv_999"})
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
