package openaiorgs

import (
	"testing"
	"time"
)

func TestListProjectUsers(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockUsers := []ProjectUser{
		{
			Object:  "project_user",
			ID:      "user_123",
			Name:    "Test User",
			Email:   "test@example.com",
			Role:    "member",
			AddedAt: UnixSeconds(now),
		},
	}

	// Register mock response
	response := ListResponse[ProjectUser]{
		Object:  "list",
		Data:    mockUsers,
		FirstID: "user_123",
		LastID:  "user_123",
		HasMore: false,
	}
	h.mockResponse("GET", "/organization/projects/proj_123/users", 200, response)

	// Make the API call
	users, err := h.client.ListProjectUsers("proj_123", 10, "")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(users.Data) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users.Data))
		return
	}
	if mockUsers[0].ID != users.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockUsers[0].ID, users.Data[0].ID)
	}
	if mockUsers[0].Name != users.Data[0].Name {
		t.Errorf("Expected Name %s, got %s", mockUsers[0].Name, users.Data[0].Name)
	}

	// Verify the request was made
	h.assertRequest("GET", "/organization/projects/proj_123/users", 1)
}

func TestCreateProjectUser(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	mockUser := ProjectUser{
		Object:  "project_user",
		ID:      "user_123",
		Name:    "New User",
		Email:   "new@example.com",
		Role:    "member",
		AddedAt: UnixSeconds(time.Now()),
	}

	h.mockResponse("POST", "/organization/projects/proj_123/users", 200, mockUser)

	// Make the API call
	user, err := h.client.CreateProjectUser("proj_123", "user_123", "member")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
		return
	}
	if mockUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", mockUser.ID, user.ID)
	}
	if mockUser.Name != user.Name {
		t.Errorf("Expected Name %s, got %s", mockUser.Name, user.Name)
	}

	// Verify the request was made
	h.assertRequest("POST", "/organization/projects/proj_123/users", 1)
}

func TestRetrieveProjectUser(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	userID := "user_123"
	mockUser := ProjectUser{
		Object:  "project_user",
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		Role:    "member",
		AddedAt: UnixSeconds(time.Now()),
	}

	h.mockResponse("GET", "/organization/projects/proj_123/users/"+userID, 200, mockUser)

	// Make the API call
	user, err := h.client.RetrieveProjectUser("proj_123", userID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
		return
	}
	if mockUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", mockUser.ID, user.ID)
	}

	// Verify the request was made
	h.assertRequest("GET", "/organization/projects/proj_123/users/"+userID, 1)
}

func TestModifyProjectUser(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	userID := "user_123"
	mockUser := ProjectUser{
		Object:  "project_user",
		ID:      userID,
		Name:    "Test User",
		Email:   "test@example.com",
		Role:    "owner",
		AddedAt: UnixSeconds(time.Now()),
	}

	h.mockResponse("POST", "/organization/projects/proj_123/users/"+userID, 200, mockUser)

	// Make the API call
	user, err := h.client.ModifyProjectUser("proj_123", userID, "owner")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
		return
	}
	if mockUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", mockUser.ID, user.ID)
	}
	if mockUser.Role != user.Role {
		t.Errorf("Expected Role %s, got %s", mockUser.Role, user.Role)
	}

	// Verify the request was made
	h.assertRequest("POST", "/organization/projects/proj_123/users/"+userID, 1)
}

func TestDeleteProjectUser(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	userID := "user_123"
	h.mockResponse("DELETE", "/organization/projects/proj_123/users/"+userID, 204, nil)

	// Make the API call
	err := h.client.DeleteProjectUser("proj_123", userID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the request was made
	h.assertRequest("DELETE", "/organization/projects/proj_123/users/"+userID, 1)
}
