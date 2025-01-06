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
	projects, err := h.client.ListProjectRateLimits(10, "", projectId)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(projects.Data) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projects.Data))
		return
	}
	if mockProjectRateLimits[0].ID != projects.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockProjectRateLimits[0].ID, projects.Data[0].ID)
	}
	if mockProjectRateLimits[0].Model != projects.Data[0].Model {
		t.Errorf("Expected Model %s, got %s", mockProjectRateLimits[0].Model, projects.Data[0].Model)
	}

	// Verify the request was made
	h.assertRequest("GET", path, 1)
}
