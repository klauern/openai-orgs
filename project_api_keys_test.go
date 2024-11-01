package openaiorgs

import (
	"testing"
	"time"
)

func TestListProjectApiKeys(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockApiKeys := []ProjectApiKey{
		{
			Object:        "api_key",
			ID:            "key_123",
			Name:          "Test API Key",
			RedactedValue: "sk-****",
			CreatedAt:     UnixSeconds(now),
			Owner: Owner{
				Object: "user",
				ID:     "user_123",
				Name:   "Test User",
				Type:   OwnerTypeUser,
			},
		},
	}

	// Register mock response
	response := ListResponse[ProjectApiKey]{
		Object:  "list",
		Data:    mockApiKeys,
		FirstID: "key_123",
		LastID:  "key_123",
		HasMore: false,
	}
	h.mockResponse("GET", "/organization/projects/proj_123/api_keys", 200, response)

	// Make the API call
	apiKeys, err := h.client.ListProjectApiKeys("proj_123", 10, "")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(apiKeys.Data) != 1 {
		t.Errorf("Expected 1 API key, got %d", len(apiKeys.Data))
		return
	}
	if mockApiKeys[0].ID != apiKeys.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockApiKeys[0].ID, apiKeys.Data[0].ID)
	}
	if mockApiKeys[0].Name != apiKeys.Data[0].Name {
		t.Errorf("Expected Name %s, got %s", mockApiKeys[0].Name, apiKeys.Data[0].Name)
	}

	// Verify the request was made
	h.assertRequest("GET", "/organization/projects/proj_123/api_keys", 1)
}

func TestRetrieveProjectApiKey(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	apiKeyID := "key_123"
	mockApiKey := ProjectApiKey{
		Object:        "api_key",
		ID:            apiKeyID,
		Name:          "Test API Key",
		RedactedValue: "sk-****",
		CreatedAt:     UnixSeconds(time.Now()),
		Owner: Owner{
			Object: "user",
			ID:     "user_123",
			Name:   "Test User",
			Type:   OwnerTypeUser,
		},
	}

	h.mockResponse("GET", "/organization/projects/proj_123/api_keys/"+apiKeyID, 200, mockApiKey)

	// Make the API call
	apiKey, err := h.client.RetrieveProjectApiKey("proj_123", apiKeyID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if apiKey == nil {
		t.Error("Expected API key, got nil")
		return
	}
	if mockApiKey.ID != apiKey.ID {
		t.Errorf("Expected ID %s, got %s", mockApiKey.ID, apiKey.ID)
	}

	// Verify the request was made
	h.assertRequest("GET", "/organization/projects/proj_123/api_keys/"+apiKeyID, 1)
}

func TestDeleteProjectApiKey(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	apiKeyID := "key_123"
	h.mockResponse("DELETE", "/organization/projects/proj_123/api_keys/"+apiKeyID, 204, nil)

	// Make the API call
	err := h.client.DeleteProjectApiKey("proj_123", apiKeyID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the request was made
	h.assertRequest("DELETE", "/organization/projects/proj_123/api_keys/"+apiKeyID, 1)
}
