package openaiorgs

import (
	"encoding/json"
	"fmt"
	"maps"
	"time"
)

// API endpoints for different types of usage data.
// Each endpoint returns usage statistics for a specific service type.
const (
	// usageCompletionsEndpoint returns usage data for text completion requests
	usageCompletionsEndpoint = "/organization/usage/completions"
	// usageEmbeddingsEndpoint returns usage data for text embedding requests
	usageEmbeddingsEndpoint = "/organization/usage/embeddings"
	// usageModerationsEndpoint returns usage data for content moderation requests
	usageModerationsEndpoint = "/organization/usage/moderations"
	// usageImagesEndpoint returns usage data for image generation requests
	usageImagesEndpoint = "/organization/usage/images"
	// usageAudioSpeechesEndpoint returns usage data for text-to-speech requests
	usageAudioSpeechesEndpoint = "/organization/usage/audio_speeches"
	// usageAudioTranscriptionsEndpoint returns usage data for speech-to-text requests
	usageAudioTranscriptionsEndpoint = "/organization/usage/audio_transcriptions"
	// usageVectorStoresEndpoint returns usage data for vector storage operations
	usageVectorStoresEndpoint = "/organization/usage/vector_stores"
	// usageCodeInterpreterEndpoint returns usage data for code interpreter sessions
	usageCodeInterpreterEndpoint = "/organization/usage/code_interpreter_sessions"
	// usageCostsEndpoint returns billing cost data across all services
	usageCostsEndpoint = "/organization/costs"
)

// UsageType represents the type of usage (completions, embeddings, etc.).
// This is used to categorize usage records and determine how to parse their details.
type UsageType string

// Supported usage types for different OpenAI services.
const (
	// UsageTypeCompletions represents text completion API usage
	UsageTypeCompletions UsageType = "completions"
	// UsageTypeEmbeddings represents text embedding API usage
	UsageTypeEmbeddings UsageType = "embeddings"
	// UsageTypeModerations represents content moderation API usage
	UsageTypeModerations UsageType = "moderations"
	// UsageTypeImages represents image generation API usage
	UsageTypeImages UsageType = "images"
	// UsageTypeAudioSpeeches represents text-to-speech API usage
	UsageTypeAudioSpeeches UsageType = "audio_speeches"
	// UsageTypeAudioTranscriptions represents speech-to-text API usage
	UsageTypeAudioTranscriptions UsageType = "audio_transcriptions"
	// UsageTypeVectorStores represents vector storage operations usage
	UsageTypeVectorStores UsageType = "vector_stores"
	// UsageTypeCodeInterpreter represents code interpreter session usage
	UsageTypeCodeInterpreter UsageType = "code_interpreter"
)

// UsageResponse represents the response from the legacy usage endpoints.
// This format is being replaced by the more detailed type-specific responses.
type UsageResponse struct {
	// Object identifies the type of this resource.
	// This will always be "list" for UsageResponse objects.
	Object string `json:"object"`

	// Data contains the list of usage records.
	Data []UsageRecord `json:"data"`

	// FirstID is the ID of the first record in this page.
	FirstID string `json:"first_id"`

	// LastID is the ID of the last record in this page.
	LastID string `json:"last_id"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`
}

// UsageRecord represents a single usage record in the legacy format.
// This contains basic information about the usage and a type-specific details field.
type UsageRecord struct {
	// ID uniquely identifies this usage record.
	ID string `json:"id"`

	// Object identifies the type of this resource.
	// This will always be "usage_record" for UsageRecord objects.
	Object string `json:"object"`

	// Timestamp indicates when this usage occurred.
	Timestamp time.Time `json:"timestamp"`

	// Type indicates what kind of usage this record represents.
	// This determines how to interpret the UsageDetails field.
	Type UsageType `json:"type"`

	// UsageDetails contains type-specific usage information.
	// The actual structure depends on the Type field.
	UsageDetails any `json:"usage_details"`

	// Cost is the monetary cost associated with this usage.
	Cost float64 `json:"cost"`

	// ProjectID identifies which project generated this usage.
	ProjectID string `json:"project_id"`

	// UserID identifies which user generated this usage, if applicable.
	UserID string `json:"user_id,omitempty"`
}

