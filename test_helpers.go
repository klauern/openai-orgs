package openaiorgs

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

const testBaseURL = "https://api.openai.com/v1"

type testHelper struct {
	client *Client
	t      *testing.T
}

// newTestHelper creates a new test helper with mocked HTTP client
func newTestHelper(t *testing.T) *testHelper {
	client := resty.New()
	// Disable retries for tests
	client.SetRetryCount(0)
	client.SetBaseURL(testBaseURL)
	client.SetAuthToken("test-token")
	client.SetHeader("Content-Type", "application/json")

	c := &Client{
		client:  client,
		BaseURL: testBaseURL,
	}

	// Enable HTTP mocking
	httpmock.ActivateNonDefault(c.client.GetClient())

	return &testHelper{
		client: c,
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
