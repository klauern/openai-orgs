package openaiorgs

import (
	"testing"
	"time"
)

func TestListInvites(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockInvites := []Invite{
		{
			ObjectType: "organization.invite",
			ID:         "invite_123421351236",
			Email:      "test@example.com",
			Role:       "admin",
			Status:     "pending",
			CreatedAt:  UnixSeconds(now),
			ExpiresAt:  UnixSeconds(now.Add(24 * time.Hour)),
		},
	}

	// Register mock response
	response := ListResponse[Invite]{
		Object:  "list",
		Data:    mockInvites,
		FirstID: "inv_123",
		LastID:  "inv_123",
		HasMore: false,
	}
	h.mockResponse("GET", InviteListEndpoint, 200, response)

	// Make the API call
	resp, err := h.client.ListInvites(100, "")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if resp == nil {
		t.Error("Expected response, got nil")
		return
	}
	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 invite, got %d", len(resp.Data))
		return
	}
	if mockInvites[0].ID != resp.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockInvites[0].ID, resp.Data[0].ID)
	}
	if mockInvites[0].Email != resp.Data[0].Email {
		t.Errorf("Expected Email %s, got %s", mockInvites[0].Email, resp.Data[0].Email)
	}

	// Verify the request was made
	h.assertRequest("GET", InviteListEndpoint, 1)
}

func TestCreateInvite(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	mockInvite := Invite{
		ObjectType: "organization_invite",
		ID:         "inv_123",
		Email:      "new@example.com",
		Role:       "admin",
		Status:     "pending",
		CreatedAt:  UnixSeconds(time.Now()),
		ExpiresAt:  UnixSeconds(time.Now().Add(24 * time.Hour)),
	}

	h.mockResponse("POST", InviteListEndpoint, 200, mockInvite)

	// Make the API call
	invite, err := h.client.CreateInvite("new@example.com", "admin")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if invite == nil {
		t.Error("Expected invite, got nil")
		return
	}
	if mockInvite.ID != invite.ID {
		t.Errorf("Expected ID %s, got %s", mockInvite.ID, invite.ID)
	}
	if mockInvite.Email != invite.Email {
		t.Errorf("Expected Email %s, got %s", mockInvite.Email, invite.Email)
	}

	// Verify the request was made
	h.assertRequest("POST", InviteListEndpoint, 1)
}

func TestRetrieveInvite(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	inviteID := "inv_123"
	mockInvite := Invite{
		ObjectType: "organization_invite",
		ID:         inviteID,
		Email:      "test@example.com",
		Role:       "admin",
		Status:     "pending",
		CreatedAt:  UnixSeconds(time.Now()),
		ExpiresAt:  UnixSeconds(time.Now().Add(24 * time.Hour)),
	}

	h.mockResponse("GET", InviteListEndpoint+"/"+inviteID, 200, mockInvite)

	// Make the API call
	invite, err := h.client.RetrieveInvite(inviteID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if invite == nil {
		t.Error("Expected invite, got nil")
		return
	}
	if mockInvite.ID != invite.ID {
		t.Errorf("Expected ID %s, got %s", mockInvite.ID, invite.ID)
	}

	// Verify the request was made
	h.assertRequest("GET", InviteListEndpoint+"/"+inviteID, 1)
}

func TestDeleteInvite(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	inviteID := "inv_123"
	h.mockResponse("DELETE", InviteListEndpoint+"/"+inviteID, 204, nil)

	// Make the API call
	err := h.client.DeleteInvite(inviteID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the request was made
	h.assertRequest("DELETE", InviteListEndpoint+"/"+inviteID, 1)
}