// BaseUsageResponse represents the common structure for new usage endpoints.
// This is the base type that type-specific responses embed and extend.
type BaseUsageResponse struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Data contains the list of usage buckets.
	Data []UsageBucket `json:"data"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`

	// NextPage is the pagination token for fetching the next page.
	NextPage string `json:"next_page"`
}

// UsageBucket represents a time-based bucket of usage data.
// Each bucket contains usage statistics for a specific time window.
type UsageBucket struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// StartTime is the start of this time bucket (Unix timestamp).
	StartTime int64 `json:"start_time"`

	// EndTime is the end of this time bucket (Unix timestamp).
	EndTime int64 `json:"end_time"`

	// Results contains the raw usage data for this time bucket.
	// The actual structure depends on the usage type.
	Results json.RawMessage `json:"results"`
}

// CompletionsUsageResponse represents the response from the completions usage endpoint.
// This provides detailed statistics about text completion API usage.
type CompletionsUsageResponse struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Data contains the list of completions usage buckets.
	Data []CompletionsUsageBucket `json:"data"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`

	// NextPage is the pagination token for fetching the next page.
	NextPage string `json:"next_page"`
}

// CompletionsUsageBucket represents a time-based bucket of completions usage data.
// Each bucket contains aggregated statistics about completion API usage.
type CompletionsUsageBucket struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// StartTime is the start of this time bucket (Unix timestamp).
	StartTime int64 `json:"start_time"`

	// EndTime is the end of this time bucket (Unix timestamp).
	EndTime int64 `json:"end_time"`

	// Results contains the detailed completions usage data.
	Results []CompletionsUsageResult `json:"results"`
}

// CompletionsUsageResult represents a single completions usage record.
// This provides detailed information about token usage and request counts.
type CompletionsUsageResult struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// InputTokens is the number of tokens in the prompt.
	InputTokens int `json:"input_tokens"`

	// OutputTokens is the number of tokens in the completion.
	OutputTokens int `json:"output_tokens"`

	// InputCachedTokens is the number of tokens served from cache.
	InputCachedTokens int `json:"input_cached_tokens"`

	// InputAudioTokens is the number of tokens from audio input.
	InputAudioTokens int `json:"input_audio_tokens"`

	// OutputAudioTokens is the number of tokens in audio output.
	OutputAudioTokens int `json:"output_audio_tokens"`

	// NumModelRequests is the number of API calls made.
	NumModelRequests int `json:"num_model_requests"`

	// ProjectID identifies which project generated this usage.
	ProjectID string `json:"project_id"`

	// UserID identifies which user generated this usage.
	UserID string `json:"user_id"`

	// APIKeyID identifies which API key was used.
	APIKeyID string `json:"api_key_id"`

	// Model identifies which model was used.
	Model string `json:"model"`

	// Batch contains information about batch processing, if applicable.
	Batch any `json:"batch"`
}

// EmbeddingsUsageResponse represents the response from the embeddings usage endpoint.
// This provides detailed statistics about text embedding API usage.
type EmbeddingsUsageResponse struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Data contains the list of embeddings usage buckets.
	Data []EmbeddingsUsageBucket `json:"data"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`

	// NextPage is the pagination token for fetching the next page.
	NextPage string `json:"next_page"`
}

// EmbeddingsUsageBucket represents a time-based bucket of embeddings usage data.
// Each bucket contains aggregated statistics about embedding API usage.
type EmbeddingsUsageBucket struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// StartTime is the start of this time bucket (Unix timestamp).
	StartTime int64 `json:"start_time"`

	// EndTime is the end of this time bucket (Unix timestamp).
	EndTime int64 `json:"end_time"`

	// Results contains the detailed embeddings usage data.
	Results []EmbeddingsUsageResult `json:"results"`
}

// EmbeddingsUsageResult represents a single embeddings usage record.
// This provides detailed information about token usage and request counts
// for text embedding operations.
type EmbeddingsUsageResult struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// InputTokens is the number of tokens processed for embedding.
	InputTokens int `json:"input_tokens"`

	// NumModelRequests is the number of API calls made.
	NumModelRequests int `json:"num_model_requests"`

	// ProjectID identifies which project generated this usage.
	ProjectID string `json:"project_id"`

	// UserID identifies which user generated this usage.
	UserID string `json:"user_id"`

	// APIKeyID identifies which API key was used.
	APIKeyID string `json:"api_key_id"`

	// Model identifies which embedding model was used.
	Model string `json:"model"`
}

