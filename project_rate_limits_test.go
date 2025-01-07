package openaiorgs

import (
	"fmt"
	"testing"
)

func TestListProjectRateLimits(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockProjectRateLimits := []ProjectRateLimit{
		{
			Object:                      "project.rate_limit",
			ID:                          "rl-babbage-002",
			Model:                       "babbage-002",
			MaxRequestsPer1Minute:       3,
			MaxTokensPer1Minute:         150000,
			MaxImagesPer1Minute:         10,
			MaxAudioMegabytesPer1Minute: 0,
			MaxRequestsPer1Day:          200,
			Batch1DayMaxInputTokens:     0,
		},
	}

	// Register mock response
	response := ListResponse[ProjectRateLimit]{
		Object:  "list",
		Data:    mockProjectRateLimits,
		FirstID: "rl-babbage-002",
		LastID:  "rl-babbage-002",
		HasMore: false,
	}
	projectId := "proj_123"
	path := fmt.Sprintf("%s/%s/rate_limits", ProjectsListEndpoint, projectId)
	h.mockResponse("GET", path, 200, response)

	// Make the API call
	projectRateLimits, err := h.client.ListProjectRateLimits(10, "", projectId)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(projectRateLimits.Data) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projectRateLimits.Data))
		return
	}
	if mockProjectRateLimits[0].ID != projectRateLimits.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockProjectRateLimits[0].ID, projectRateLimits.Data[0].ID)
	}
	if mockProjectRateLimits[0].Model != projectRateLimits.Data[0].Model {
		t.Errorf("Expected Model %s, got %s", mockProjectRateLimits[0].Model, projectRateLimits.Data[0].Model)
	}

	// Verify the request was made
	h.assertRequest("GET", path, 1)
}

func TestModifyProjectRateLimits(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	rateLimitId := "rl-babbage-002"

	// Mock response data
	mockProjectRateLimit := ProjectRateLimit{
		Object:                      "project.rate_limit",
		ID:                          rateLimitId,
		Model:                       "babbage-002",
		MaxRequestsPer1Minute:       3,
		MaxTokensPer1Minute:         150000,
		MaxImagesPer1Minute:         10,
		MaxAudioMegabytesPer1Minute: 0,
		MaxRequestsPer1Day:          200,
		Batch1DayMaxInputTokens:     0,
	}

	expectedMaxRequestsPer1Minute := int64(2)

	// Register mock response
	response := ProjectRateLimit{
		Object:                      "project.rate_limit",
		ID:                          rateLimitId,
		Model:                       "babbage-002",
		MaxRequestsPer1Minute:       expectedMaxRequestsPer1Minute,
		MaxTokensPer1Minute:         150000,
		MaxImagesPer1Minute:         10,
		MaxAudioMegabytesPer1Minute: 0,
		MaxRequestsPer1Day:          200,
		Batch1DayMaxInputTokens:     0,
	}
	projectId := "proj_123"
	path := fmt.Sprintf("%s/%s/rate_limits/%s", ProjectsListEndpoint, projectId, rateLimitId)
	h.mockResponse("POST", path, 200, response)

	fields := ProjectRateLimitRequestFields{
		MaxRequestsPer1Minute: expectedMaxRequestsPer1Minute,
	}
	// Make the API call
	projectRateLimit, err := h.client.ModifyProjectRateLimit(projectId, rateLimitId, fields)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if mockProjectRateLimit.ID != projectRateLimit.ID {
		t.Errorf("Expected ID %s, got %s", mockProjectRateLimit.ID, projectRateLimit.ID)
	}
	if mockProjectRateLimit.Model != projectRateLimit.Model {
		t.Errorf("Expected Model %s, got %s", mockProjectRateLimit.Model, projectRateLimit.Model)
	}
	if expectedMaxRequestsPer1Minute != projectRateLimit.MaxRequestsPer1Minute {
		t.Errorf("Expected MaxRequestsPer1Minute %d, got %d", mockProjectRateLimit.MaxRequestsPer1Minute, projectRateLimit.MaxRequestsPer1Minute)
	}

	// Verify the request was made
	h.assertRequest("POST", path, 1)
}
