package openaiorgs

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestGetCompletionsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResults := []CompletionsUsageResult{
		{
			Object:           "usage_result",
			InputTokens:      10,
			OutputTokens:     20,
			NumModelRequests: 1,
			ProjectID:        "proj_123",
			Model:            "gpt-4",
		},
	}

	mockBucket := CompletionsUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672531200, // 2023-01-01 00:00:00 UTC
		EndTime:   1672617600, // 2023-01-02 00:00:00 UTC
		Results:   mockResults,
	}

	response := CompletionsUsageResponse{
		Object:   "list",
		Data:     []CompletionsUsageBucket{mockBucket},
		HasMore:  false,
		NextPage: "",
	}

	// Register mock response
	h.mockResponse("GET", usageCompletionsEndpoint, 200, response)

	// Test GetCompletionsUsage with no params
	usage, err := h.client.GetCompletionsUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage bucket, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in bucket, got %d", len(usage.Data[0].Results))
		return
	}

	result := usage.Data[0].Results[0]
	if result.InputTokens != 10 {
		t.Errorf("Expected InputTokens 10, got %d", result.InputTokens)
	}
	if result.OutputTokens != 20 {
		t.Errorf("Expected OutputTokens 20, got %d", result.OutputTokens)
	}
	if result.ProjectID != "proj_123" {
		t.Errorf("Expected ProjectID proj_123, got %s", result.ProjectID)
	}
	if result.Model != "gpt-4" {
		t.Errorf("Expected Model gpt-4, got %s", result.Model)
	}

	// Verify the first request was made
	h.assertRequest("GET", usageCompletionsEndpoint, 1)

	// Test with query parameters
	queryParams := map[string]string{
		"start_time": "1672531200",
		"end_time":   "1672617600",
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

	// Mock response data
	mockResults := []EmbeddingsUsageResult{
		{
			Object:           "usage_result",
			InputTokens:      50,
			NumModelRequests: 1,
			ProjectID:        "proj_456",
			Model:            "text-embedding-ada-002",
		},
	}

	mockBucket := EmbeddingsUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672531200, // 2023-01-01 00:00:00 UTC
		EndTime:   1672617600, // 2023-01-02 00:00:00 UTC
		Results:   mockResults,
	}

	response := EmbeddingsUsageResponse{
		Object:   "list",
		Data:     []EmbeddingsUsageBucket{mockBucket},
		HasMore:  false,
		NextPage: "",
	}

	// Register mock response
	h.mockResponse("GET", usageEmbeddingsEndpoint, 200, response)

	// Test GetEmbeddingsUsage with no params
	usage, err := h.client.GetEmbeddingsUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage bucket, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in bucket, got %d", len(usage.Data[0].Results))
		return
	}

	result := usage.Data[0].Results[0]
	if result.InputTokens != 50 {
		t.Errorf("Expected InputTokens 50, got %d", result.InputTokens)
	}
	if result.ProjectID != "proj_456" {
		t.Errorf("Expected ProjectID proj_456, got %s", result.ProjectID)
	}
	if result.Model != "text-embedding-ada-002" {
		t.Errorf("Expected Model text-embedding-ada-002, got %s", result.Model)
	}

	// Test with query parameters
	queryParams := map[string]string{
		"start_time": "1672531200",
		"end_time":   "1672617600",
		"project_id": "proj_456",
	}
	_, err = h.client.GetEmbeddingsUsage(queryParams)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	// Verify requests were made
	h.assertRequest("GET", usageEmbeddingsEndpoint, 2)
}

func TestGetModerationsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResults := []ModerationsUsageResult{
		{
			Object:           "usage_result",
			InputTokens:      25,
			NumModelRequests: 1,
			ProjectID:        "proj_789",
			Model:            "text-moderation-latest",
		},
	}

	mockBucket := ModerationsUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672531200,
		EndTime:   1672617600,
		Results:   mockResults,
	}

	response := ModerationsUsageResponse{
		Object:   "list",
		Data:     []ModerationsUsageBucket{mockBucket},
		HasMore:  false,
		NextPage: "",
	}

	// Register mock response
	h.mockResponse("GET", usageModerationsEndpoint, 200, response)

	// Test GetModerationsUsage
	usage, err := h.client.GetModerationsUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage bucket, got %d", len(usage.Data))
		return
	}
	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(usage.Data[0].Results))
		return
	}
	result := usage.Data[0].Results[0]
	if result.InputTokens != 25 {
		t.Errorf("Expected InputTokens 25, got %d", result.InputTokens)
	}
	if result.Model != "text-moderation-latest" {
		t.Errorf("Expected Model text-moderation-latest, got %s", result.Model)
	}

	// Verify the request was made
	h.assertRequest("GET", usageModerationsEndpoint, 1)
}

func TestGetImagesUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResults := []ImagesUsageResult{
		{
			Object:           "usage_result",
			Images:           5,
			NumModelRequests: 1,
			Size:             "1024x1024",
			ProjectID:        "proj_987",
			Model:            "dall-e-3",
		},
	}

	mockBucket := ImagesUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672531200,
		EndTime:   1672617600,
		Results:   mockResults,
	}

	response := ImagesUsageResponse{
		Object:   "list",
		Data:     []ImagesUsageBucket{mockBucket},
		HasMore:  false,
		NextPage: "",
	}

	// Register mock response
	h.mockResponse("GET", usageImagesEndpoint, 200, response)

	// Test GetImagesUsage
	usage, err := h.client.GetImagesUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage bucket, got %d", len(usage.Data))
		return
	}
	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(usage.Data[0].Results))
		return
	}
	result := usage.Data[0].Results[0]
	if result.Images != 5 {
		t.Errorf("Expected Images 5, got %d", result.Images)
	}
	if result.Size != "1024x1024" {
		t.Errorf("Expected Size 1024x1024, got %s", result.Size)
	}

	// Verify the request was made
	h.assertRequest("GET", usageImagesEndpoint, 1)
}

func TestGetAudioSpeechesUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResults := []AudioSpeechesUsageResult{
		{
			Object:           "usage_result",
			Characters:       1000,
			NumModelRequests: 1,
			ProjectID:        "proj_audio",
			Model:            "tts-1",
		},
	}

	mockBucket := AudioSpeechesUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672531200,
		EndTime:   1672617600,
		Results:   mockResults,
	}

	response := AudioSpeechesUsageResponse{
		Object:   "list",
		Data:     []AudioSpeechesUsageBucket{mockBucket},
		HasMore:  false,
		NextPage: "",
	}

	// Register mock response
	h.mockResponse("GET", usageAudioSpeechesEndpoint, 200, response)

	// Test GetAudioSpeechesUsage
	usage, err := h.client.GetAudioSpeechesUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage bucket, got %d", len(usage.Data))
		return
	}
	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(usage.Data[0].Results))
		return
	}
	result := usage.Data[0].Results[0]
	if result.Characters != 1000 {
		t.Errorf("Expected Characters 1000, got %d", result.Characters)
	}
	if result.Model != "tts-1" {
		t.Errorf("Expected Model tts-1, got %s", result.Model)
	}

	// Verify the request was made
	h.assertRequest("GET", usageAudioSpeechesEndpoint, 1)
}

func TestGetAudioTranscriptionsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResults := []AudioTranscriptionsUsageResult{
		{
			Object:           "usage_result",
			Seconds:          120,
			NumModelRequests: 1,
			ProjectID:        "proj_trans",
			Model:            "whisper-1",
		},
	}

	mockBucket := AudioTranscriptionsUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672531200,
		EndTime:   1672617600,
		Results:   mockResults,
	}

	response := AudioTranscriptionsUsageResponse{
		Object:   "list",
		Data:     []AudioTranscriptionsUsageBucket{mockBucket},
		HasMore:  false,
		NextPage: "",
	}

	// Register mock response
	h.mockResponse("GET", usageAudioTranscriptionsEndpoint, 200, response)

	// Test GetAudioTranscriptionsUsage
	usage, err := h.client.GetAudioTranscriptionsUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage bucket, got %d", len(usage.Data))
		return
	}
	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(usage.Data[0].Results))
		return
	}
	result := usage.Data[0].Results[0]
	if result.Seconds != 120 {
		t.Errorf("Expected Seconds 120, got %d", result.Seconds)
	}
	if result.Model != "whisper-1" {
		t.Errorf("Expected Model whisper-1, got %s", result.Model)
	}

	// Verify the request was made
	h.assertRequest("GET", usageAudioTranscriptionsEndpoint, 1)
}

func TestGetVectorStoresUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResults := []VectorStoresUsageResult{
		{
			Object:     "usage_result",
			UsageBytes: 10000000,
			ProjectID:  "proj_vector",
		},
	}

	mockBucket := VectorStoresUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672531200,
		EndTime:   1672617600,
		Results:   mockResults,
	}

	response := VectorStoresUsageResponse{
		Object:   "list",
		Data:     []VectorStoresUsageBucket{mockBucket},
		HasMore:  false,
		NextPage: "",
	}

	// Register mock response
	h.mockResponse("GET", usageVectorStoresEndpoint, 200, response)

	// Test GetVectorStoresUsage
	usage, err := h.client.GetVectorStoresUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage bucket, got %d", len(usage.Data))
		return
	}
	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(usage.Data[0].Results))
		return
	}
	result := usage.Data[0].Results[0]
	if result.UsageBytes != 10000000 {
		t.Errorf("Expected UsageBytes 10000000, got %d", result.UsageBytes)
	}

	// Verify the request was made
	h.assertRequest("GET", usageVectorStoresEndpoint, 1)
}

func TestGetCodeInterpreterUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResults := []CodeInterpreterUsageResult{
		{
			Object:      "usage_result",
			NumSessions: 5,
			ProjectID:   "proj_code",
		},
	}

	mockBucket := CodeInterpreterUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672531200,
		EndTime:   1672617600,
		Results:   mockResults,
	}

	response := CodeInterpreterUsageResponse{
		Object:   "list",
		Data:     []CodeInterpreterUsageBucket{mockBucket},
		HasMore:  false,
		NextPage: "",
	}

	// Register mock response
	h.mockResponse("GET", usageCodeInterpreterEndpoint, 200, response)

	// Test GetCodeInterpreterUsage
	usage, err := h.client.GetCodeInterpreterUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage bucket, got %d", len(usage.Data))
		return
	}
	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(usage.Data[0].Results))
		return
	}
	result := usage.Data[0].Results[0]
	if result.NumSessions != 5 {
		t.Errorf("Expected NumSessions 5, got %d", result.NumSessions)
	}

	// Verify the request was made
	h.assertRequest("GET", usageCodeInterpreterEndpoint, 1)
}

func TestGetCostsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResults := []CostsUsageResult{
		{
			Object: "usage_result",
			Amount: CostAmount{
				Value:    150.25,
				Currency: "USD",
			},
			ProjectID: "proj_costs",
		},
	}

	mockBucket := CostsUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672531200, // 2023-01-01 00:00:00 UTC
		EndTime:   1672617600, // 2023-01-02 00:00:00 UTC
		Results:   mockResults,
	}

	response := CostsUsageResponse{
		Object:   "list",
		Data:     []CostsUsageBucket{mockBucket},
		HasMore:  false,
		NextPage: "",
	}

	// Register mock response
	h.mockResponse("GET", usageCostsEndpoint, 200, response)

	// Test GetCostsUsage
	usage, err := h.client.GetCostsUsage(nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(usage.Data) != 1 {
		t.Errorf("Expected 1 usage bucket, got %d", len(usage.Data))
		return
	}

	if len(usage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 result in bucket, got %d", len(usage.Data[0].Results))
		return
	}

	result := usage.Data[0].Results[0]
	if result.Amount.Value != 150.25 {
		t.Errorf("Expected Amount 150.25, got %f", result.Amount.Value)
	}
	if result.Amount.Currency != "USD" {
		t.Errorf("Expected Currency USD, got %s", result.Amount.Currency)
	}
	if result.ProjectID != "proj_costs" {
		t.Errorf("Expected ProjectID proj_costs, got %s", result.ProjectID)
	}

	// Verify request was made
	h.assertRequest("GET", usageCostsEndpoint, 1)
}