// ModerationsUsageResponse represents the response from the moderations usage endpoint.
// This provides detailed statistics about content moderation API usage.
type ModerationsUsageResponse struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Data contains the list of moderations usage buckets.
	Data []ModerationsUsageBucket `json:"data"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`

	// NextPage is the pagination token for fetching the next page.
	NextPage string `json:"next_page"`
}

// ModerationsUsageBucket represents a time-based bucket of moderations usage data.
// Each bucket contains aggregated statistics about content moderation API usage.
type ModerationsUsageBucket struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// StartTime is the start of this time bucket (Unix timestamp).
	StartTime int64 `json:"start_time"`

	// EndTime is the end of this time bucket (Unix timestamp).
	EndTime int64 `json:"end_time"`

	// Results contains the detailed moderations usage data.
	Results []ModerationsUsageResult `json:"results"`
}

// ModerationsUsageResult represents a single moderations usage record.
// This provides detailed information about token usage and request counts
// for content moderation operations.
type ModerationsUsageResult struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// InputTokens is the number of tokens processed for moderation.
	InputTokens int `json:"input_tokens"`

	// NumModelRequests is the number of API calls made.
	NumModelRequests int `json:"num_model_requests"`

	// ProjectID identifies which project generated this usage.
	ProjectID string `json:"project_id"`

	// UserID identifies which user generated this usage.
	UserID string `json:"user_id"`

	// APIKeyID identifies which API key was used.
	APIKeyID string `json:"api_key_id"`

	// Model identifies which moderation model was used.
	Model string `json:"model"`
}

// ImagesUsageResponse represents the response from the images usage endpoint.
// This provides detailed statistics about image generation API usage.
type ImagesUsageResponse struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Data contains the list of images usage buckets.
	Data []ImagesUsageBucket `json:"data"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`

	// NextPage is the pagination token for fetching the next page.
	NextPage string `json:"next_page"`
}

// ImagesUsageBucket represents a time-based bucket of images usage data.
// Each bucket contains aggregated statistics about image generation API usage.
type ImagesUsageBucket struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// StartTime is the start of this time bucket (Unix timestamp).
	StartTime int64 `json:"start_time"`

	// EndTime is the end of this time bucket (Unix timestamp).
	EndTime int64 `json:"end_time"`

	// Results contains the detailed images usage data.
	Results []ImagesUsageResult `json:"results"`
}

// ImagesUsageResult represents a single images usage record.
// This provides detailed information about image generation counts
// and configuration options used.
type ImagesUsageResult struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Images is the number of images generated.
	Images int `json:"images"`

	// NumModelRequests is the number of API calls made.
	NumModelRequests int `json:"num_model_requests"`

	// Size specifies the resolution of generated images.
	// Common values are "256x256", "512x512", "1024x1024".
	Size string `json:"size"`

	// Source indicates the image generation method.
	// Can be "generation" for new images or "edit" for modifications.
	Source string `json:"source"`

	// ProjectID identifies which project generated this usage.
	ProjectID string `json:"project_id"`

	// UserID identifies which user generated this usage.
	UserID string `json:"user_id"`

	// APIKeyID identifies which API key was used.
	APIKeyID string `json:"api_key_id"`

	// Model identifies which image generation model was used.
	Model string `json:"model"`
}

// AudioSpeechesUsageResponse represents the response from the audio speeches usage endpoint.
// This provides detailed statistics about text-to-speech API usage.
type AudioSpeechesUsageResponse struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Data contains the list of audio speeches usage buckets.
	Data []AudioSpeechesUsageBucket `json:"data"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`

	// NextPage is the pagination token for fetching the next page.
	NextPage string `json:"next_page"`
}

// AudioSpeechesUsageBucket represents a time-based bucket of audio speeches usage data.
// Each bucket contains aggregated statistics about text-to-speech API usage.
type AudioSpeechesUsageBucket struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// StartTime is the start of this time bucket (Unix timestamp).
	StartTime int64 `json:"start_time"`

	// EndTime is the end of this time bucket (Unix timestamp).
	EndTime int64 `json:"end_time"`

	// Results contains the detailed audio speeches usage data.
	Results []AudioSpeechesUsageResult `json:"results"`
}

