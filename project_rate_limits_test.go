package openaiorgs

import (
	"testing"
	"time"
)

func TestListProjectRateLimits(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockProjectRateLimits := []ProjectRateLimit{
		{
			Object:    "project",
			ID:        "proj_123",
			Name:      "Test ProjectRateLimit",
			CreatedAt: UnixSeconds(now),
			Status:    "active",
		},
	}

	// Register mock response
	response := ListResponse[ProjectRateLimit]{
		Object:  "list",
		Data:    mockProjectRateLimits,
		FirstID: "proj_123",
		LastID:  "proj_123",
		HasMore: false,
	}
	h.mockResponse("GET", ProjectRateLimitsListEndpoint, 200, response)

	// Make the API call
	projects, err := h.client.ListProjectRateLimits(10, "", false)
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
	if mockProjectRateLimits[0].Name != projects.Data[0].Name {
		t.Errorf("Expected Name %s, got %s", mockProjectRateLimits[0].Name, projects.Data[0].Name)
	}

	// Verify the request was made
	h.assertRequest("GET", ProjectRateLimitsListEndpoint, 1)
}
