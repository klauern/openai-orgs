package openaiorgs

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestGetCompletionsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for the new API structure
	compResult := CompletionsUsageResult{
		Object:            "completions_usage_result",
		InputTokens:       10,
		OutputTokens:      20,
		InputCachedTokens: 0,
		NumModelRequests:  1,
		ProjectID:         "proj_123",
		Model:             "gpt-4",
	}

	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	unixTimestamp := mockTimestamp.Unix()

	// Create response for the new completions usage API structure
	response := CompletionsUsageResponse{
		Object: "list",
		Data: []CompletionsUsageBucket{
			{
				Object:    "completions_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []CompletionsUsageResult{compResult},
			},
		},
		HasMore: false,
	}

	// Register mock responses
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

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in the bucket, got %d", len(usage.Data[0].Results))
		return
	}

	if usage.Data[0].Results[0].ProjectID != "proj_123" {
		t.Errorf("Expected ProjectID proj_123, got %s", usage.Data[0].Results[0].ProjectID)
	}

	if usage.Data[0].Results[0].InputTokens != 10 {
		t.Errorf("Expected InputTokens 10, got %d", usage.Data[0].Results[0].InputTokens)
	}

	if usage.Data[0].Results[0].OutputTokens != 20 {
		t.Errorf("Expected OutputTokens 20, got %d", usage.Data[0].Results[0].OutputTokens)
	}

	if usage.Data[0].Results[0].Model != "gpt-4" {
		t.Errorf("Expected Model gpt-4, got %s", usage.Data[0].Results[0].Model)
	}

	// Verify the first request was made
	h.assertRequest("GET", usageCompletionsEndpoint, 1)

	// Test with query parameters
	queryParams := map[string]string{
		"start_date": "2023-01-01",
		"end_date":   "2023-01-31",
	}
	_, err = h.client.GetCompletionsUsage(queryParams)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	// Verify the second request was made with query parameters
	h.assertRequest("GET", usageCompletionsEndpoint, 2)
}

func TestGetEmbeddingsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for the new API structure
	embResult := EmbeddingsUsageResult{
		Object:           "embeddings_usage_result",
		InputTokens:      50,
		NumModelRequests: 1,
		ProjectID:        "proj_456",
		Model:            "text-embedding-ada-002",
	}

	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	unixTimestamp := mockTimestamp.Unix()

	// Create response for the new embeddings usage API structure
	response := EmbeddingsUsageResponse{
		Object: "list",
		Data: []EmbeddingsUsageBucket{
			{
				Object:    "embeddings_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []EmbeddingsUsageResult{embResult},
			},
		},
		HasMore: false,
	}

	// Register mock responses
	h.mockResponse("GET", usageEmbeddingsEndpoint, 200, response)
	h.mockResponse("GET", usageEmbeddingsEndpoint, 200, response)

	// Test GetEmbeddingsUsage with no params
	usage, err := h.client.GetEmbeddingsUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in the bucket, got %d", len(usage.Data[0].Results))
		return
	}

	if usage.Data[0].Results[0].ProjectID != "proj_456" {
		t.Errorf("Expected ProjectID proj_456, got %s", usage.Data[0].Results[0].ProjectID)
	}

	if usage.Data[0].Results[0].InputTokens != 50 {
		t.Errorf("Expected InputTokens 50, got %d", usage.Data[0].Results[0].InputTokens)
	}

	if usage.Data[0].Results[0].Model != "text-embedding-ada-002" {
		t.Errorf("Expected Model text-embedding-ada-002, got %s", usage.Data[0].Results[0].Model)
	}

	// Verify the first request was made
	h.assertRequest("GET", usageEmbeddingsEndpoint, 1)

	// Test with query parameters
	queryParams := map[string]string{
		"start_date": "2023-01-01",
		"end_date":   "2023-01-31",
		"project_id": "proj_456",
	}
	_, err = h.client.GetEmbeddingsUsage(queryParams)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	// Verify the second request was made with query parameters
	h.assertRequest("GET", usageEmbeddingsEndpoint, 2)
}

func TestGetModerationsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for the new API structure
	modResult := ModerationsUsageResult{
		Object:           "moderations_usage_result",
		InputTokens:      25,
		NumModelRequests: 1,
		ProjectID:        "proj_789",
		Model:            "text-moderation-latest",
	}

	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	unixTimestamp := mockTimestamp.Unix()

	// Create response for the new moderations usage API structure
	response := ModerationsUsageResponse{
		Object: "list",
		Data: []ModerationsUsageBucket{
			{
				Object:    "moderations_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []ModerationsUsageResult{modResult},
			},
		},
		HasMore: false,
	}

	// Register mock response
	h.mockResponse("GET", usageModerationsEndpoint, 200, response)

	// Test GetModerationsUsage with no params
	usage, err := h.client.GetModerationsUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in the bucket, got %d", len(usage.Data[0].Results))
		return
	}

	if usage.Data[0].Results[0].ProjectID != "proj_789" {
		t.Errorf("Expected ProjectID proj_789, got %s", usage.Data[0].Results[0].ProjectID)
	}

	if usage.Data[0].Results[0].Model != "text-moderation-latest" {
		t.Errorf("Expected Model text-moderation-latest, got %s", usage.Data[0].Results[0].Model)
	}

	// Verify the request was made
	h.assertRequest("GET", usageModerationsEndpoint, 1)
}

func TestGetImagesUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for the new API structure
	imgResult := ImagesUsageResult{
		Object:           "images_usage_result",
		Images:           5,
		NumModelRequests: 1,
		Size:             "1024x1024",
		Source:           "dall-e-3",
		ProjectID:        "proj_987",
		Model:            "dall-e-3",
	}

	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	unixTimestamp := mockTimestamp.Unix()

	// Create response for the new images usage API structure
	response := ImagesUsageResponse{
		Object: "list",
		Data: []ImagesUsageBucket{
			{
				Object:    "images_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []ImagesUsageResult{imgResult},
			},
		},
		HasMore: false,
	}

	// Register mock response
	h.mockResponse("GET", usageImagesEndpoint, 200, response)

	// Test GetImagesUsage with no params
	usage, err := h.client.GetImagesUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in the bucket, got %d", len(usage.Data[0].Results))
		return
	}

	if usage.Data[0].Results[0].ProjectID != "proj_987" {
		t.Errorf("Expected ProjectID proj_987, got %s", usage.Data[0].Results[0].ProjectID)
	}

	if usage.Data[0].Results[0].Model != "dall-e-3" {
		t.Errorf("Expected Model dall-e-3, got %s", usage.Data[0].Results[0].Model)
	}

	// Verify the request was made
	h.assertRequest("GET", usageImagesEndpoint, 1)
}

func TestGetAudioSpeechesUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for the new API structure
	speechResult := AudioSpeechesUsageResult{
		Object:           "audio_speeches_usage_result",
		Characters:       1000,
		NumModelRequests: 1,
		ProjectID:        "proj_audio",
		Model:            "tts-1",
	}

	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	unixTimestamp := mockTimestamp.Unix()

	// Create response for the new audio speeches usage API structure
	response := AudioSpeechesUsageResponse{
		Object: "list",
		Data: []AudioSpeechesUsageBucket{
			{
				Object:    "audio_speeches_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []AudioSpeechesUsageResult{speechResult},
			},
		},
		HasMore: false,
	}

	// Register mock response
	h.mockResponse("GET", usageAudioSpeechesEndpoint, 200, response)

	// Test GetAudioSpeechesUsage with no params
	usage, err := h.client.GetAudioSpeechesUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in the bucket, got %d", len(usage.Data[0].Results))
		return
	}

	if usage.Data[0].Results[0].ProjectID != "proj_audio" {
		t.Errorf("Expected ProjectID proj_audio, got %s", usage.Data[0].Results[0].ProjectID)
	}

	if usage.Data[0].Results[0].Model != "tts-1" {
		t.Errorf("Expected Model tts-1, got %s", usage.Data[0].Results[0].Model)
	}

	// Verify the request was made
	h.assertRequest("GET", usageAudioSpeechesEndpoint, 1)
}

func TestGetAudioTranscriptionsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for the new API structure
	transResult := AudioTranscriptionsUsageResult{
		Object:           "audio_transcriptions_usage_result",
		Seconds:          120,
		NumModelRequests: 1,
		ProjectID:        "proj_trans",
		Model:            "whisper-1",
	}

	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	unixTimestamp := mockTimestamp.Unix()

	// Create response for the new audio transcriptions usage API structure
	response := AudioTranscriptionsUsageResponse{
		Object: "list",
		Data: []AudioTranscriptionsUsageBucket{
			{
				Object:    "audio_transcriptions_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []AudioTranscriptionsUsageResult{transResult},
			},
		},
		HasMore: false,
	}

	// Register mock response
	h.mockResponse("GET", usageAudioTranscriptionsEndpoint, 200, response)

	// Test GetAudioTranscriptionsUsage with no params
	usage, err := h.client.GetAudioTranscriptionsUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in the bucket, got %d", len(usage.Data[0].Results))
		return
	}

	if usage.Data[0].Results[0].ProjectID != "proj_trans" {
		t.Errorf("Expected ProjectID proj_trans, got %s", usage.Data[0].Results[0].ProjectID)
	}

	if usage.Data[0].Results[0].Model != "whisper-1" {
		t.Errorf("Expected Model whisper-1, got %s", usage.Data[0].Results[0].Model)
	}

	// Verify the request was made
	h.assertRequest("GET", usageAudioTranscriptionsEndpoint, 1)
}

func TestGetVectorStoresUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for the new API structure
	vectorResult := VectorStoresUsageResult{
		Object:     "vector_stores_usage_result",
		UsageBytes: 10000000,
		ProjectID:  "proj_vector",
	}

	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	unixTimestamp := mockTimestamp.Unix()

	// Create response for the new vector stores usage API structure
	response := VectorStoresUsageResponse{
		Object: "list",
		Data: []VectorStoresUsageBucket{
			{
				Object:    "vector_stores_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []VectorStoresUsageResult{vectorResult},
			},
		},
		HasMore: false,
	}

	// Register mock response
	h.mockResponse("GET", usageVectorStoresEndpoint, 200, response)

	// Test GetVectorStoresUsage with no params
	usage, err := h.client.GetVectorStoresUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in the bucket, got %d", len(usage.Data[0].Results))
		return
	}

	if usage.Data[0].Results[0].ProjectID != "proj_vector" {
		t.Errorf("Expected ProjectID proj_vector, got %s", usage.Data[0].Results[0].ProjectID)
	}

	if usage.Data[0].Results[0].UsageBytes != 10000000 {
		t.Errorf("Expected UsageBytes 10000000, got %d", usage.Data[0].Results[0].UsageBytes)
	}

	// Verify the request was made
	h.assertRequest("GET", usageVectorStoresEndpoint, 1)
}

func TestGetCodeInterpreterUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for the new API structure
	codeResult := CodeInterpreterUsageResult{
		Object:      "code_interpreter_usage_result",
		NumSessions: 3,
		ProjectID:   "proj_code",
	}

	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	unixTimestamp := mockTimestamp.Unix()

	// Create response for the new code interpreter usage API structure
	response := CodeInterpreterUsageResponse{
		Object: "list",
		Data: []CodeInterpreterUsageBucket{
			{
				Object:    "code_interpreter_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []CodeInterpreterUsageResult{codeResult},
			},
		},
		HasMore: false,
	}

	// Register mock response
	h.mockResponse("GET", usageCodeInterpreterEndpoint, 200, response)

	// Test GetCodeInterpreterUsage with no params
	usage, err := h.client.GetCodeInterpreterUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in the bucket, got %d", len(usage.Data[0].Results))
		return
	}

	if usage.Data[0].Results[0].ProjectID != "proj_code" {
		t.Errorf("Expected ProjectID proj_code, got %s", usage.Data[0].Results[0].ProjectID)
	}

	if usage.Data[0].Results[0].NumSessions != 3 {
		t.Errorf("Expected NumSessions 3, got %d", usage.Data[0].Results[0].NumSessions)
	}

	// Verify the request was made
	h.assertRequest("GET", usageCodeInterpreterEndpoint, 1)
}

func TestGetCostsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for the new API structure
	costResult := CostsUsageResult{
		Object: "costs_usage_result",
		Amount: CostAmount{
			Value:    150.25,
			Currency: "USD",
		},
		ProjectID: "proj_costs",
	}

	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	unixTimestamp := mockTimestamp.Unix()

	// Create response for the new costs usage API structure
	response := CostsUsageResponse{
		Object: "list",
		Data: []CostsUsageBucket{
			{
				Object:    "costs_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []CostsUsageResult{costResult},
			},
		},
		HasMore: false,
	}

	// Register mock response
	h.mockResponse("GET", usageCostsEndpoint, 200, response)

	// Test GetCostsUsage with no params
	usage, err := h.client.GetCostsUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in the bucket, got %d", len(usage.Data[0].Results))
		return
	}

	if usage.Data[0].Results[0].ProjectID != "proj_costs" {
		t.Errorf("Expected ProjectID proj_costs, got %s", usage.Data[0].Results[0].ProjectID)
	}

	if usage.Data[0].Results[0].Amount.Value != 150.25 {
		t.Errorf("Expected Amount.Value 150.25, got %f", usage.Data[0].Results[0].Amount.Value)
	}

	if usage.Data[0].Results[0].Amount.Currency != "USD" {
		t.Errorf("Expected Amount.Currency USD, got %s", usage.Data[0].Results[0].Amount.Currency)
	}

	// Verify the request was made
	h.assertRequest("GET", usageCostsEndpoint, 1)
}