// AudioSpeechesUsageResult represents a single audio speeches usage record.
// This provides detailed information about text-to-speech operations,
// including character counts and request statistics.
type AudioSpeechesUsageResult struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Characters is the number of text characters processed.
	Characters int `json:"characters"`

	// NumModelRequests is the number of API calls made.
	NumModelRequests int `json:"num_model_requests"`

	// ProjectID identifies which project generated this usage.
	ProjectID string `json:"project_id"`

	// UserID identifies which user generated this usage.
	UserID string `json:"user_id"`

	// APIKeyID identifies which API key was used.
	APIKeyID string `json:"api_key_id"`

	// Model identifies which text-to-speech model was used.
	Model string `json:"model"`
}

// AudioTranscriptionsUsageResponse represents the response from the audio transcriptions usage endpoint.
// This provides detailed statistics about speech-to-text API usage.
type AudioTranscriptionsUsageResponse struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Data contains the list of audio transcriptions usage buckets.
	Data []AudioTranscriptionsUsageBucket `json:"data"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`

	// NextPage is the pagination token for fetching the next page.
	NextPage string `json:"next_page"`
}

// AudioTranscriptionsUsageBucket represents a time-based bucket of audio transcriptions usage data.
// Each bucket contains aggregated statistics about speech-to-text API usage.
type AudioTranscriptionsUsageBucket struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// StartTime is the start of this time bucket (Unix timestamp).
	StartTime int64 `json:"start_time"`

	// EndTime is the end of this time bucket (Unix timestamp).
	EndTime int64 `json:"end_time"`

	// Results contains the detailed audio transcriptions usage data.
	Results []AudioTranscriptionsUsageResult `json:"results"`
}

// AudioTranscriptionsUsageResult represents a single audio transcriptions usage record.
// This provides detailed information about speech-to-text operations,
// including audio duration and request statistics.
type AudioTranscriptionsUsageResult struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Seconds is the duration of audio processed.
	Seconds int `json:"seconds"`

	// NumModelRequests is the number of API calls made.
	NumModelRequests int `json:"num_model_requests"`

	// ProjectID identifies which project generated this usage.
	ProjectID string `json:"project_id"`

	// UserID identifies which user generated this usage.
	UserID string `json:"user_id"`

	// APIKeyID identifies which API key was used.
	APIKeyID string `json:"api_key_id"`

	// Model identifies which speech-to-text model was used.
	Model string `json:"model"`
}

// VectorStoresUsageResponse represents the response from the vector stores usage endpoint.
// This provides detailed statistics about vector storage operations and capacity usage.
type VectorStoresUsageResponse struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Data contains the list of vector stores usage buckets.
	Data []VectorStoresUsageBucket `json:"data"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`

	// NextPage is the pagination token for fetching the next page.
	NextPage string `json:"next_page"`
}

// VectorStoresUsageBucket represents a time-based bucket of vector stores usage data.
// Each bucket contains aggregated statistics about vector storage operations.
type VectorStoresUsageBucket struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// StartTime is the start of this time bucket (Unix timestamp).
	StartTime int64 `json:"start_time"`

	// EndTime is the end of this time bucket (Unix timestamp).
	EndTime int64 `json:"end_time"`

	// Results contains the detailed vector stores usage data.
	Results []VectorStoresUsageResult `json:"results"`
}

// VectorStoresUsageResult represents a single vector stores usage record.
// This provides detailed information about vector storage capacity
// and data volume used.
type VectorStoresUsageResult struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// UsageBytes is the total storage space used in bytes.
	UsageBytes int `json:"usage_bytes"`

	// ProjectID identifies which project generated this usage.
	ProjectID string `json:"project_id"`
}

// CodeInterpreterUsageResponse represents the response from the code interpreter usage endpoint.
// This provides detailed statistics about code interpreter session usage.
type CodeInterpreterUsageResponse struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Data contains the list of code interpreter usage buckets.
	Data []CodeInterpreterUsageBucket `json:"data"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`

	// NextPage is the pagination token for fetching the next page.
	NextPage string `json:"next_page"`
}

