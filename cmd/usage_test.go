package cmd

import (
	"fmt"
	"strings"
	"testing"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Mock client interface for usage testing
type mockUsageClient interface {
	GetCompletionsUsage(params map[string]string) (*openaiorgs.CompletionsUsageResponse, error)
	GetEmbeddingsUsage(params map[string]string) (*openaiorgs.EmbeddingsUsageResponse, error)
	GetModerationsUsage(params map[string]string) (*openaiorgs.ModerationsUsageResponse, error)
	GetImagesUsage(params map[string]string) (*openaiorgs.ImagesUsageResponse, error)
	GetAudioSpeechesUsage(params map[string]string) (*openaiorgs.AudioSpeechesUsageResponse, error)
	GetAudioTranscriptionsUsage(params map[string]string) (*openaiorgs.AudioTranscriptionsUsageResponse, error)
	GetVectorStoresUsage(params map[string]string) (*openaiorgs.VectorStoresUsageResponse, error)
	GetCodeInterpreterUsage(params map[string]string) (*openaiorgs.CodeInterpreterUsageResponse, error)
	GetCostsUsage(params map[string]string) (*openaiorgs.CostsUsageResponse, error)
}

// Mock implementation
type mockUsageClientImpl struct {
	GetCompletionsUsageFunc         func(params map[string]string) (*openaiorgs.CompletionsUsageResponse, error)
	GetEmbeddingsUsageFunc          func(params map[string]string) (*openaiorgs.EmbeddingsUsageResponse, error)
	GetModerationsUsageFunc         func(params map[string]string) (*openaiorgs.ModerationsUsageResponse, error)
	GetImagesUsageFunc              func(params map[string]string) (*openaiorgs.ImagesUsageResponse, error)
	GetAudioSpeechesUsageFunc       func(params map[string]string) (*openaiorgs.AudioSpeechesUsageResponse, error)
	GetAudioTranscriptionsUsageFunc func(params map[string]string) (*openaiorgs.AudioTranscriptionsUsageResponse, error)
	GetVectorStoresUsageFunc        func(params map[string]string) (*openaiorgs.VectorStoresUsageResponse, error)
	GetCodeInterpreterUsageFunc     func(params map[string]string) (*openaiorgs.CodeInterpreterUsageResponse, error)
	GetCostsUsageFunc               func(params map[string]string) (*openaiorgs.CostsUsageResponse, error)
}

func (m *mockUsageClientImpl) GetCompletionsUsage(params map[string]string) (*openaiorgs.CompletionsUsageResponse, error) {
	if m.GetCompletionsUsageFunc != nil {
		return m.GetCompletionsUsageFunc(params)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUsageClientImpl) GetEmbeddingsUsage(params map[string]string) (*openaiorgs.EmbeddingsUsageResponse, error) {
	if m.GetEmbeddingsUsageFunc != nil {
		return m.GetEmbeddingsUsageFunc(params)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUsageClientImpl) GetModerationsUsage(params map[string]string) (*openaiorgs.ModerationsUsageResponse, error) {
	if m.GetModerationsUsageFunc != nil {
		return m.GetModerationsUsageFunc(params)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUsageClientImpl) GetImagesUsage(params map[string]string) (*openaiorgs.ImagesUsageResponse, error) {
	if m.GetImagesUsageFunc != nil {
		return m.GetImagesUsageFunc(params)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUsageClientImpl) GetAudioSpeechesUsage(params map[string]string) (*openaiorgs.AudioSpeechesUsageResponse, error) {
	if m.GetAudioSpeechesUsageFunc != nil {
		return m.GetAudioSpeechesUsageFunc(params)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUsageClientImpl) GetAudioTranscriptionsUsage(params map[string]string) (*openaiorgs.AudioTranscriptionsUsageResponse, error) {
	if m.GetAudioTranscriptionsUsageFunc != nil {
		return m.GetAudioTranscriptionsUsageFunc(params)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUsageClientImpl) GetVectorStoresUsage(params map[string]string) (*openaiorgs.VectorStoresUsageResponse, error) {
	if m.GetVectorStoresUsageFunc != nil {
		return m.GetVectorStoresUsageFunc(params)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUsageClientImpl) GetCodeInterpreterUsage(params map[string]string) (*openaiorgs.CodeInterpreterUsageResponse, error) {
	if m.GetCodeInterpreterUsageFunc != nil {
		return m.GetCodeInterpreterUsageFunc(params)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUsageClientImpl) GetCostsUsage(params map[string]string) (*openaiorgs.CostsUsageResponse, error) {
	if m.GetCostsUsageFunc != nil {
		return m.GetCostsUsageFunc(params)
	}
	return nil, fmt.Errorf("not implemented")
}

// Helper to create a test CompletionsUsageResponse
func createTestCompletionsResponse(buckets ...openaiorgs.CompletionsUsageBucket) *openaiorgs.CompletionsUsageResponse {
	return &openaiorgs.CompletionsUsageResponse{
		Object:   "page",
		Data:     buckets,
		HasMore:  false,
		NextPage: "",
	}
}

func createTestCompletionsBucket(startTime, endTime int64, results ...openaiorgs.CompletionsUsageResult) openaiorgs.CompletionsUsageBucket {
	return openaiorgs.CompletionsUsageBucket{
		Object:    "bucket",
		StartTime: startTime,
		EndTime:   endTime,
		Results:   results,
	}
}

func createTestCompletionsResult(inputTokens, outputTokens, numRequests int) openaiorgs.CompletionsUsageResult {
	return openaiorgs.CompletionsUsageResult{
		Object:           "result",
		InputTokens:      inputTokens,
		OutputTokens:     outputTokens,
		NumModelRequests: numRequests,
		ProjectID:        "proj_test",
		Model:            "gpt-4",
	}
}

// Testable handler for completions usage
func getCompletionsUsageHandler(client mockUsageClient, params map[string]string, outputFormat string, verbose bool) error {
	usage, err := client.GetCompletionsUsage(params)
	if err != nil {
		return wrapError("get completions usage", err)
	}
	return outputCompletionsUsageResponse(usage, outputFormat, verbose)
}

func TestOutputCompletionsUsageJSON(t *testing.T) {
	result := createTestCompletionsResult(100, 50, 5)
	bucket := createTestCompletionsBucket(1700000000, 1700003600, result)
	response := createTestCompletionsResponse(bucket)

	output := captureOutput(func() {
		err := outputCompletionsUsageJSON(response, false)
		if err != nil {
			t.Errorf("outputCompletionsUsageJSON() error = %v", err)
		}
	})

	if !strings.Contains(output, `"object"`) {
		t.Errorf("Expected JSON to contain 'object' field, got: %s", output)
	}
	if !strings.Contains(output, `"input_tokens"`) {
		t.Errorf("Expected JSON to contain 'input_tokens', got: %s", output)
	}
	if !strings.Contains(output, `"output_tokens"`) {
		t.Errorf("Expected JSON to contain 'output_tokens', got: %s", output)
	}
	if !strings.Contains(output, `"start_time"`) {
		t.Errorf("Expected JSON to contain 'start_time', got: %s", output)
	}
}

func TestOutputCompletionsUsagePretty(t *testing.T) {
	t.Run("non-verbose", func(t *testing.T) {
		result := createTestCompletionsResult(100, 50, 5)
		bucket := createTestCompletionsBucket(1700000000, 1700003600, result)
		response := createTestCompletionsResponse(bucket)

		output := captureOutput(func() {
			err := outputCompletionsUsagePretty(response, false)
			if err != nil {
				t.Errorf("outputCompletionsUsagePretty() error = %v", err)
			}
		})

		if !strings.Contains(output, "=== Time Bucket ===") {
			t.Errorf("Expected time bucket header, got: %s", output)
		}
		if !strings.Contains(output, "--- Usage Record ---") {
			t.Errorf("Expected usage record header, got: %s", output)
		}
		if !strings.Contains(output, "Input tokens:        100") {
			t.Errorf("Expected input tokens count, got: %s", output)
		}
		if !strings.Contains(output, "Output tokens:       50") {
			t.Errorf("Expected output tokens count, got: %s", output)
		}
		if !strings.Contains(output, "Model requests:      5") {
			t.Errorf("Expected model requests count, got: %s", output)
		}
		if !strings.Contains(output, "Project ID:         proj_test") {
			t.Errorf("Expected project ID, got: %s", output)
		}
		if !strings.Contains(output, "Model:              gpt-4") {
			t.Errorf("Expected model name, got: %s", output)
		}
		if !strings.Contains(output, "Total records: 1") {
			t.Errorf("Expected total records count, got: %s", output)
		}
	})

	t.Run("verbose", func(t *testing.T) {
		result := createTestCompletionsResult(200, 100, 10)
		bucket := createTestCompletionsBucket(1700000000, 1700003600, result)
		response := &openaiorgs.CompletionsUsageResponse{
			Object:   "page",
			Data:     []openaiorgs.CompletionsUsageBucket{bucket},
			HasMore:  true,
			NextPage: "page_2",
		}

		output := captureOutput(func() {
			err := outputCompletionsUsagePretty(response, true)
			if err != nil {
				t.Errorf("outputCompletionsUsagePretty() error = %v", err)
			}
		})

		if !strings.Contains(output, "=== Completions Usage Summary ===") {
			t.Errorf("Expected summary header in verbose output, got: %s", output)
		}
		if !strings.Contains(output, "Total buckets: 1") {
			t.Errorf("Expected total buckets count, got: %s", output)
		}
		if !strings.Contains(output, "Has more: true") {
			t.Errorf("Expected has_more flag, got: %s", output)
		}
		if !strings.Contains(output, "Next page: page_2") {
			t.Errorf("Expected next page token, got: %s", output)
		}
	})
}

func TestOutputCompletionsUsageJSONL(t *testing.T) {
	t.Run("non-verbose", func(t *testing.T) {
		result := createTestCompletionsResult(100, 50, 5)
		bucket := createTestCompletionsBucket(1700000000, 1700003600, result)
		response := createTestCompletionsResponse(bucket)

		output := captureOutput(func() {
			err := outputCompletionsUsageJSONL(response, false)
			if err != nil {
				t.Errorf("outputCompletionsUsageJSONL() error = %v", err)
			}
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")
		// Should have: 1 bucket info line + 1 result line
		if len(lines) != 2 {
			t.Errorf("Expected 2 lines, got %d: %s", len(lines), output)
		}
		if !strings.Contains(lines[0], `"start_time"`) {
			t.Errorf("Expected first line to contain bucket info, got: %s", lines[0])
		}
		if !strings.Contains(lines[1], `"input_tokens"`) {
			t.Errorf("Expected second line to contain result, got: %s", lines[1])
		}
	})

	t.Run("verbose", func(t *testing.T) {
		result := createTestCompletionsResult(100, 50, 5)
		bucket := createTestCompletionsBucket(1700000000, 1700003600, result)
		response := createTestCompletionsResponse(bucket)

		output := captureOutput(func() {
			err := outputCompletionsUsageJSONL(response, true)
			if err != nil {
				t.Errorf("outputCompletionsUsageJSONL() error = %v", err)
			}
		})

		lines := strings.Split(strings.TrimSpace(output), "\n")
		// Should have: 1 metadata line + 1 bucket info line + 1 result line
		if len(lines) != 3 {
			t.Errorf("Expected 3 lines, got %d: %s", len(lines), output)
		}
		if !strings.Contains(lines[0], `"total"`) {
			t.Errorf("Expected metadata line first, got: %s", lines[0])
		}
	})
}

func TestGetCompletionsUsageHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockUsageClientImpl{
			GetCompletionsUsageFunc: func(params map[string]string) (*openaiorgs.CompletionsUsageResponse, error) {
				if params["start_time"] != "1700000000" {
					t.Errorf("unexpected start_time: %s", params["start_time"])
				}
				result := createTestCompletionsResult(100, 50, 5)
				bucket := createTestCompletionsBucket(1700000000, 1700003600, result)
				return createTestCompletionsResponse(bucket), nil
			},
		}

		output := captureOutput(func() {
			err := getCompletionsUsageHandler(mock, map[string]string{"start_time": "1700000000"}, "pretty", false)
			if err != nil {
				t.Errorf("getCompletionsUsageHandler() error = %v", err)
			}
		})

		if !strings.Contains(output, "=== Time Bucket ===") {
			t.Errorf("Expected time bucket header in output, got: %s", output)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock := &mockUsageClientImpl{
			GetCompletionsUsageFunc: func(params map[string]string) (*openaiorgs.CompletionsUsageResponse, error) {
				return nil, fmt.Errorf("API error: unauthorized")
			},
		}

		err := getCompletionsUsageHandler(mock, map[string]string{"start_time": "1700000000"}, "pretty", false)
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if !strings.Contains(err.Error(), "get completions usage") {
			t.Errorf("Expected wrapped error message, got: %v", err)
		}
	})

	t.Run("json output", func(t *testing.T) {
		mock := &mockUsageClientImpl{
			GetCompletionsUsageFunc: func(params map[string]string) (*openaiorgs.CompletionsUsageResponse, error) {
				result := createTestCompletionsResult(100, 50, 5)
				bucket := createTestCompletionsBucket(1700000000, 1700003600, result)
				return createTestCompletionsResponse(bucket), nil
			},
		}

		output := captureOutput(func() {
			err := getCompletionsUsageHandler(mock, map[string]string{"start_time": "1700000000"}, "json", false)
			if err != nil {
				t.Errorf("getCompletionsUsageHandler() error = %v", err)
			}
		})

		if !strings.Contains(output, `"input_tokens"`) {
			t.Errorf("Expected JSON output with input_tokens, got: %s", output)
		}
	})

	t.Run("unknown format", func(t *testing.T) {
		mock := &mockUsageClientImpl{
			GetCompletionsUsageFunc: func(params map[string]string) (*openaiorgs.CompletionsUsageResponse, error) {
				return createTestCompletionsResponse(), nil
			},
		}

		err := getCompletionsUsageHandler(mock, map[string]string{"start_time": "1700000000"}, "xml", false)
		if err == nil {
			t.Error("Expected error for unknown format, got nil")
		}
		if !strings.Contains(err.Error(), "unknown output format: xml") {
			t.Errorf("Expected descriptive error, got: %v", err)
		}
	})
}

func TestUsageOutputFormatRouting(t *testing.T) {
	result := createTestCompletionsResult(100, 50, 5)
	bucket := createTestCompletionsBucket(1700000000, 1700003600, result)
	response := createTestCompletionsResponse(bucket)

	tests := []struct {
		name          string
		format        string
		expectedInOut string
	}{
		{
			name:          "json routing",
			format:        "json",
			expectedInOut: `"input_tokens"`,
		},
		{
			name:          "jsonl routing",
			format:        "jsonl",
			expectedInOut: `"start_time"`,
		},
		{
			name:          "pretty routing",
			format:        "pretty",
			expectedInOut: "=== Time Bucket ===",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				err := outputCompletionsUsageResponse(response, tt.format, false)
				if err != nil {
					t.Errorf("outputCompletionsUsageResponse(%q) error = %v", tt.format, err)
				}
			})

			if !strings.Contains(output, tt.expectedInOut) {
				t.Errorf("Expected output to contain %q for format %q, got: %s", tt.expectedInOut, tt.format, output)
			}
		})
	}
}

func TestEmptyUsageResponse(t *testing.T) {
	response := createTestCompletionsResponse()

	t.Run("json empty", func(t *testing.T) {
		output := captureOutput(func() {
			err := outputCompletionsUsageJSON(response, false)
			if err != nil {
				t.Errorf("outputCompletionsUsageJSON() error = %v", err)
			}
		})
		if !strings.Contains(output, `"data"`) {
			t.Errorf("Expected 'data' field, got: %s", output)
		}
	})

	t.Run("pretty empty", func(t *testing.T) {
		output := captureOutput(func() {
			err := outputCompletionsUsagePretty(response, false)
			if err != nil {
				t.Errorf("outputCompletionsUsagePretty() error = %v", err)
			}
		})
		if !strings.Contains(output, "Total records: 0") {
			t.Errorf("Expected zero records in output, got: %s", output)
		}
		if strings.Contains(output, "=== Time Bucket ===") {
			t.Errorf("Expected no time buckets for empty data, got: %s", output)
		}
	})
}

func TestCostsUsageOutput(t *testing.T) {
	costsResponse := &openaiorgs.CostsUsageResponse{
		Object: "page",
		Data: []openaiorgs.CostsUsageBucket{
			{
				Object:    "bucket",
				StartTime: 1700000000,
				EndTime:   1700003600,
				Results: []openaiorgs.CostsUsageResult{
					{
						Object: "result",
						Amount: openaiorgs.CostAmount{
							Value:    12.50,
							Currency: "USD",
						},
						ProjectID: "proj_costs",
					},
				},
			},
		},
		HasMore: false,
	}

	t.Run("pretty", func(t *testing.T) {
		output := captureOutput(func() {
			err := outputCostsUsagePretty(costsResponse, false)
			if err != nil {
				t.Errorf("outputCostsUsagePretty() error = %v", err)
			}
		})

		if !strings.Contains(output, "Amount: 12.50 USD") {
			t.Errorf("Expected cost amount in output, got: %s", output)
		}
		if !strings.Contains(output, "Project ID: proj_costs") {
			t.Errorf("Expected project ID in output, got: %s", output)
		}
	})

	t.Run("json", func(t *testing.T) {
		output := captureOutput(func() {
			err := outputCostsUsageJSON(costsResponse, false)
			if err != nil {
				t.Errorf("outputCostsUsageJSON() error = %v", err)
			}
		})

		if !strings.Contains(output, `"value"`) {
			t.Errorf("Expected 'value' field in JSON, got: %s", output)
		}
		if !strings.Contains(output, `"currency"`) {
			t.Errorf("Expected 'currency' field in JSON, got: %s", output)
		}
	})
}

func TestEmbeddingsUsageOutput(t *testing.T) {
	response := &openaiorgs.EmbeddingsUsageResponse{
		Object: "page",
		Data: []openaiorgs.EmbeddingsUsageBucket{
			{
				Object:    "bucket",
				StartTime: 1700000000,
				EndTime:   1700003600,
				Results: []openaiorgs.EmbeddingsUsageResult{
					{
						Object:           "result",
						InputTokens:      500,
						NumModelRequests: 3,
						ProjectID:        "proj_embed",
						Model:            "text-embedding-ada-002",
					},
				},
			},
		},
		HasMore: false,
	}

	output := captureOutput(func() {
		err := outputEmbeddingsUsagePretty(response, false)
		if err != nil {
			t.Errorf("outputEmbeddingsUsagePretty() error = %v", err)
		}
	})

	if !strings.Contains(output, "Input tokens:   500") {
		t.Errorf("Expected input tokens in output, got: %s", output)
	}
	if !strings.Contains(output, "Model:         text-embedding-ada-002") {
		t.Errorf("Expected model name in output, got: %s", output)
	}
}

func TestImagesUsageOutput(t *testing.T) {
	response := &openaiorgs.ImagesUsageResponse{
		Object: "page",
		Data: []openaiorgs.ImagesUsageBucket{
			{
				Object:    "bucket",
				StartTime: 1700000000,
				EndTime:   1700003600,
				Results: []openaiorgs.ImagesUsageResult{
					{
						Object:           "result",
						Images:           10,
						NumModelRequests: 10,
						Size:             "1024x1024",
						Source:           "image.generations",
						ProjectID:        "proj_img",
						Model:            "dall-e-3",
					},
				},
			},
		},
		HasMore: false,
	}

	output := captureOutput(func() {
		err := outputImagesUsagePretty(response, false)
		if err != nil {
			t.Errorf("outputImagesUsagePretty() error = %v", err)
		}
	})

	if !strings.Contains(output, "Images:        10") {
		t.Errorf("Expected images count in output, got: %s", output)
	}
	if !strings.Contains(output, "Size:          1024x1024") {
		t.Errorf("Expected size in output, got: %s", output)
	}
}

func TestUsageResponseUnknownFormats(t *testing.T) {
	tests := []struct {
		name      string
		outputFn  func() error
		expectErr string
	}{
		{
			name: "completions unknown format",
			outputFn: func() error {
				return outputCompletionsUsageResponse(createTestCompletionsResponse(), "yaml", false)
			},
			expectErr: "unknown output format: yaml",
		},
		{
			name: "embeddings unknown format",
			outputFn: func() error {
				return outputEmbeddingsUsageResponse(&openaiorgs.EmbeddingsUsageResponse{}, "csv", false)
			},
			expectErr: "unknown output format: csv",
		},
		{
			name: "moderations unknown format",
			outputFn: func() error {
				return outputModerationsUsageResponse(&openaiorgs.ModerationsUsageResponse{}, "table", false)
			},
			expectErr: "unknown output format: table",
		},
		{
			name: "images unknown format",
			outputFn: func() error {
				return outputImagesUsageResponse(&openaiorgs.ImagesUsageResponse{}, "xml", false)
			},
			expectErr: "unknown output format: xml",
		},
		{
			name: "audio speeches unknown format",
			outputFn: func() error {
				return outputAudioSpeechesUsageResponse(&openaiorgs.AudioSpeechesUsageResponse{}, "html", false)
			},
			expectErr: "unknown output format: html",
		},
		{
			name: "audio transcriptions unknown format",
			outputFn: func() error {
				return outputAudioTranscriptionsUsageResponse(&openaiorgs.AudioTranscriptionsUsageResponse{}, "proto", false)
			},
			expectErr: "unknown output format: proto",
		},
		{
			name: "vector stores unknown format",
			outputFn: func() error {
				return outputVectorStoresUsageResponse(&openaiorgs.VectorStoresUsageResponse{}, "tsv", false)
			},
			expectErr: "unknown output format: tsv",
		},
		{
			name: "code interpreter unknown format",
			outputFn: func() error {
				return outputCodeInterpreterUsageResponse(&openaiorgs.CodeInterpreterUsageResponse{}, "toml", false)
			},
			expectErr: "unknown output format: toml",
		},
		{
			name: "costs unknown format",
			outputFn: func() error {
				return outputCostsUsageResponse(&openaiorgs.CostsUsageResponse{}, "ini", false)
			},
			expectErr: "unknown output format: ini",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.outputFn()
			if err == nil {
				t.Error("Expected error for unknown format, got nil")
			}
			if !strings.Contains(err.Error(), tt.expectErr) {
				t.Errorf("Expected error %q, got: %v", tt.expectErr, err)
			}
		})
	}
}