func TestUsageWithPagination(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for the new API pagination structure
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	unixTimestamp := mockTimestamp.Unix()

	// Create first page response
	result1 := CompletionsUsageResult{
		Object:           "completions_usage_result",
		InputTokens:      10,
		OutputTokens:     20,
		NumModelRequests: 1,
		ProjectID:        "proj_123",
		Model:            "gpt-4",
	}

	result2 := CompletionsUsageResult{
		Object:           "completions_usage_result",
		InputTokens:      15,
		OutputTokens:     25,
		NumModelRequests: 1,
		ProjectID:        "proj_123",
		Model:            "gpt-4",
	}

	firstPageResponse := CompletionsUsageResponse{
		Object: "list",
		Data: []CompletionsUsageBucket{
			{
				Object:    "completions_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []CompletionsUsageResult{result1, result2},
			},
		},
		HasMore:  true,
		NextPage: "page2",
	}

	// Create second page response
	result3 := CompletionsUsageResult{
		Object:           "completions_usage_result",
		InputTokens:      30,
		OutputTokens:     40,
		NumModelRequests: 1,
		ProjectID:        "proj_123",
		Model:            "gpt-4",
	}

	secondPageResponse := CompletionsUsageResponse{
		Object: "list",
		Data: []CompletionsUsageBucket{
			{
				Object:    "completions_usage_bucket",
				StartTime: unixTimestamp,
				EndTime:   unixTimestamp + 3600,
				Results:   []CompletionsUsageResult{result3},
			},
		},
		HasMore:  false,
		NextPage: "",
	}

	// Create a map to store responses based on query parameters
	httpmock.Reset()
	httpmock.RegisterResponder("GET", testBaseURL+usageCompletionsEndpoint,
		func(req *http.Request) (*http.Response, error) {
			query := req.URL.Query()
			page := query.Get("page")

			if page == "page2" {
				// Return second page if "page" parameter is set to page2
				return httpmock.NewJsonResponse(200, secondPageResponse)
			}

			// Otherwise return first page
			return httpmock.NewJsonResponse(200, firstPageResponse)
		})

	// Test first page
	firstPage, err := h.client.GetCompletionsUsage(map[string]string{"limit": "2"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	// Validate first page response
	totalFirstPageResults := 0
	for _, bucket := range firstPage.Data {
		totalFirstPageResults += len(bucket.Results)
	}

	if totalFirstPageResults != 2 {
		t.Errorf("Expected 2 results in first page, got %d", totalFirstPageResults)
		return
	}

	if !firstPage.HasMore {
		t.Errorf("Expected HasMore to be true for first page")
		return
	}

	if firstPage.NextPage != "page2" {
		t.Errorf("Expected NextPage to be 'page2', got %s", firstPage.NextPage)
		return
	}

	// Test second page using page parameter with the NextPage from first page
	secondPage, err := h.client.GetCompletionsUsage(map[string]string{
		"limit": "2",
		"page":  firstPage.NextPage,
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	// Validate second page response
	totalSecondPageResults := 0
	for _, bucket := range secondPage.Data {
		totalSecondPageResults += len(bucket.Results)
	}

	if totalSecondPageResults != 1 {
		t.Errorf("Expected 1 result in second page, got %d", totalSecondPageResults)
		return
	}

	if secondPage.HasMore {
		t.Errorf("Expected HasMore to be false for second page")
		return
	}

	// Verify requests were made
	callCountInfo := httpmock.GetCallCountInfo()
	count := 0
	for key := range callCountInfo {
		if key[:len("GET "+testBaseURL+usageCompletionsEndpoint)] == "GET "+testBaseURL+usageCompletionsEndpoint {
			count += callCountInfo[key]
		}
	}

	if count != 2 {
		t.Errorf("Expected 2 calls to GET %s, got %d", usageCompletionsEndpoint, count)
	}
}

func TestUsageErrorHandling(t *testing.T) {
	endpoints := []string{
		usageCompletionsEndpoint,
		usageEmbeddingsEndpoint,
		usageModerationsEndpoint,
		usageImagesEndpoint,
		usageAudioSpeechesEndpoint,
		usageAudioTranscriptionsEndpoint,
		usageCodeInterpreterEndpoint,
		usageCostsEndpoint,
	}

	testCases := []struct {
		name      string
		setupMock func(h *testHelper, endpoint string)
		wantErr   string
	}{
		{
			name: "Network Error",
			setupMock: func(h *testHelper, endpoint string) {
				httpmock.RegisterResponder("GET", h.client.BaseURL+endpoint,
					func(req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("network error")
					})
			},
			wantErr: "network error",
		},
		{
			name: "Invalid JSON Response",
			setupMock: func(h *testHelper, endpoint string) {
				httpmock.RegisterResponder("GET", h.client.BaseURL+endpoint,
					func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:     "200 OK",
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(`{"invalid json`)),
							Header:     http.Header{"Content-Type": []string{"application/json"}},
							Request:    req,
						}, nil
					})
			},
			wantErr: "error unmarshaling response",
		},
		{
			name: "Unauthorized Error",
			setupMock: func(h *testHelper, endpoint string) {
				httpmock.RegisterResponder("GET", h.client.BaseURL+endpoint,
					httpmock.NewStringResponder(http.StatusUnauthorized, `{"error": "unauthorized"}`))
			},
			wantErr: "API request failed with status code 401",
		},
		{
			name: "Server Error",
			setupMock: func(h *testHelper, endpoint string) {
				httpmock.RegisterResponder("GET", h.client.BaseURL+endpoint,
					httpmock.NewStringResponder(http.StatusInternalServerError, `{"error": "internal server error"}`))
			},
			wantErr: "API request failed with status code 500",
		},
		{
			name: "Invalid Content Type",
			setupMock: func(h *testHelper, endpoint string) {
				httpmock.RegisterResponder("GET", h.client.BaseURL+endpoint,
					func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:     "200 OK",
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(`{"data": []}`)),
							Header:     http.Header{"Content-Type": []string{"text/plain"}},
							Request:    req,
						}, nil
					})
			},
			wantErr: "expected Content-Type \"application/json\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, endpoint := range endpoints {
				t.Run(endpoint, func(t *testing.T) {
					// Create a new test helper for each test case
					h := newTestHelper(t)
					defer h.cleanup()

					// Set up the mock for this test case
					tc.setupMock(h, endpoint)

					// Create query params
					queryParams := map[string]string{
						"start_time": "1234567890",
					}

					// Call the appropriate usage function based on the endpoint
					var err error
					switch endpoint {
					case usageCompletionsEndpoint:
						_, err = h.client.GetCompletionsUsage(queryParams)
					case usageEmbeddingsEndpoint:
						_, err = h.client.GetEmbeddingsUsage(queryParams)
					case usageModerationsEndpoint:
						_, err = h.client.GetModerationsUsage(queryParams)
					case usageImagesEndpoint:
						_, err = h.client.GetImagesUsage(queryParams)
					case usageAudioSpeechesEndpoint:
						_, err = h.client.GetAudioSpeechesUsage(queryParams)
					case usageAudioTranscriptionsEndpoint:
						_, err = h.client.GetAudioTranscriptionsUsage(queryParams)
					case usageCodeInterpreterEndpoint:
						_, err = h.client.GetCodeInterpreterUsage(queryParams)
					case usageCostsEndpoint:
						_, err = h.client.GetCostsUsage(queryParams)
					}

					// Verify the error
					if err == nil {
						t.Errorf("expected error containing %q, got nil", tc.wantErr)
						return
					}
					if !strings.Contains(err.Error(), tc.wantErr) {
						t.Errorf("expected error containing %q, got %q", tc.wantErr, err.Error())
					}
				})
			}
		})
	}
}