// CodeInterpreterUsageBucket represents a time-based bucket of code interpreter usage data.
// Each bucket contains aggregated statistics about code interpreter sessions.
type CodeInterpreterUsageBucket struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// StartTime is the start of this time bucket (Unix timestamp).
	StartTime int64 `json:"start_time"`

	// EndTime is the end of this time bucket (Unix timestamp).
	EndTime int64 `json:"end_time"`

	// Results contains the detailed code interpreter usage data.
	Results []CodeInterpreterUsageResult `json:"results"`
}

// CodeInterpreterUsageResult represents a single code interpreter usage record.
// This provides detailed information about code interpreter session counts
// and resource utilization.
type CodeInterpreterUsageResult struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// NumSessions is the number of code interpreter sessions used.
	NumSessions int `json:"num_sessions"`

	// ProjectID identifies which project generated this usage.
	ProjectID string `json:"project_id"`
}

// CostsUsageResponse represents the response from the costs usage endpoint.
// This provides detailed statistics about billing and costs across all services.
type CostsUsageResponse struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Data contains the list of costs usage buckets.
	Data []CostsUsageBucket `json:"data"`

	// HasMore indicates if there are more records available.
	HasMore bool `json:"has_more"`

	// NextPage is the pagination token for fetching the next page.
	NextPage string `json:"next_page"`
}

// CostsUsageBucket represents a time-based bucket of costs usage data.
// Each bucket contains aggregated billing information for all services.
type CostsUsageBucket struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// StartTime is the start of this time bucket (Unix timestamp).
	StartTime int64 `json:"start_time"`

	// EndTime is the end of this time bucket (Unix timestamp).
	EndTime int64 `json:"end_time"`

	// Results contains the detailed costs usage data.
	Results []CostsUsageResult `json:"results"`
}

