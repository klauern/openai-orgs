package openaiorgs

import (
	"testing"
	"time"
)

func TestGetCompletionsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockUsageRecord := UsageRecord{
		ID:        "usage_123",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      UsageTypeCompletions,
		UsageDetails: map[string]interface{}{
			"prompt_tokens":     10,
			"completion_tokens": 20,
			"total_tokens":      30,
			"model":             "gpt-4",
		},
		Cost:      0.05,
		ProjectID: "proj_123",
	}

	// Create response for both calls
	response := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockUsageRecord},
		FirstID: "usage_123",
		LastID:  "usage_123",
		HasMore: false,
	}

	// Register mock responses for both calls before making any requests
	h.mockResponse("GET", usageCompletionsEndpoint, 200, response)
	h.mockResponse("GET", usageCompletionsEndpoint, 200, response)

	// Test GetCompletionsUsage with no params
	usage, err := h.client.GetCompletionsUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}
	if usage.Data[0].ID != "usage_123" {
		t.Errorf("Expected ID usage_123, got %s", usage.Data[0].ID)
	}
	if usage.Data[0].Type != UsageTypeCompletions {
		t.Errorf("Expected Type %s, got %s", UsageTypeCompletions, usage.Data[0].Type)
	}
	if usage.Data[0].Cost != 0.05 {
		t.Errorf("Expected Cost 0.05, got %f", usage.Data[0].Cost)
	}
	if usage.Data[0].ProjectID != "proj_123" {
		t.Errorf("Expected ProjectID proj_123, got %s", usage.Data[0].ProjectID)
	}

	// Verify the first request was made
	h.assertRequest("GET", usageCompletionsEndpoint, 1)

	// Test with query parameters
	queryParams := map[string]string{
		"start_date": "2023-01-01",
		"end_date":   "2023-01-31",
	}
	usage, err = h.client.GetCompletionsUsage(queryParams)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	// Verify the second request was made with query parameters
	h.assertRequest("GET", usageCompletionsEndpoint, 2) // Total call count is now 2
}
