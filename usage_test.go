package openaiorgs

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestGetCompletionsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockUsageRecord := UsageRecord{
		ID:        "usage_123",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      UsageTypeCompletions,
		UsageDetails: map[string]interface{}{
			"prompt_tokens":     10,
			"completion_tokens": 20,
			"total_tokens":      30,
			"model":             "gpt-4",
		},
		Cost:      0.05,
		ProjectID: "proj_123",
	}

	// Create response for both calls
	response := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockUsageRecord},
		FirstID: "usage_123",
		LastID:  "usage_123",
		HasMore: false,
	}

	// Register mock responses for both calls before making any requests
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
	if usage.Data[0].ID != "usage_123" {
		t.Errorf("Expected ID usage_123, got %s", usage.Data[0].ID)
	}
	if usage.Data[0].Type != UsageTypeCompletions {
		t.Errorf("Expected Type %s, got %s", UsageTypeCompletions, usage.Data[0].Type)
	}
	if usage.Data[0].Cost != 0.05 {
		t.Errorf("Expected Cost 0.05, got %f", usage.Data[0].Cost)
	}
	if usage.Data[0].ProjectID != "proj_123" {
		t.Errorf("Expected ProjectID proj_123, got %s", usage.Data[0].ProjectID)
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
	h.assertRequest("GET", usageCompletionsEndpoint, 2) // Total call count is now 2
}

func TestGetEmbeddingsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockUsageRecord := UsageRecord{
		ID:        "usage_456",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      UsageTypeEmbeddings,
		UsageDetails: map[string]interface{}{
			"prompt_tokens": 50,
			"model":         "text-embedding-ada-002",
		},
		Cost:      0.02,
		ProjectID: "proj_456",
	}

	// Create response for both calls
	response := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockUsageRecord},
		FirstID: "usage_456",
		LastID:  "usage_456",
		HasMore: false,
	}

	// Register mock responses for both calls before making any requests
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
	if usage.Data[0].ID != "usage_456" {
		t.Errorf("Expected ID usage_456, got %s", usage.Data[0].ID)
	}
	if usage.Data[0].Type != UsageTypeEmbeddings {
		t.Errorf("Expected Type %s, got %s", UsageTypeEmbeddings, usage.Data[0].Type)
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

	// Mock response data
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockUsageRecord := UsageRecord{
		ID:        "usage_789",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      UsageTypeModerations,
		UsageDetails: map[string]interface{}{
			"prompt_tokens": 25,
			"model":         "text-moderation-latest",
		},
		Cost:      0.01,
		ProjectID: "proj_789",
	}

	// Create response
	response := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockUsageRecord},
		FirstID: "usage_789",
		LastID:  "usage_789",
		HasMore: false,
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
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}
	if usage.Data[0].Type != UsageTypeModerations {
		t.Errorf("Expected Type %s, got %s", UsageTypeModerations, usage.Data[0].Type)
	}

	// Verify the request was made
	h.assertRequest("GET", usageModerationsEndpoint, 1)
}

func TestGetImagesUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockUsageRecord := UsageRecord{
		ID:        "usage_images_123",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      UsageTypeImages,
		UsageDetails: map[string]interface{}{
			"images": 5,
			"size":   "1024x1024",
			"model":  "dall-e-3",
		},
		Cost:      0.20,
		ProjectID: "proj_987",
	}

	// Create response
	response := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockUsageRecord},
		FirstID: "usage_images_123",
		LastID:  "usage_images_123",
		HasMore: false,
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
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}
	if usage.Data[0].Type != UsageTypeImages {
		t.Errorf("Expected Type %s, got %s", UsageTypeImages, usage.Data[0].Type)
	}

	// Verify the request was made
	h.assertRequest("GET", usageImagesEndpoint, 1)
}

func TestGetAudioSpeechesUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockUsageRecord := UsageRecord{
		ID:        "usage_speech_123",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      UsageTypeAudioSpeeches,
		UsageDetails: map[string]interface{}{
			"characters": 1000,
			"model":      "tts-1",
		},
		Cost:      0.015,
		ProjectID: "proj_audio",
	}

	// Create response
	response := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockUsageRecord},
		FirstID: "usage_speech_123",
		LastID:  "usage_speech_123",
		HasMore: false,
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
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}
	if usage.Data[0].Type != UsageTypeAudioSpeeches {
		t.Errorf("Expected Type %s, got %s", UsageTypeAudioSpeeches, usage.Data[0].Type)
	}

	// Verify the request was made
	h.assertRequest("GET", usageAudioSpeechesEndpoint, 1)
}

func TestGetAudioTranscriptionsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockUsageRecord := UsageRecord{
		ID:        "usage_transcription_123",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      UsageTypeAudioTranscriptions,
		UsageDetails: map[string]interface{}{
			"seconds": 120,
			"model":   "whisper-1",
		},
		Cost:      0.10,
		ProjectID: "proj_trans",
	}

	// Create response
	response := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockUsageRecord},
		FirstID: "usage_transcription_123",
		LastID:  "usage_transcription_123",
		HasMore: false,
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
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}
	if usage.Data[0].Type != UsageTypeAudioTranscriptions {
		t.Errorf("Expected Type %s, got %s", UsageTypeAudioTranscriptions, usage.Data[0].Type)
	}

	// Verify the request was made
	h.assertRequest("GET", usageAudioTranscriptionsEndpoint, 1)
}

func TestGetVectorStoresUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockUsageRecord := UsageRecord{
		ID:        "usage_vector_123",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      UsageTypeVectorStores,
		UsageDetails: map[string]interface{}{
			"vectors": 5000,
			"size":    10000000,
			"model":   "text-embedding-ada-002",
		},
		Cost:      0.25,
		ProjectID: "proj_vector",
	}

	// Create response
	response := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockUsageRecord},
		FirstID: "usage_vector_123",
		LastID:  "usage_vector_123",
		HasMore: false,
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
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}
	if usage.Data[0].Type != UsageTypeVectorStores {
		t.Errorf("Expected Type %s, got %s", UsageTypeVectorStores, usage.Data[0].Type)
	}

	// Verify the request was made
	h.assertRequest("GET", usageVectorStoresEndpoint, 1)
}

func TestGetCodeInterpreterUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockUsageRecord := UsageRecord{
		ID:        "usage_code_123",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      UsageTypeCodeInterpreter,
		UsageDetails: map[string]interface{}{
			"session_duration": 600,
			"model":            "gpt-4",
		},
		Cost:      0.30,
		ProjectID: "proj_code",
	}

	// Create response
	response := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockUsageRecord},
		FirstID: "usage_code_123",
		LastID:  "usage_code_123",
		HasMore: false,
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
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}
	if usage.Data[0].Type != UsageTypeCodeInterpreter {
		t.Errorf("Expected Type %s, got %s", UsageTypeCodeInterpreter, usage.Data[0].Type)
	}

	// Verify the request was made
	h.assertRequest("GET", usageCodeInterpreterEndpoint, 1)
}

func TestGetCostsUsage(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockUsageRecord := UsageRecord{
		ID:        "usage_costs_123",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      "costs",
		UsageDetails: map[string]interface{}{
			"amount":   150.25,
			"currency": "USD",
			"period":   "2023-01",
		},
		ProjectID: "proj_costs",
	}

	// Create response
	response := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockUsageRecord},
		FirstID: "usage_costs_123",
		LastID:  "usage_costs_123",
		HasMore: false,
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
		t.Errorf("Expected 1 usage record, got %d", len(usage.Data))
		return
	}
	if usage.Data[0].ID != "usage_costs_123" {
		t.Errorf("Expected ID usage_costs_123, got %s", usage.Data[0].ID)
	}

	// Verify the request was made
	h.assertRequest("GET", usageCostsEndpoint, 1)
}

func TestUsageWithPagination(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data for first page
	mockTimestamp := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	mockRecord1 := UsageRecord{
		ID:        "usage_page1_1",
		Object:    "usage",
		Timestamp: mockTimestamp,
		Type:      UsageTypeCompletions,
		Cost:      0.05,
		ProjectID: "proj_123",
	}
	mockRecord2 := UsageRecord{
		ID:        "usage_page1_2",
		Object:    "usage",
		Timestamp: mockTimestamp.Add(1 * time.Hour),
		Type:      UsageTypeCompletions,
		Cost:      0.10,
		ProjectID: "proj_123",
	}

	// First page response (with has_more=true)
	firstPageResponse := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockRecord1, mockRecord2},
		FirstID: "usage_page1_1",
		LastID:  "usage_page1_2",
		HasMore: true,
	}

	// Second page data
	mockRecord3 := UsageRecord{
		ID:        "usage_page2_1",
		Object:    "usage",
		Timestamp: mockTimestamp.Add(2 * time.Hour),
		Type:      UsageTypeCompletions,
		Cost:      0.15,
		ProjectID: "proj_123",
	}

	// Second page response
	secondPageResponse := ListResponse[UsageRecord]{
		Object:  "list",
		Data:    []UsageRecord{mockRecord3},
		FirstID: "usage_page2_1",
		LastID:  "usage_page2_1",
		HasMore: false,
	}

	// Create a map to store responses based on query parameters
	httpmock.Reset()
	httpmock.RegisterResponder("GET", testBaseURL+usageCompletionsEndpoint,
		func(req *http.Request) (*http.Response, error) {
			query := req.URL.Query()
			after := query.Get("after")

			if after == "usage_page1_2" {
				// Return second page if "after" is the last ID of the first page
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
	if len(firstPage.Data) != 2 {
		t.Errorf("Expected 2 usage records in first page, got %d", len(firstPage.Data))
		return
	}
	if !firstPage.HasMore {
		t.Errorf("Expected HasMore to be true for first page")
		return
	}

	// Test second page using after parameter with the LastID from first page
	secondPage, err := h.client.GetCompletionsUsage(map[string]string{
		"limit": "2",
		"after": firstPage.LastID,
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}

	// Validate second page response
	if len(secondPage.Data) != 1 {
		t.Errorf("Expected 1 usage record in second page, got %d", len(secondPage.Data))
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
