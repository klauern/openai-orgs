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

func TestListProjectRateLimitsPaginated(t *testing.T) {
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
		HasMore: true,
	}
	projectId := "proj_123"
	path := fmt.Sprintf("%s/%s/rate_limits", ProjectsListEndpoint, projectId)
	h.mockResponse("GET", path, 200, response)

	// Make the API call
	projectRateLimits, err := h.client.ListProjectRateLimits(1, "", projectId)
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

	if !projectRateLimits.HasMore {
		t.Errorf("Expected there to be more project rate limits")
		return
	}

	// Verify the request was made
	h.assertRequest("GET", path, 1)

	// Mock response data
	mockProjectRateLimits = []ProjectRateLimit{
		{
			Object:                      "project.rate_limit",
			ID:                          "rl-babbage-003",
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
	response = ListResponse[ProjectRateLimit]{
		Object:  "list",
		Data:    mockProjectRateLimits,
		FirstID: "rl-babbage-003",
		LastID:  "rl-babbage-003",
		HasMore: false,
	}
	h.mockResponse("GET", path, 200, response)

	projectRateLimits, err = h.client.ListProjectRateLimits(1, projectRateLimits.LastID, projectId)
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

	if projectRateLimits.HasMore {
		t.Errorf("Expected there to be no more project rate limits")
		return
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
	expectedMaxTokensPer1Minute := int64(70000)
	expectedMaxImagesPer1Minute := int64(5)
	expectedMaxAudioMegabytesPer1Minute := int64(2)
	expectedMaxRequestsPer1Day := int64(100)
	expectedBatch1DayMaxInputTokens := int64(5)

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
		MaxTokensPer1Minute: expectedMaxTokensPer1Minute,
		MaxImagesPer1Minute: expectedMaxImagesPer1Minute,
		MaxAudioMegabytesPer1Minute: expectedMaxAudioMegabytesPer1Minute,
		MaxRequestsPer1Day: expectedMaxRequestsPer1Day,
		Batch1DayMaxInputTokens: expectedBatch1DayMaxInputTokens,
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
