package openaiorgs

import (
	"net/http"
	"testing"
)

func TestClient_Get_WithPagination(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock first page response
	firstPageResponse := ListResponse[map[string]interface{}]{
		Object: "list",
		Data: []map[string]interface{}{
			{"id": "obj1", "name": "First"},
			{"id": "obj2", "name": "Second"},
		},
		FirstID: "obj1",
		LastID:  "obj2",
		HasMore: true,
	}

	h.mockResponse("GET", "/test-endpoint", http.StatusOK, firstPageResponse)

	// Mock second page response
	secondPageResponse := ListResponse[map[string]interface{}]{
		Object: "list",
		Data: []map[string]interface{}{
			{"id": "obj3", "name": "Third"},
			{"id": "obj4", "name": "Fourth"},
		},
		FirstID: "obj3",
		LastID:  "obj4",
		HasMore: false,
	}
	h.mockResponse("GET", "/test-endpoint?after=obj2", http.StatusOK, secondPageResponse)

	// First page request
	firstPage, err := Get[map[string]interface{}](h.client.client, "/test-endpoint", nil)
	if err != nil {
		t.Fatalf("Failed to get first page: %v", err)
	}

	// Verify first page
	if len(firstPage.Data) != 2 {
		t.Errorf("Expected 2 items in first page, got %d", len(firstPage.Data))
	}
	if !firstPage.HasMore {
		t.Error("Expected HasMore to be true for first page")
	}

	// Second page request with after parameter
	secondPage, err := Get[map[string]interface{}](h.client.client, "/test-endpoint", map[string]string{
		"after": firstPage.LastID,
	})
	if err != nil {
		t.Fatalf("Failed to get second page: %v", err)
	}

	// Verify second page
	if len(secondPage.Data) != 2 {
		t.Errorf("Expected 2 items in second page, got %d", len(secondPage.Data))
	}
	if secondPage.HasMore {
		t.Error("Expected HasMore to be false for second page")
	}

	// Verify total number of requests
	h.assertRequest("GET", "/test-endpoint", 2)
}
