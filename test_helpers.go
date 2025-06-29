package openaiorgs

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

// testBaseURL is the mock API endpoint used for testing.
const testBaseURL = "https://api.openai.com/v1"

// testHelper provides utilities for testing the OpenAI Organizations API client.
// It handles HTTP mocking, response registration, and test assertions.
type testHelper struct {
	// client is the OpenAI Organizations API client being tested.
	client *Client
	// t is the testing context for making assertions and reporting failures.
	t *testing.T
}

// newTestHelper creates a new test helper with a mocked HTTP client.
// It disables retries and activates HTTP mocking for the client.
// The returned helper should be used for a single test and cleaned up afterward.
func newTestHelper(t *testing.T) *testHelper {
	client := NewClient(testBaseURL, "test-token")
	// Disable retries for tests
	client.client.SetRetryCount(0)
	// Enable HTTP mocking
	httpmock.ActivateNonDefault(client.client.GetClient())

	return &testHelper{
		client: client,
		t:      t,
	}
}

// mockResponse registers a mock response for a given HTTP method and endpoint.
// It will return the specified response with the given status code when the endpoint is called.
//
// Parameters:
//   - method: HTTP method (GET, POST, etc.)
//   - endpoint: API endpoint path (without base URL)
//   - statusCode: HTTP status code to return
//   - response: Response body to return (will be JSON encoded)
func (h *testHelper) mockResponse(method, endpoint string, statusCode int, response any) {
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


// cleanup removes all registered HTTP mocks.
// This should be called after each test to ensure a clean state.
func (h *testHelper) cleanup() {
	httpmock.Reset()
}

// assertRequest verifies that a specific HTTP request was made the expected number of times.
// It checks both exact endpoint matches and endpoints with query parameters.
//
// Parameters:
//   - method: HTTP method to check (GET, POST, etc.)
//   - endpoint: API endpoint path to check (without base URL)
//   - times: Expected number of calls to this endpoint
//
// The test will fail if the actual number of calls doesn't match the expected count.
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
