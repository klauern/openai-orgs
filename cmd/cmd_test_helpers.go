package cmd

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/urfave/cli/v3"

	openaiorgs "github.com/klauern/openai-orgs"
)

const testBaseURL = "https://api.openai.com/v1"

// cmdTestHelper provides utilities for testing CLI commands with httpmock.
// It intercepts newClient to return an httpmock-backed client, so that
// actual CLI action functions are exercised end-to-end.
type cmdTestHelper struct {
	client *openaiorgs.Client
	t      *testing.T
}

// newCmdTestHelper creates a test helper that overrides newClientFunc
// to return a pre-configured client with httpmock enabled.
// Call cleanup() via defer to restore the original newClientFunc and reset mocks.
func newCmdTestHelper(t *testing.T) *cmdTestHelper {
	t.Helper()
	client := openaiorgs.NewClient(testBaseURL, "test-token")
	// Disable retries for tests
	client.SetRetryCount(0)
	// Enable HTTP mocking
	httpmock.ActivateNonDefault(client.GetHTTPClient())

	// Override the client factory so CLI actions use our mocked client
	newClientFunc = func(_ context.Context, _ *cli.Command) *openaiorgs.Client {
		return client
	}

	return &cmdTestHelper{
		client: client,
		t:      t,
	}
}

// mockResponse registers a mock JSON response for a given HTTP method and endpoint.
func (h *cmdTestHelper) mockResponse(method, endpoint string, statusCode int, response any) {
	h.t.Helper()
	responder := func(req *http.Request) (*http.Response, error) {
		resp, err := httpmock.NewJsonResponse(statusCode, response)
		if err != nil {
			h.t.Fatalf("Failed to create mock response: %v", err)
		}
		return resp, nil
	}
	httpmock.RegisterResponder(method, testBaseURL+endpoint, responder)
}

// cleanup restores the original newClientFunc and resets all httpmock state.
func (h *cmdTestHelper) cleanup() {
	newClientFunc = defaultNewClient
	httpmock.DeactivateAndReset()
}

// assertRequest verifies that a specific HTTP request was made the expected number of times.
func (h *cmdTestHelper) assertRequest(method, endpoint string, times int) {
	h.t.Helper()
	fullURL := method + " " + testBaseURL + endpoint
	count := httpmock.GetCallCountInfo()[fullURL]

	// Also check for calls with query parameters
	if times > count {
		for key, c := range httpmock.GetCallCountInfo() {
			if key != fullURL && strings.HasPrefix(key, fullURL+"?") {
				count += c
			}
		}
	}

	if count != times {
		h.t.Errorf("Expected %d calls to %s %s, got %d", times, method, endpoint, count)
	}
}

// runCmd runs a CLI command with the given arguments using the test helper's context.
// It builds a minimal urfave/cli app with the provided command and runs it.
// The root command includes global flags (output, api-key) matching the real app.
func (h *cmdTestHelper) runCmd(command *cli.Command, args []string) error {
	h.t.Helper()
	root := &cli.Command{
		Name: "test",
		Commands: []*cli.Command{
			command,
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output",
				Usage: "Output format (default: pretty)",
				Value: "pretty",
			},
			&cli.StringFlag{
				Name:  "api-key",
				Usage: "OpenAI API key",
				Value: "test-token",
			},
		},
	}
	return root.Run(context.Background(), append([]string{"test"}, args...))
}

// captureOutput captures stdout during the execution of f and returns it as a string.
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return ""
	}
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String()
}