func TestUsageWithInvalidQueryParams(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Set up mock response for invalid query parameters
	h.mockResponse("GET", usageCompletionsEndpoint, http.StatusBadRequest, map[string]interface{}{
		"error": "invalid query parameters",
	})

	// Test with invalid query parameters
	queryParams := map[string]string{
		"invalid_param": "value",
	}

	_, err := h.client.GetCompletionsUsage(queryParams)
	if err == nil {
		t.Error("expected error for invalid query parameters, got nil")
		return
	}

	if !strings.Contains(err.Error(), "API request failed with status code 400") {
		t.Errorf("expected error containing 'API request failed with status code 400', got %q", err.Error())
	}
}

func TestUsageWithEmptyResponse(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Set up mock response with empty data
	h.mockResponse("GET", usageCompletionsEndpoint, http.StatusOK, map[string]interface{}{
		"object": "list",
		"data":   []interface{}{},
	})

	// Test with valid query parameters
	queryParams := map[string]string{
		"start_time": "1234567890",
	}

	resp, err := h.client.GetCompletionsUsage(queryParams)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if resp == nil {
		t.Error("expected non-nil response")
		return
	}

	if len(resp.Data) != 0 {
		t.Errorf("expected empty data array, got %d items", len(resp.Data))
	}
}

func TestUsageWithNilQueryParams(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Set up mock response for successful empty response
	h.mockResponse("GET", usageCompletionsEndpoint, http.StatusOK, map[string]interface{}{
		"object": "list",
		"data":   []interface{}{},
	})

	// Test with nil query parameters
	resp, err := h.client.GetCompletionsUsage(nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if resp == nil {
		t.Error("expected non-nil response")
		return
	}

	if len(resp.Data) != 0 {
		t.Errorf("expected empty data array, got %d items", len(resp.Data))
	}
}

func TestGetVectorStoresUsageErrors(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(h *testHelper)
		wantErr   string
	}{
		{
			name: "GET request error",
			setupMock: func(h *testHelper) {
				httpmock.RegisterResponder("GET", h.client.BaseURL+usageVectorStoresEndpoint,
					func(req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("network error")
					})
			},
			wantErr: "network error",
		},
		{
			name: "API error response",
			setupMock: func(h *testHelper) {
				httpmock.RegisterResponder("GET", h.client.BaseURL+usageVectorStoresEndpoint,
					httpmock.NewStringResponder(http.StatusInternalServerError, `{"error": "server error"}`))
			},
			wantErr: "API request failed with status code 500",
		},
		{
			name: "Invalid content type",
			setupMock: func(h *testHelper) {
				httpmock.RegisterResponder("GET", h.client.BaseURL+usageVectorStoresEndpoint,
					func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:     "200 OK",
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(`{"data": []}`)),
							Header:     http.Header{"Content-Type": []string{"text/plain"}},
							Request:    req,
						}, nil
					})
			},
			wantErr: "expected Content-Type \"application/json\"",
		},
		{
			name: "JSON unmarshal error",
			setupMock: func(h *testHelper) {
				httpmock.RegisterResponder("GET", h.client.BaseURL+usageVectorStoresEndpoint,
					func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							Status:     "200 OK",
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(`{"invalid json`)),
							Header:     http.Header{"Content-Type": []string{"application/json"}},
							Request:    req,
						}, nil
					})
			},
			wantErr: "error unmarshaling response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHelper(t)
			defer h.cleanup()

			tt.setupMock(h)

			_, err := h.client.GetVectorStoresUsage(map[string]string{"start_time": "1234567890"})
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
