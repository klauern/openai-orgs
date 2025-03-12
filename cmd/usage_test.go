package cmd

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestUsageCommand(t *testing.T) {
	cmd := UsageCommand()
	assert.Equal(t, "usage", cmd.Name)
	assert.Equal(t, 9, len(cmd.Commands)) // Should have 9 commands for different usage types
}

func TestBuildUsageQueryParams(t *testing.T) {
	app := &cli.Command{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "start-date"},
			&cli.StringFlag{Name: "end-date"},
			&cli.StringFlag{Name: "project-id"},
			&cli.IntFlag{Name: "limit"},
			&cli.StringFlag{Name: "after"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			params := buildUsageQueryParams(cmd)
			assert.Equal(t, "2023-01-01", params["start_date"])
			assert.Equal(t, "2023-01-31", params["end_date"])
			assert.Equal(t, "proj_123", params["project_id"])
			assert.Equal(t, "10", params["limit"])
			assert.Equal(t, "usage_abc", params["after"])
			return nil
		},
	}

	// Create CLI context with flags
	args := []string{
		"test",
		"--start-date", "2023-01-01",
		"--end-date", "2023-01-31",
		"--project-id", "proj_123",
		"--limit", "10",
		"--after", "usage_abc",
	}

	err := app.Run(context.Background(), args)
	assert.NoError(t, err)
}

func TestOutputUsageJSON(t *testing.T) {
	// Create a mock UsageResponse
	response := createTestUsageResponse()

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	err := outputUsageJSON(response, false)
	assert.NoError(t, err)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Validate output contains expected data
	assert.Contains(t, output, `"id": "usage_test_123"`)
	assert.Contains(t, output, `"cost": 0.05`)
	assert.Contains(t, output, `"project_id": "proj_test"`)
}

func TestOutputUsagePretty(t *testing.T) {
	// Create a mock UsageResponse
	response := createTestUsageResponse()

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function
	err := outputUsagePretty(response, true)
	assert.NoError(t, err)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Validate output contains expected data
	assert.Contains(t, output, "=== Usage Summary ===")
	assert.Contains(t, output, "Total records: 1")
	assert.Contains(t, output, "ID:        usage_test_123")
	assert.Contains(t, output, "Type:      completions")
	assert.Contains(t, output, "Cost:      $0.0500")
	assert.Contains(t, output, "Project:   proj_test")
}

// Helper function to create a test UsageResponse
func createTestUsageResponse() *openaiorgs.UsageResponse {
	return &openaiorgs.UsageResponse{
		Object: "list",
		Data: []openaiorgs.UsageRecord{
			{
				ID:     "usage_test_123",
				Object: "usage",
				Type:   openaiorgs.UsageTypeCompletions,
				UsageDetails: map[string]interface{}{
					"prompt_tokens":     float64(10),
					"completion_tokens": float64(20),
					"total_tokens":      float64(30),
					"model":             "gpt-4",
				},
				Cost:      0.05,
				ProjectID: "proj_test",
			},
		},
		FirstID: "usage_test_123",
		LastID:  "usage_test_123",
		HasMore: false,
	}
}
