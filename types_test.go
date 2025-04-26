package openaiorgs

import (
	"strings"
	"testing"
)

func TestListResponse_String(t *testing.T) {
	// Create a test ListResponse with some sample data
	response := &ListResponse[Project]{
		Object: "list",
		Data: []Project{
			{
				ID:     "proj_123",
				Name:   "Test Project 1",
				Status: "active",
			},
			{
				ID:     "proj_456",
				Name:   "Test Project 2",
				Status: "archived",
			},
		},
		FirstID: "proj_123",
		LastID:  "proj_456",
		HasMore: false,
	}

	// Get the string representation
	result := response.String()

	// Verify expected content
	expectedParts := []string{
		"Object: list",
		"First ID: proj_123",
		"Last ID: proj_456",
		"Has More: false",
		"Data:",
		"[0]",
		"proj_123",
		"Test Project 1",
		"[1]",
		"proj_456",
		"Test Project 2",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected string representation to contain %q, but it didn't.\nGot: %s", part, result)
		}
	}
}
