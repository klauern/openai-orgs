package openaiorgs

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestAPIClient_ErrorHandling(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	t.Run("ListAdminAPIKeys_NetworkError", func(t *testing.T) {
		// Mock network error
		httpmock.RegisterResponder("GET", testBaseURL+AdminAPIKeysEndpoint,
			httpmock.NewErrorResponder(http.ErrServerClosed))

		_, err := h.client.ListAdminAPIKeys(0, "")
		if err == nil {
			t.Error("Expected network error, got nil")
		}
	})

	t.Run("ListAdminAPIKeys_UnauthorizedError", func(t *testing.T) {
		// Mock 401 Unauthorized
		httpmock.RegisterResponder("GET", testBaseURL+AdminAPIKeysEndpoint,
			httpmock.NewStringResponder(401, `{"error": {"message": "Invalid API key", "type": "invalid_request_error"}}`))

		_, err := h.client.ListAdminAPIKeys(0, "")
		if err == nil {
			t.Error("Expected unauthorized error, got nil")
		}
		if !strings.Contains(err.Error(), "401") {
			t.Errorf("Expected error to mention status code 401, got: %v", err)
		}
	})

	t.Run("ListAdminAPIKeys_ServerError", func(t *testing.T) {
		// Mock 500 Internal Server Error
		httpmock.RegisterResponder("GET", testBaseURL+AdminAPIKeysEndpoint,
			httpmock.NewStringResponder(500, `{"error": {"message": "Internal server error", "type": "api_error"}}`))

		_, err := h.client.ListAdminAPIKeys(0, "")
		if err == nil {
			t.Error("Expected server error, got nil")
		}
		if !strings.Contains(err.Error(), "500") {
			t.Errorf("Expected error to mention status code 500, got: %v", err)
		}
	})

	t.Run("RetrieveAdminAPIKey_NotFound", func(t *testing.T) {
		// Mock 404 Not Found
		httpmock.RegisterResponder("GET", testBaseURL+AdminAPIKeysEndpoint+"/nonexistent",
			httpmock.NewStringResponder(404, `{"error": {"message": "API key not found", "type": "invalid_request_error"}}`))

		_, err := h.client.RetrieveAdminAPIKey("nonexistent")
		if err == nil {
			t.Error("Expected not found error, got nil")
		}
		if !strings.Contains(err.Error(), "404") {
			t.Errorf("Expected error to mention status code 404, got: %v", err)
		}
	})

	t.Run("CreateAdminAPIKey_ValidationError", func(t *testing.T) {
		// Mock 400 Bad Request for validation error
		httpmock.RegisterResponder("POST", testBaseURL+AdminAPIKeysEndpoint,
			httpmock.NewStringResponder(400, `{"error": {"message": "Invalid name parameter", "type": "invalid_request_error"}}`))

		_, err := h.client.CreateAdminAPIKey("", []string{})
		if err == nil {
			t.Error("Expected validation error, got nil")
		}
		if !strings.Contains(err.Error(), "400") {
			t.Errorf("Expected error to mention status code 400, got: %v", err)
		}
	})

	t.Run("ListUsers_MalformedJSON", func(t *testing.T) {
		// Mock response with malformed JSON
		httpmock.RegisterResponder("GET", testBaseURL+"/organization/users",
			httpmock.NewStringResponder(200, `{"object": "list", "data": [{"invalid": json}]}`))

		_, err := h.client.ListUsers(0, "")
		if err == nil {
			t.Error("Expected JSON parse error, got nil")
		}
		if !strings.Contains(err.Error(), "invalid character") {
			t.Errorf("Expected JSON parse error, got: %v", err)
		}
	})
}

func TestNewClient_EdgeCases(t *testing.T) {
	t.Run("EmptyBaseURL_UsesDefault", func(t *testing.T) {
		client := NewClient("", "test-token")
		if client == nil {
			t.Error("Expected client to be created with default baseURL")
			return
		}
		// Check that it uses the default base URL
		if !strings.Contains(client.BaseURL, "api.openai.com") {
			t.Errorf("Expected default baseURL to contain 'api.openai.com', got: %s", client.BaseURL)
		}
	})

	t.Run("CustomBaseURL", func(t *testing.T) {
		customURL := "https://custom.api.com/v1"
		client := NewClient(customURL, "test-token")
		if client == nil {
			t.Error("Expected client to be created with custom baseURL")
			return
		}
		if client.BaseURL != customURL {
			t.Errorf("Expected baseURL to be %s, got: %s", customURL, client.BaseURL)
		}
	})

	t.Run("EmptyAPIKey", func(t *testing.T) {
		client := NewClient("https://api.openai.com/v1", "")
		if client == nil {
			t.Error("Expected client to be created even with empty API key")
		}
	})
}
