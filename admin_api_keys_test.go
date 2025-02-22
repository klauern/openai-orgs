package openaiorgs

import (
	"testing"
	"time"
)

func TestListAdminAPIKeys(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockAPIKeys := []AdminAPIKey{
		{
			Object:        "api_key",
			ID:            "key_123",
			Name:          "Test API Key",
			RedactedValue: "sk-****",
			CreatedAt:     UnixSeconds(now),
			LastUsedAt:    UnixSeconds(now.Add(-1 * time.Hour)),
			Scopes:        []string{"organization.read", "organization.write"},
		},
	}

	// Register mock response
	response := ListResponse[AdminAPIKey]{
		Object:  "list",
		Data:    mockAPIKeys,
		FirstID: "key_123",
		LastID:  "key_123",
		HasMore: false,
	}
	h.mockResponse("GET", AdminAPIKeysEndpoint, 200, response)

	// Make the API call
	apiKeys, err := h.client.ListAdminAPIKeys(10, "")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(apiKeys.Data) != 1 {
		t.Errorf("Expected 1 API key, got %d", len(apiKeys.Data))
		return
	}
	if mockAPIKeys[0].ID != apiKeys.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockAPIKeys[0].ID, apiKeys.Data[0].ID)
	}
	if mockAPIKeys[0].Name != apiKeys.Data[0].Name {
		t.Errorf("Expected Name %s, got %s", mockAPIKeys[0].Name, apiKeys.Data[0].Name)
	}

	// Verify the request was made
	h.assertRequest("GET", AdminAPIKeysEndpoint, 1)
}

func TestCreateAdminAPIKey(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	mockAPIKey := AdminAPIKey{
		Object:        "api_key",
		ID:            "key_123",
		Name:          "New API Key",
		RedactedValue: "sk-****",
		CreatedAt:     UnixSeconds(time.Now()),
		LastUsedAt:    UnixSeconds(time.Now()),
		Scopes:        []string{"organization.read"},
	}

	h.mockResponse("POST", AdminAPIKeysEndpoint, 200, mockAPIKey)

	// Make the API call
	apiKey, err := h.client.CreateAdminAPIKey("New API Key", []string{"organization.read"})
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiKey == nil {
		t.Error("Expected API key, got nil")
		return
	}
	if mockAPIKey.ID != apiKey.ID {
		t.Errorf("Expected ID %s, got %s", mockAPIKey.ID, apiKey.ID)
	}
	if mockAPIKey.Name != apiKey.Name {
		t.Errorf("Expected Name %s, got %s", mockAPIKey.Name, apiKey.Name)
	}

	// Verify the request was made
	h.assertRequest("POST", AdminAPIKeysEndpoint, 1)
}

func TestRetrieveAdminAPIKey(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	apiKeyID := "key_123"
	mockAPIKey := AdminAPIKey{
		Object:        "api_key",
		ID:            apiKeyID,
		Name:          "Test API Key",
		RedactedValue: "sk-****",
		CreatedAt:     UnixSeconds(time.Now()),
		LastUsedAt:    UnixSeconds(time.Now()),
		Scopes:        []string{"organization.read", "organization.write"},
	}

	h.mockResponse("GET", AdminAPIKeysEndpoint+"/"+apiKeyID, 200, mockAPIKey)

	// Make the API call
	apiKey, err := h.client.RetrieveAdminAPIKey(apiKeyID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiKey == nil {
		t.Error("Expected API key, got nil")
		return
	}
	if mockAPIKey.ID != apiKey.ID {
		t.Errorf("Expected ID %s, got %s", mockAPIKey.ID, apiKey.ID)
	}

	// Verify the request was made
	h.assertRequest("GET", AdminAPIKeysEndpoint+"/"+apiKeyID, 1)
}

func TestDeleteAdminAPIKey(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	apiKeyID := "key_123"
	h.mockResponse("DELETE", AdminAPIKeysEndpoint+"/"+apiKeyID, 204, nil)

	// Make the API call
	err := h.client.DeleteAdminAPIKey(apiKeyID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the request was made
	h.assertRequest("DELETE", AdminAPIKeysEndpoint+"/"+apiKeyID, 1)
}
