package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Helper
func createMockProjectRateLimit(id, model string) openaiorgs.ProjectRateLimit {
	return openaiorgs.ProjectRateLimit{
		Object:                      "project.rate_limit",
		ID:                          id,
		Model:                       model,
		MaxRequestsPer1Minute:       100,
		MaxTokensPer1Minute:         50000,
		MaxImagesPer1Minute:         10,
		MaxAudioMegabytesPer1Minute: 5,
		MaxRequestsPer1Day:          10000,
		Batch1DayMaxInputTokens:     1000000,
	}
}

func TestListProjectRateLimitsTableCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful list table",
			args:       []string{"project-rate-limits", "list", "--project-id", "proj_123"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.ProjectRateLimit]{
				Object:  "list",
				Data:    []openaiorgs.ProjectRateLimit{createMockProjectRateLimit("rl_1", "gpt-4")},
				FirstID: "rl_1",
				LastID:  "rl_1",
				HasMore: false,
			},
			wantContains: []string{
				"ID | Model",
				"rl_1",
				"gpt-4",
				"100",
				"50000",
				"10",
				"5",
				"10000",
				"1000000",
			},
		},
		{
			name:       "empty list table",
			args:       []string{"project-rate-limits", "list", "--project-id", "proj_123"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.ProjectRateLimit]{
				Object: "list",
				Data:   []openaiorgs.ProjectRateLimit{},
			},
			wantContains: []string{"ID | Model"},
		},
		{
			name:       "error from API",
			args:       []string{"project-rate-limits", "list", "--project-id", "proj_123"},
			statusCode: 500,
			response:   map[string]string{"error": "API error"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/projects/proj_123/rate_limits", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectRateLimitsCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestListProjectRateLimitsJSONCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		statusCode int
		response   any
		wantErr   bool
	}{
		{
			name:       "successful list json",
			args:       []string{"--output", "json", "project-rate-limits", "list", "--project-id", "proj_123"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.ProjectRateLimit]{
				Object:  "list",
				Data:    []openaiorgs.ProjectRateLimit{createMockProjectRateLimit("rl_1", "gpt-4")},
				FirstID: "rl_1",
				LastID:  "rl_1",
				HasMore: false,
			},
		},
		{
			name:       "error from API json",
			args:       []string{"--output", "json", "project-rate-limits", "list", "--project-id", "proj_123"},
			statusCode: 500,
			response:   map[string]string{"error": "API error"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/projects/proj_123/rate_limits", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectRateLimitsCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify valid JSON array output
				var result []map[string]interface{}
				if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &result); err != nil {
					t.Errorf("Expected valid JSON array output, got: %s", output)
				}
				if !strings.Contains(output, "rl_1") {
					t.Errorf("Expected output to contain rl_1, got: %s", output)
				}
				if !strings.Contains(output, "gpt-4") {
					t.Errorf("Expected output to contain gpt-4, got: %s", output)
				}
			}
		})
	}
}

func TestModifyProjectRateLimitTableCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name: "successful modify table",
			args: []string{
				"project-rate-limits", "modify",
				"--project-id", "proj_123",
				"--rate-limit-id", "rl_1",
				"--max-requests-per-1-minute", "200",
				"--max-tokens-per-1-minute", "100000",
			},
			statusCode: 200,
			response: func() openaiorgs.ProjectRateLimit {
				rl := createMockProjectRateLimit("rl_1", "gpt-4")
				rl.MaxRequestsPer1Minute = 200
				rl.MaxTokensPer1Minute = 100000
				return rl
			}(),
			wantContains: []string{"rl_1", "gpt-4", "200", "100000"},
		},
		{
			name: "error from API",
			args: []string{
				"project-rate-limits", "modify",
				"--project-id", "proj_123",
				"--rate-limit-id", "rl_1",
				"--max-requests-per-1-minute", "200",
			},
			statusCode: 500,
			response:   map[string]string{"error": "modify failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("POST", "/organization/projects/proj_123/rate_limits/rl_1", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectRateLimitsCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestModifyProjectRateLimitJSONCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		statusCode int
		response   any
		wantErr   bool
	}{
		{
			name: "successful modify json",
			args: []string{
				"--output", "json",
				"project-rate-limits", "modify",
				"--project-id", "proj_123",
				"--rate-limit-id", "rl_1",
				"--max-requests-per-1-minute", "200",
				"--max-tokens-per-1-minute", "100000",
			},
			statusCode: 200,
			response: func() openaiorgs.ProjectRateLimit {
				rl := createMockProjectRateLimit("rl_1", "gpt-4")
				rl.MaxRequestsPer1Minute = 200
				rl.MaxTokensPer1Minute = 100000
				return rl
			}(),
		},
		{
			name: "error from API json",
			args: []string{
				"--output", "json",
				"project-rate-limits", "modify",
				"--project-id", "proj_123",
				"--rate-limit-id", "rl_1",
				"--max-requests-per-1-minute", "200",
			},
			statusCode: 500,
			response:   map[string]string{"error": "modify failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("POST", "/organization/projects/proj_123/rate_limits/rl_1", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectRateLimitsCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify valid JSON output
				var result map[string]interface{}
				if err := json.Unmarshal([]byte(strings.TrimSpace(output)), &result); err != nil {
					t.Errorf("Expected valid JSON output, got: %s", output)
				}
				if !strings.Contains(output, "rl_1") {
					t.Errorf("Expected output to contain rl_1, got: %s", output)
				}
				if !strings.Contains(output, "gpt-4") {
					t.Errorf("Expected output to contain gpt-4, got: %s", output)
				}
			}
		})
	}
}
