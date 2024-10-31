package openaiorgs

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

const testBaseURL = "https://api.openai.com/v1"

type testHelper struct {
	client *Client
	t      *testing.T
}

// newTestHelper creates a new test helper with mocked HTTP client
func newTestHelper(t *testing.T) *testHelper {
	client := NewClient(testBaseURL, "test-token")
	// Enable HTTP mocking
	httpmock.ActivateNonDefault(client.client.GetClient())

	return &testHelper{
		client: client,
		t:      t,
	}
}

// mockResponse registers a mock response for a given method, endpoint, and response
func (h *testHelper) mockResponse(method, endpoint string, statusCode int, response interface{}) {
	responder := func(req *http.Request) (*http.Response, error) {
		// Return the response directly without any conditions
		resp, err := httpmock.NewJsonResponse(statusCode, response)
		if err != nil {
			h.t.Fatalf("Failed to create mock response: %v", err)
		}
		return resp, nil
	}

	// Register the mock responder
	httpmock.RegisterResponder(method, testBaseURL+endpoint, responder)
}

// mockListResponse is a helper for mocking paginated list responses
func (h *testHelper) mockListResponse(method, endpoint string, items interface{}) { //nolint:unused
	response := ListResponse[interface{}]{
		Object:  "list",
		Data:    []interface{}{items},
		FirstID: "first_id",
		LastID:  "last_id",
		HasMore: false,
	}
	h.mockResponse(method, endpoint, http.StatusOK, response)
}

// cleanup removes all registered mocks
func (h *testHelper) cleanup() {
	httpmock.Reset()
}

// assertRequest verifies that a specific request was made
func (h *testHelper) assertRequest(method, endpoint string, times int) {
	// Original code only checks exact endpoint match
	count := httpmock.GetCallCountInfo()[method+" "+testBaseURL+endpoint]

	// Need to also check for endpoint with query parameters
	if times > count {
		// Check if there are any calls with additional query parameters
		for key := range httpmock.GetCallCountInfo() {
			if key != method+" "+testBaseURL+endpoint &&
				key[:len(method+" "+testBaseURL+endpoint)] == method+" "+testBaseURL+endpoint {
				count += httpmock.GetCallCountInfo()[key]
			}
		}
	}

	if count != times {
		h.t.Errorf("Expected %d calls to %s %s, got %d", times, method, endpoint, count)
	}
}

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