// CostsUsageResult represents a single costs usage record.
// This provides detailed information about billing amounts
// and associated line items.
type CostsUsageResult struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`

	// Amount contains the cost value and currency information.
	Amount CostAmount `json:"amount"`

	// LineItem contains additional details about what the cost is for.
	// The structure varies depending on the service type.
	LineItem any `json:"line_item"`

	// ProjectID identifies which project generated this cost.
	ProjectID string `json:"project_id"`
}

// CostAmount represents a monetary amount with its currency.
type CostAmount struct {
	// Value is the numeric amount.
	Value float64 `json:"value"`

	// Currency is the three-letter currency code (e.g., "USD").
	Currency string `json:"currency"`
}

// CompletionsUsage represents the usage details for a completions request.
// This is used in the UsageDetails field of UsageRecord when Type is UsageTypeCompletions.
type CompletionsUsage struct {
	// PromptTokens is the number of tokens in the input prompt.
	PromptTokens int `json:"prompt_tokens"`

	// CompletionTokens is the number of tokens in the generated completion.
	CompletionTokens int `json:"completion_tokens"`

	// TotalTokens is the sum of prompt and completion tokens.
	TotalTokens int `json:"total_tokens"`

	// Model identifies which completion model was used.
	Model string `json:"model"`
}

// EmbeddingsUsage represents the usage details for an embeddings request.
// This is used in the UsageDetails field of UsageRecord when Type is UsageTypeEmbeddings.
type EmbeddingsUsage struct {
	// PromptTokens is the number of tokens processed for embedding.
	PromptTokens int `json:"prompt_tokens"`

	// Model identifies which embedding model was used.
	Model string `json:"model"`
}

// ModerationsUsage represents the usage details for a moderations request.
// This is used in the UsageDetails field of UsageRecord when Type is UsageTypeModerations.
type ModerationsUsage struct {
	// PromptTokens is the number of tokens processed for moderation.
	PromptTokens int `json:"prompt_tokens"`

	// Model identifies which moderation model was used.
	Model string `json:"model"`
}

// ImagesUsage represents the usage details for an image generation request.
// This is used in the UsageDetails field of UsageRecord when Type is UsageTypeImages.
type ImagesUsage struct {
	// Images is the number of images generated.
	Images int `json:"images"`

	// Size specifies the resolution of generated images.
	// Common values are "256x256", "512x512", "1024x1024".
	Size string `json:"size"`

	// Model identifies which image generation model was used.
	Model string `json:"model"`
}

// AudioSpeechesUsage represents the usage details for a text-to-speech request.
// This is used in the UsageDetails field of UsageRecord when Type is UsageTypeAudioSpeeches.
type AudioSpeechesUsage struct {
	// Characters is the number of text characters processed.
	Characters int `json:"characters"`

	// Model identifies which text-to-speech model was used.
	Model string `json:"model"`
}

// AudioTranscriptionsUsage represents the usage details for a speech-to-text request.
// This is used in the UsageDetails field of UsageRecord when Type is UsageTypeAudioTranscriptions.
type AudioTranscriptionsUsage struct {
	// Seconds is the duration of audio processed.
	Seconds int `json:"seconds"`

	// Model identifies which speech-to-text model was used.
	Model string `json:"model"`
}

// VectorStoresUsage represents the usage details for vector storage operations.
// This is used in the UsageDetails field of UsageRecord when Type is UsageTypeVectorStores.
type VectorStoresUsage struct {
	// Vectors is the number of vectors stored.
	Vectors int `json:"vectors"`

	// Size is the total storage space used in bytes.
	Size int `json:"size"`

	// Model identifies which vector storage model was used.
	Model string `json:"model"`
}

// CodeInterpreterUsage represents the usage details for a code interpreter session.
// This is used in the UsageDetails field of UsageRecord when Type is UsageTypeCodeInterpreter.
type CodeInterpreterUsage struct {
	// SessionDuration is the duration of the session in seconds.
	SessionDuration int `json:"session_duration"`

	// Model identifies which code interpreter model was used.
	Model string `json:"model"`
}

// CostsUsage represents the usage details for billing costs.
// This is used in the UsageDetails field of UsageRecord when Type is related to costs.
type CostsUsage struct {
	// Amount is the monetary value of the cost.
	Amount float64 `json:"amount"`

	// Currency is the three-letter currency code (e.g., "USD").
	Currency string `json:"currency"`

	// Period is the billing period this cost applies to.
	Period string `json:"period"`
}

// String returns a human-readable string representation of the UsageRecord.
// This is useful for logging and debugging purposes. The format varies based
// on the Type field and corresponding UsageDetails structure.
func (ur *UsageRecord) String() string {
	userInfo := ""
	if ur.UserID != "" {
		userInfo = fmt.Sprintf(", UserID: %s", ur.UserID)
	}
	return fmt.Sprintf("UsageRecord{ID: %s, Type: %s, Cost: %.2f, ProjectID: %s%s, Time: %s}",
		ur.ID, ur.Type, ur.Cost, ur.ProjectID, userInfo, ur.Timestamp.Format(time.RFC3339))
}

// GetCompletionsUsage retrieves usage statistics for text completion API calls.
// The response includes token counts, request counts, and other metrics aggregated
// by time buckets.
//
// Parameters:
//   - queryParams: Optional query parameters for filtering and pagination.
//     Supported parameters include:
//   - start_time: Start of the time range (Unix timestamp)
//   - end_time: End of the time range (Unix timestamp)
//   - limit: Maximum number of buckets to return
//   - after: Pagination token for the next page
//   - project_id: Filter by specific project
//   - user_id: Filter by specific user
//   - api_key_id: Filter by specific API key
//
// Returns the usage statistics or an error if the request fails.
func (c *Client) GetCompletionsUsage(queryParams map[string]string) (*CompletionsUsageResponse, error) {
	// Create a copy of the query parameters to avoid modifying the original
	params := make(map[string]string)
	if queryParams != nil {
		maps.Copy(params, queryParams)
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		Get(usageCompletionsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected Content-Type \"application/json\", got %q", contentType)
	}

	var completionsResp CompletionsUsageResponse
	err = json.Unmarshal(resp.Body(), &completionsResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &completionsResp, nil
}

// GetEmbeddingsUsage retrieves usage statistics for text embedding API calls.
// The response includes token counts, request counts, and other metrics aggregated
// by time buckets.
//
// Parameters:
//   - queryParams: Optional query parameters for filtering and pagination.
//     Supported parameters include:
//   - start_time: Start of the time range (Unix timestamp)
//   - end_time: End of the time range (Unix timestamp)
//   - limit: Maximum number of buckets to return
//   - after: Pagination token for the next page
//   - project_id: Filter by specific project
//   - user_id: Filter by specific user
//   - api_key_id: Filter by specific API key
//
// Returns the usage statistics or an error if the request fails.
func (c *Client) GetEmbeddingsUsage(queryParams map[string]string) (*EmbeddingsUsageResponse, error) {
	// Create a copy of the query parameters to avoid modifying the original
	params := make(map[string]string)
	if queryParams != nil {
		maps.Copy(params, queryParams)
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		Get(usageEmbeddingsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected Content-Type \"application/json\", got %q", contentType)
	}

	var embeddingsResp EmbeddingsUsageResponse
	err = json.Unmarshal(resp.Body(), &embeddingsResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &embeddingsResp, nil
}

// GetModerationsUsage retrieves usage statistics for content moderation API calls.
// The response includes token counts, request counts, and other metrics aggregated
// by time buckets.
//
// Parameters:
//   - queryParams: Optional query parameters for filtering and pagination.
//     Supported parameters include:
//   - start_time: Start of the time range (Unix timestamp)
//   - end_time: End of the time range (Unix timestamp)
//   - limit: Maximum number of buckets to return
//   - after: Pagination token for the next page
//   - project_id: Filter by specific project
//   - user_id: Filter by specific user
//   - api_key_id: Filter by specific API key
//
// Returns the usage statistics or an error if the request fails.
func (c *Client) GetModerationsUsage(queryParams map[string]string) (*ModerationsUsageResponse, error) {
	// Create a copy of the query parameters to avoid modifying the original
	params := make(map[string]string)
	if queryParams != nil {
		maps.Copy(params, queryParams)
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		Get(usageModerationsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected Content-Type \"application/json\", got %q", contentType)
	}

	var moderationsResp ModerationsUsageResponse
	err = json.Unmarshal(resp.Body(), &moderationsResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &moderationsResp, nil
}

// GetImagesUsage retrieves usage statistics for image generation API calls.
// The response includes image counts, request counts, and other metrics aggregated
// by time buckets.
//
// Parameters:
//   - queryParams: Optional query parameters for filtering and pagination.
//     Supported parameters include:
//   - start_time: Start of the time range (Unix timestamp)
//   - end_time: End of the time range (Unix timestamp)
//   - limit: Maximum number of buckets to return
//   - after: Pagination token for the next page
//   - project_id: Filter by specific project
//   - user_id: Filter by specific user
//   - api_key_id: Filter by specific API key
//
// Returns the usage statistics or an error if the request fails.
func (c *Client) GetImagesUsage(queryParams map[string]string) (*ImagesUsageResponse, error) {
	// Create a copy of the query parameters to avoid modifying the original
	params := make(map[string]string)
	if queryParams != nil {
		maps.Copy(params, queryParams)
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		Get(usageImagesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected Content-Type \"application/json\", got %q", contentType)
	}

	var imagesResp ImagesUsageResponse
	err = json.Unmarshal(resp.Body(), &imagesResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &imagesResp, nil
}

// GetAudioSpeechesUsage retrieves usage statistics for text-to-speech API calls.
// The response includes character counts, request counts, and other metrics aggregated
// by time buckets.
//
// Parameters:
//   - queryParams: Optional query parameters for filtering and pagination.
//     Supported parameters include:
//   - start_time: Start of the time range (Unix timestamp)
//   - end_time: End of the time range (Unix timestamp)
//   - limit: Maximum number of buckets to return
//   - after: Pagination token for the next page
//   - project_id: Filter by specific project
//   - user_id: Filter by specific user
//   - api_key_id: Filter by specific API key
//
// Returns the usage statistics or an error if the request fails.
func (c *Client) GetAudioSpeechesUsage(queryParams map[string]string) (*AudioSpeechesUsageResponse, error) {
	// Create a copy of the query parameters to avoid modifying the original
	params := make(map[string]string)
	if queryParams != nil {
		maps.Copy(params, queryParams)
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		Get(usageAudioSpeechesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected Content-Type \"application/json\", got %q", contentType)
	}

	var audioSpeechesResp AudioSpeechesUsageResponse
	err = json.Unmarshal(resp.Body(), &audioSpeechesResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &audioSpeechesResp, nil
}

// GetAudioTranscriptionsUsage retrieves usage statistics for speech-to-text API calls.
// The response includes audio duration, request counts, and other metrics aggregated
// by time buckets.
//
// Parameters:
//   - queryParams: Optional query parameters for filtering and pagination.
//     Supported parameters include:
//   - start_time: Start of the time range (Unix timestamp)
//   - end_time: End of the time range (Unix timestamp)
//   - limit: Maximum number of buckets to return
//   - after: Pagination token for the next page
//   - project_id: Filter by specific project
//   - user_id: Filter by specific user
//   - api_key_id: Filter by specific API key
//
// Returns the usage statistics or an error if the request fails.
func (c *Client) GetAudioTranscriptionsUsage(queryParams map[string]string) (*AudioTranscriptionsUsageResponse, error) {
	// Create a copy of the query parameters to avoid modifying the original
	params := make(map[string]string)
	if queryParams != nil {
		maps.Copy(params, queryParams)
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		Get(usageAudioTranscriptionsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected Content-Type \"application/json\", got %q", contentType)
	}

	var audioTranscriptionsResp AudioTranscriptionsUsageResponse
	err = json.Unmarshal(resp.Body(), &audioTranscriptionsResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &audioTranscriptionsResp, nil
}

// GetVectorStoresUsage retrieves usage statistics for vector storage operations.
// The response includes storage size, vector counts, and other metrics aggregated
// by time buckets.
//
// Parameters:
//   - queryParams: Optional query parameters for filtering and pagination.
//     Supported parameters include:
//   - start_time: Start of the time range (Unix timestamp)
//   - end_time: End of the time range (Unix timestamp)
//   - limit: Maximum number of buckets to return
//   - after: Pagination token for the next page
//   - project_id: Filter by specific project
//
// Returns the usage statistics or an error if the request fails.
func (c *Client) GetVectorStoresUsage(queryParams map[string]string) (*VectorStoresUsageResponse, error) {
	// Create a copy of the query parameters to avoid modifying the original
	params := make(map[string]string)
	if queryParams != nil {
		maps.Copy(params, queryParams)
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		Get(usageVectorStoresEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected Content-Type \"application/json\", got %q", contentType)
	}

	var vectorStoresResp VectorStoresUsageResponse
	err = json.Unmarshal(resp.Body(), &vectorStoresResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &vectorStoresResp, nil
}

// GetCodeInterpreterUsage retrieves usage statistics for code interpreter sessions.
// The response includes session counts, duration, and other metrics aggregated
// by time buckets.
//
// Parameters:
//   - queryParams: Optional query parameters for filtering and pagination.
//     Supported parameters include:
//   - start_time: Start of the time range (Unix timestamp)
//   - end_time: End of the time range (Unix timestamp)
//   - limit: Maximum number of buckets to return
//   - after: Pagination token for the next page
//   - project_id: Filter by specific project
//
// Returns the usage statistics or an error if the request fails.
func (c *Client) GetCodeInterpreterUsage(queryParams map[string]string) (*CodeInterpreterUsageResponse, error) {
	// Create a copy of the query parameters to avoid modifying the original
	params := make(map[string]string)
	if queryParams != nil {
		maps.Copy(params, queryParams)
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		Get(usageCodeInterpreterEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected Content-Type \"application/json\", got %q", contentType)
	}

	var codeInterpreterResp CodeInterpreterUsageResponse
	err = json.Unmarshal(resp.Body(), &codeInterpreterResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &codeInterpreterResp, nil
}

// GetCostsUsage retrieves billing cost statistics across all services.
// The response includes monetary amounts, line items, and other metrics aggregated
// by time buckets.
//
// Parameters:
//   - queryParams: Optional query parameters for filtering and pagination.
//     Supported parameters include:
//   - start_time: Start of the time range (Unix timestamp)
//   - end_time: End of the time range (Unix timestamp)
//   - limit: Maximum number of buckets to return
//   - after: Pagination token for the next page
//   - project_id: Filter by specific project
//
// Returns the cost statistics or an error if the request fails.
func (c *Client) GetCostsUsage(queryParams map[string]string) (*CostsUsageResponse, error) {
	// Create a copy of the query parameters to avoid modifying the original
	params := make(map[string]string)
	if queryParams != nil {
		maps.Copy(params, queryParams)
	}

	resp, err := c.client.R().
		SetQueryParams(params).
		Get(usageCostsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	contentType := resp.Header().Get("Content-Type")
	if contentType != "application/json" {
		return nil, fmt.Errorf("expected Content-Type \"application/json\", got %q", contentType)
	}

	var costsResp CostsUsageResponse
	err = json.Unmarshal(resp.Body(), &costsResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &costsResp, nil
}
