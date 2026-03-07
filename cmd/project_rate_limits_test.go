package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Mock client interface for testing project rate limits
type mockProjectRateLimitClient interface {
	ListProjectRateLimits(limit int, after string, projectID string) (*openaiorgs.ListResponse[openaiorgs.ProjectRateLimit], error)
	ModifyProjectRateLimit(projectID, rateLimitID string, fields openaiorgs.ProjectRateLimitRequestFields) (*openaiorgs.ProjectRateLimit, error)
}

// Mock implementation
type mockProjectRateLimitClientImpl struct {
	ListProjectRateLimitsFunc  func(limit int, after string, projectID string) (*openaiorgs.ListResponse[openaiorgs.ProjectRateLimit], error)
	ModifyProjectRateLimitFunc func(projectID, rateLimitID string, fields openaiorgs.ProjectRateLimitRequestFields) (*openaiorgs.ProjectRateLimit, error)
}

func (m *mockProjectRateLimitClientImpl) ListProjectRateLimits(limit int, after string, projectID string) (*openaiorgs.ListResponse[openaiorgs.ProjectRateLimit], error) {
	if m.ListProjectRateLimitsFunc != nil {
		return m.ListProjectRateLimitsFunc(limit, after, projectID)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectRateLimitClientImpl) ModifyProjectRateLimit(projectID, rateLimitID string, fields openaiorgs.ProjectRateLimitRequestFields) (*openaiorgs.ProjectRateLimit, error) {
	if m.ModifyProjectRateLimitFunc != nil {
		return m.ModifyProjectRateLimitFunc(projectID, rateLimitID, fields)
	}
	return nil, fmt.Errorf("not implemented")
}

// Testable handlers

func listProjectRateLimitsTableHandler(client mockProjectRateLimitClient, limit int, after, projectID string) error {
	rateLimits, err := client.ListProjectRateLimits(limit, after, projectID)
	if err != nil {
		return wrapError("list project rate limits", err)
	}
	return printProjectRateLimitsTable(rateLimits)
}

func listProjectRateLimitsJSONHandler(client mockProjectRateLimitClient, limit int, after, projectID string) error {
	rateLimits, err := client.ListProjectRateLimits(limit, after, projectID)
	if err != nil {
		return wrapError("list project rate limits", err)
	}
	return printProjectRateLimitsJSON(rateLimits)
}

func modifyProjectRateLimitTableHandler(client mockProjectRateLimitClient, projectID, rateLimitID string, fields openaiorgs.ProjectRateLimitRequestFields) error {
	rateLimit, err := client.ModifyProjectRateLimit(projectID, rateLimitID, fields)
	if err != nil {
		return wrapError("modify project rate limit", err)
	}
	return printProjectRateLimitTable(rateLimit)
}

func modifyProjectRateLimitJSONHandler(client mockProjectRateLimitClient, projectID, rateLimitID string, fields openaiorgs.ProjectRateLimitRequestFields) error {
	rateLimit, err := client.ModifyProjectRateLimit(projectID, rateLimitID, fields)
	if err != nil {
		return wrapError("modify project rate limit", err)
	}
	return printProjectRateLimitJSON(rateLimit)
}

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

// Tests

func TestListProjectRateLimitsTableHandler(t *testing.T) {
	tests := []struct {
		name         string
		limit        int
		after        string
		projectID    string
		mockFn       func(*mockProjectRateLimitClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:      "successful list table",
			limit:     10,
			after:     "",
			projectID: "proj_123",
			mockFn: func(m *mockProjectRateLimitClientImpl) {
				rl := createMockProjectRateLimit("rl_1", "gpt-4")
				m.ListProjectRateLimitsFunc = func(limit int, after string, projectID string) (*openaiorgs.ListResponse[openaiorgs.ProjectRateLimit], error) {
					if projectID != "proj_123" {
						t.Errorf("unexpected projectID: %s", projectID)
					}
					return &openaiorgs.ListResponse[openaiorgs.ProjectRateLimit]{
						Object:  "list",
						Data:    []openaiorgs.ProjectRateLimit{rl},
						FirstID: "rl_1",
						LastID:  "rl_1",
						HasMore: false,
					}, nil
				}
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
			name:      "empty list table",
			limit:     10,
			after:     "",
			projectID: "proj_123",
			mockFn: func(m *mockProjectRateLimitClientImpl) {
				m.ListProjectRateLimitsFunc = func(limit int, after string, projectID string) (*openaiorgs.ListResponse[openaiorgs.ProjectRateLimit], error) {
					return &openaiorgs.ListResponse[openaiorgs.ProjectRateLimit]{
						Object: "list",
						Data:   []openaiorgs.ProjectRateLimit{},
					}, nil
				}
			},
			wantContains: []string{"ID | Model"},
		},
		{
			name:      "error from client",
			limit:     10,
			after:     "",
			projectID: "proj_123",
			mockFn: func(m *mockProjectRateLimitClientImpl) {
				m.ListProjectRateLimitsFunc = func(limit int, after string, projectID string) (*openaiorgs.ListResponse[openaiorgs.ProjectRateLimit], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectRateLimitClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listProjectRateLimitsTableHandler(mock, tt.limit, tt.after, tt.projectID)
				if (err != nil) != tt.wantErr {
					t.Errorf("listProjectRateLimitsTableHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestListProjectRateLimitsJSONHandler(t *testing.T) {
	tests := []struct {
		name      string
		limit     int
		after     string
		projectID string
		mockFn    func(*mockProjectRateLimitClientImpl)
		wantErr   bool
	}{
		{
			name:      "successful list json",
			limit:     10,
			after:     "",
			projectID: "proj_123",
			mockFn: func(m *mockProjectRateLimitClientImpl) {
				rl := createMockProjectRateLimit("rl_1", "gpt-4")
				m.ListProjectRateLimitsFunc = func(limit int, after string, projectID string) (*openaiorgs.ListResponse[openaiorgs.ProjectRateLimit], error) {
					return &openaiorgs.ListResponse[openaiorgs.ProjectRateLimit]{
						Object:  "list",
						Data:    []openaiorgs.ProjectRateLimit{rl},
						FirstID: "rl_1",
						LastID:  "rl_1",
						HasMore: false,
					}, nil
				}
			},
		},
		{
			name:      "error from client json",
			limit:     10,
			after:     "",
			projectID: "proj_123",
			mockFn: func(m *mockProjectRateLimitClientImpl) {
				m.ListProjectRateLimitsFunc = func(limit int, after string, projectID string) (*openaiorgs.ListResponse[openaiorgs.ProjectRateLimit], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectRateLimitClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listProjectRateLimitsJSONHandler(mock, tt.limit, tt.after, tt.projectID)
				if (err != nil) != tt.wantErr {
					t.Errorf("listProjectRateLimitsJSONHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

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

func TestModifyProjectRateLimitTableHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectID    string
		rateLimitID  string
		fields       openaiorgs.ProjectRateLimitRequestFields
		mockFn       func(*mockProjectRateLimitClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:        "successful modify table",
			projectID:   "proj_123",
			rateLimitID: "rl_1",
			fields: openaiorgs.ProjectRateLimitRequestFields{
				MaxRequestsPer1Minute: 200,
				MaxTokensPer1Minute:   100000,
			},
			mockFn: func(m *mockProjectRateLimitClientImpl) {
				m.ModifyProjectRateLimitFunc = func(projectID, rateLimitID string, fields openaiorgs.ProjectRateLimitRequestFields) (*openaiorgs.ProjectRateLimit, error) {
					if projectID != "proj_123" || rateLimitID != "rl_1" {
						t.Errorf("unexpected params: projectID=%s, rateLimitID=%s", projectID, rateLimitID)
					}
					if fields.MaxRequestsPer1Minute != 200 || fields.MaxTokensPer1Minute != 100000 {
						t.Errorf("unexpected fields: %+v", fields)
					}
					rl := createMockProjectRateLimit("rl_1", "gpt-4")
					rl.MaxRequestsPer1Minute = 200
					rl.MaxTokensPer1Minute = 100000
					return &rl, nil
				}
			},
			wantContains: []string{"rl_1", "gpt-4", "200", "100000"},
		},
		{
			name:        "error from client",
			projectID:   "proj_123",
			rateLimitID: "rl_1",
			fields: openaiorgs.ProjectRateLimitRequestFields{
				MaxRequestsPer1Minute: 200,
			},
			mockFn: func(m *mockProjectRateLimitClientImpl) {
				m.ModifyProjectRateLimitFunc = func(projectID, rateLimitID string, fields openaiorgs.ProjectRateLimitRequestFields) (*openaiorgs.ProjectRateLimit, error) {
					return nil, fmt.Errorf("modify failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectRateLimitClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := modifyProjectRateLimitTableHandler(mock, tt.projectID, tt.rateLimitID, tt.fields)
				if (err != nil) != tt.wantErr {
					t.Errorf("modifyProjectRateLimitTableHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestModifyProjectRateLimitJSONHandler(t *testing.T) {
	tests := []struct {
		name        string
		projectID   string
		rateLimitID string
		fields      openaiorgs.ProjectRateLimitRequestFields
		mockFn      func(*mockProjectRateLimitClientImpl)
		wantErr     bool
	}{
		{
			name:        "successful modify json",
			projectID:   "proj_123",
			rateLimitID: "rl_1",
			fields: openaiorgs.ProjectRateLimitRequestFields{
				MaxRequestsPer1Minute: 200,
				MaxTokensPer1Minute:   100000,
			},
			mockFn: func(m *mockProjectRateLimitClientImpl) {
				m.ModifyProjectRateLimitFunc = func(projectID, rateLimitID string, fields openaiorgs.ProjectRateLimitRequestFields) (*openaiorgs.ProjectRateLimit, error) {
					rl := createMockProjectRateLimit("rl_1", "gpt-4")
					rl.MaxRequestsPer1Minute = 200
					rl.MaxTokensPer1Minute = 100000
					return &rl, nil
				}
			},
		},
		{
			name:        "error from client json",
			projectID:   "proj_123",
			rateLimitID: "rl_1",
			fields: openaiorgs.ProjectRateLimitRequestFields{
				MaxRequestsPer1Minute: 200,
			},
			mockFn: func(m *mockProjectRateLimitClientImpl) {
				m.ModifyProjectRateLimitFunc = func(projectID, rateLimitID string, fields openaiorgs.ProjectRateLimitRequestFields) (*openaiorgs.ProjectRateLimit, error) {
					return nil, fmt.Errorf("modify failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectRateLimitClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := modifyProjectRateLimitJSONHandler(mock, tt.projectID, tt.rateLimitID, tt.fields)
				if (err != nil) != tt.wantErr {
					t.Errorf("modifyProjectRateLimitJSONHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

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
