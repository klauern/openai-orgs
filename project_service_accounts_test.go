package openaiorgs

import (
	"testing"
	"time"
)

func TestListProjectServiceAccounts(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockAccounts := []ProjectServiceAccount{
		{
			Object:    "project_service_account",
			ID:        "service_account_123",
			Name:      "Test Service Account",
			Role:      "admin",
			CreatedAt: UnixSeconds(now),
			APIKey:    nil,
		},
	}

	// Register mock response
	response := ListResponse[ProjectServiceAccount]{
		Object:  "list",
		Data:    mockAccounts,
		FirstID: "service_account_123",
		LastID:  "service_account_123",
		HasMore: false,
	}
	h.mockResponse("GET", "/organization/projects/proj_123/service_accounts", 200, response)

	// Make the API call
	accounts, err := h.client.ListProjectServiceAccounts("proj_123", 10, "")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(accounts.Data) != 1 {
		t.Errorf("Expected 1 service account, got %d", len(accounts.Data))
		return
	}
	if mockAccounts[0].ID != accounts.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockAccounts[0].ID, accounts.Data[0].ID)
	}
	if mockAccounts[0].Name != accounts.Data[0].Name {
		t.Errorf("Expected Name %s, got %s", mockAccounts[0].Name, accounts.Data[0].Name)
	}

	// Verify the request was made
	h.assertRequest("GET", "/organization/projects/proj_123/service_accounts", 1)
}

func TestCreateProjectServiceAccount(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	mockAccount := ProjectServiceAccount{
		Object:    "project_service_account",
		ID:        "service_account_123",
		Name:      "New Service Account",
		Role:      "admin",
		CreatedAt: UnixSeconds(time.Now()),
		APIKey: &ProjectServiceAccountAPIKey{
			Object:    "project_service_account_api_key",
			Value:     "sk-api-key-123",
			Name:      nil,
			CreatedAt: UnixSeconds(time.Now()),
			ID:        "api_key_123",
		},
	}

	h.mockResponse("POST", "/organization/projects/proj_123/service_accounts", 200, mockAccount)

	// Make the API call
	account, err := h.client.CreateProjectServiceAccount("proj_123", "New Service Account")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if account == nil {
		t.Error("Expected service account, got nil")
		return
	}
	if mockAccount.ID != account.ID {
		t.Errorf("Expected ID %s, got %s", mockAccount.ID, account.ID)
	}
	if mockAccount.Name != account.Name {
		t.Errorf("Expected Name %s, got %s", mockAccount.Name, account.Name)
	}

	// Verify the request was made
	h.assertRequest("POST", "/organization/projects/proj_123/service_accounts", 1)
}

func TestRetrieveProjectServiceAccount(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	accountID := "service_account_123"
	mockAccount := ProjectServiceAccount{
		Object:    "project_service_account",
		ID:        accountID,
		Name:      "Test Service Account",
		Role:      "admin",
		CreatedAt: UnixSeconds(time.Now()),
		APIKey:    nil,
	}

	h.mockResponse("GET", "/organization/projects/proj_123/service_accounts/"+accountID, 200, mockAccount)

	// Make the API call
	account, err := h.client.RetrieveProjectServiceAccount("proj_123", accountID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if account == nil {
		t.Error("Expected service account, got nil")
		return
	}
	if mockAccount.ID != account.ID {
		t.Errorf("Expected ID %s, got %s", mockAccount.ID, account.ID)
	}

	// Verify the request was made
	h.assertRequest("GET", "/organization/projects/proj_123/service_accounts/"+accountID, 1)
}

func TestDeleteProjectServiceAccount(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	accountID := "service_account_123"
	h.mockResponse("DELETE", "/organization/projects/proj_123/service_accounts/"+accountID, 204, nil)

	// Make the API call
	err := h.client.DeleteProjectServiceAccount("proj_123", accountID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the request was made
	h.assertRequest("DELETE", "/organization/projects/proj_123/service_accounts/"+accountID, 1)
}