func TestUsageWithPagination(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for first page
	mockResults1 := []CompletionsUsageResult{
		{
			Object:           "usage_result",
			InputTokens:      10,
			OutputTokens:     20,
			NumModelRequests: 1,
			ProjectID:        "proj_123",
			Model:            "gpt-4",
		},
		{
			Object:           "usage_result",
			InputTokens:      15,
			OutputTokens:     25,
			NumModelRequests: 1,
			ProjectID:        "proj_123",
			Model:            "gpt-4",
		},
	}

	firstPageBucket := CompletionsUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672531200,
		EndTime:   1672617600,
		Results:   mockResults1,
	}

	firstPageResponse := CompletionsUsageResponse{
		Object:   "list",
		Data:     []CompletionsUsageBucket{firstPageBucket},
		HasMore:  true,
		NextPage: "next_page_token",
	}

	// Second page data
	mockResults2 := []CompletionsUsageResult{
		{
			Object:           "usage_result",
			InputTokens:      20,
			OutputTokens:     30,
			NumModelRequests: 1,
			ProjectID:        "proj_123",
			Model:            "gpt-4",
		},
	}

	secondPageBucket := CompletionsUsageBucket{
		Object:    "usage_bucket",
		StartTime: 1672617600,
		EndTime:   1672704000,
		Results:   mockResults2,
	}

	secondPageResponse := CompletionsUsageResponse{
		Object:   "list",
		Data:     []CompletionsUsageBucket{secondPageBucket},
		HasMore:  false,
		NextPage: "",
	}

	// Create a map to store responses based on query parameters
	httpmock.Reset()
	httpmock.RegisterResponder("GET", testBaseURL+usageCompletionsEndpoint,
		func(req *http.Request) (*http.Response, error) {
			query := req.URL.Query()
			page := query.Get("page")

			if page == "next_page_token" {
				return httpmock.NewJsonResponse(200, secondPageResponse)
			}
			return httpmock.NewJsonResponse(200, firstPageResponse)
		})

	// Test first page
	firstPage, err := h.client.GetCompletionsUsage(map[string]string{"limit": "2"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(firstPage.Data) != 1 || len(firstPage.Data[0].Results) != 2 {
		t.Errorf("Expected 1 bucket with 2 results in first page")
		return
	}
	if !firstPage.HasMore {
		t.Errorf("Expected HasMore to be true for first page")
		return
	}

	// Test second page using page parameter
	secondPage, err := h.client.GetCompletionsUsage(map[string]string{
		"limit": "2",
		"page":  firstPage.NextPage,
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	if len(secondPage.Data) != 1 || len(secondPage.Data[0].Results) != 1 {
		t.Errorf("Expected 1 bucket with 1 result in second page")
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
	h := newTestHelper(t)
	defer h.cleanup()

	// Register error response
	h.mockErrorResponse("GET", usageCompletionsEndpoint, 401, "Unauthorized")

	// Test error handling
	_, err := h.client.GetCompletionsUsage(nil)
	if err == nil {
		t.Errorf("Expected error, got nil")
		return
	}

	// Verify request was made
	h.assertRequest("GET", usageCompletionsEndpoint, 1)
}
