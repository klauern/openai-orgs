package openaiorgs

import (
	"encoding/json"
	"fmt"
	"maps"
	"time"
)

const (
	usageCompletionsEndpoint         = "/organization/usage/completions"
	usageEmbeddingsEndpoint          = "/organization/usage/embeddings"
	usageModerationsEndpoint         = "/organization/usage/moderations"
	usageImagesEndpoint              = "/organization/usage/images"
	usageAudioSpeechesEndpoint       = "/organization/usage/audio_speeches"
	usageAudioTranscriptionsEndpoint = "/organization/usage/audio_transcriptions"
	usageVectorStoresEndpoint        = "/organization/usage/vector_stores"
	usageCodeInterpreterEndpoint     = "/organization/usage/code_interpreter_sessions"
	usageCostsEndpoint               = "/organization/costs"
)

// UsageType represents the type of usage (completions, embeddings, etc.)
type UsageType string

const (
	UsageTypeCompletions         UsageType = "completions"
	UsageTypeEmbeddings          UsageType = "embeddings"
	UsageTypeModerations         UsageType = "moderations"
	UsageTypeImages              UsageType = "images"
	UsageTypeAudioSpeeches       UsageType = "audio_speeches"
	UsageTypeAudioTranscriptions UsageType = "audio_transcriptions"
	UsageTypeVectorStores        UsageType = "vector_stores"
	UsageTypeCodeInterpreter     UsageType = "code_interpreter"
)

// UsageResponse represents the response from the usage endpoints
type UsageResponse struct {
	Object  string        `json:"object"`
	Data    []UsageRecord `json:"data"`
	FirstID string        `json:"first_id"`
	LastID  string        `json:"last_id"`
	HasMore bool          `json:"has_more"`
}

// UsageRecord represents a single usage record
type UsageRecord struct {
	ID           string      `json:"id"`
	Object       string      `json:"object"`
	Timestamp    time.Time   `json:"timestamp"`
	Type         UsageType   `json:"type"`
	UsageDetails interface{} `json:"usage_details"`
	Cost         float64     `json:"cost"`
	ProjectID    string      `json:"project_id"`
	UserID       string      `json:"user_id,omitempty"`
}

// BaseUsageResponse represents the common structure for new usage endpoints
type BaseUsageResponse struct {
	Object   string        `json:"object"`
	Data     []UsageBucket `json:"data"`
	HasMore  bool          `json:"has_more"`
	NextPage string        `json:"next_page"`
}

// UsageBucket represents a time-based bucket of usage data
type UsageBucket struct {
	Object    string          `json:"object"`
	StartTime int64           `json:"start_time"`
	EndTime   int64           `json:"end_time"`
	Results   json.RawMessage `json:"results"`
}

// CompletionsUsageResponse represents the response from the completions usage endpoint
type CompletionsUsageResponse struct {
	Object   string                   `json:"object"`
	Data     []CompletionsUsageBucket `json:"data"`
	HasMore  bool                     `json:"has_more"`
	NextPage string                   `json:"next_page"`
}

// CompletionsUsageBucket represents a time-based bucket of completions usage data
type CompletionsUsageBucket struct {
	Object    string                   `json:"object"`
	StartTime int64                    `json:"start_time"`
	EndTime   int64                    `json:"end_time"`
	Results   []CompletionsUsageResult `json:"results"`
}

// CompletionsUsageResult represents a single completions usage record
type CompletionsUsageResult struct {
	Object            string      `json:"object"`
	InputTokens       int         `json:"input_tokens"`
	OutputTokens      int         `json:"output_tokens"`
	InputCachedTokens int         `json:"input_cached_tokens"`
	InputAudioTokens  int         `json:"input_audio_tokens"`
	OutputAudioTokens int         `json:"output_audio_tokens"`
	NumModelRequests  int         `json:"num_model_requests"`
	ProjectID         string      `json:"project_id"`
	UserID            string      `json:"user_id"`
	APIKeyID          string      `json:"api_key_id"`
	Model             string      `json:"model"`
	Batch             interface{} `json:"batch"`
}

// EmbeddingsUsageResponse represents the response from the embeddings usage endpoint
type EmbeddingsUsageResponse struct {
	Object   string                  `json:"object"`
	Data     []EmbeddingsUsageBucket `json:"data"`
	HasMore  bool                    `json:"has_more"`
	NextPage string                  `json:"next_page"`
}

// EmbeddingsUsageBucket represents a time-based bucket of embeddings usage data
type EmbeddingsUsageBucket struct {
	Object    string                  `json:"object"`
	StartTime int64                   `json:"start_time"`
	EndTime   int64                   `json:"end_time"`
	Results   []EmbeddingsUsageResult `json:"results"`
}

// EmbeddingsUsageResult represents a single embeddings usage record
type EmbeddingsUsageResult struct {
	Object           string `json:"object"`
	InputTokens      int    `json:"input_tokens"`
	NumModelRequests int    `json:"num_model_requests"`
	ProjectID        string `json:"project_id"`
	UserID           string `json:"user_id"`
	APIKeyID         string `json:"api_key_id"`
	Model            string `json:"model"`
}

// ModerationsUsageResponse represents the response from the moderations usage endpoint
type ModerationsUsageResponse struct {
	Object   string                   `json:"object"`
	Data     []ModerationsUsageBucket `json:"data"`
	HasMore  bool                     `json:"has_more"`
	NextPage string                   `json:"next_page"`
}

// ModerationsUsageBucket represents a time-based bucket of moderations usage data
type ModerationsUsageBucket struct {
	Object    string                   `json:"object"`
	StartTime int64                    `json:"start_time"`
	EndTime   int64                    `json:"end_time"`
	Results   []ModerationsUsageResult `json:"results"`
}

// ModerationsUsageResult represents a single moderations usage record
type ModerationsUsageResult struct {
	Object           string `json:"object"`
	InputTokens      int    `json:"input_tokens"`
	NumModelRequests int    `json:"num_model_requests"`
	ProjectID        string `json:"project_id"`
	UserID           string `json:"user_id"`
	APIKeyID         string `json:"api_key_id"`
	Model            string `json:"model"`
}

// ImagesUsageResponse represents the response from the images usage endpoint
type ImagesUsageResponse struct {
	Object   string              `json:"object"`
	Data     []ImagesUsageBucket `json:"data"`
	HasMore  bool                `json:"has_more"`
	NextPage string              `json:"next_page"`
}

// ImagesUsageBucket represents a time-based bucket of images usage data
type ImagesUsageBucket struct {
	Object    string              `json:"object"`
	StartTime int64               `json:"start_time"`
	EndTime   int64               `json:"end_time"`
	Results   []ImagesUsageResult `json:"results"`
}

// ImagesUsageResult represents a single images usage record
type ImagesUsageResult struct {
	Object           string `json:"object"`
	Images           int    `json:"images"`
	NumModelRequests int    `json:"num_model_requests"`
	Size             string `json:"size"`
	Source           string `json:"source"`
	ProjectID        string `json:"project_id"`
	UserID           string `json:"user_id"`
	APIKeyID         string `json:"api_key_id"`
	Model            string `json:"model"`
}

// AudioSpeechesUsageResponse represents the response from the audio speeches usage endpoint
type AudioSpeechesUsageResponse struct {
	Object   string                     `json:"object"`
	Data     []AudioSpeechesUsageBucket `json:"data"`
	HasMore  bool                       `json:"has_more"`
	NextPage string                     `json:"next_page"`
}

// AudioSpeechesUsageBucket represents a time-based bucket of audio speeches usage data
type AudioSpeechesUsageBucket struct {
	Object    string                     `json:"object"`
	StartTime int64                      `json:"start_time"`
	EndTime   int64                      `json:"end_time"`
	Results   []AudioSpeechesUsageResult `json:"results"`
}

// AudioSpeechesUsageResult represents a single audio speeches usage record
type AudioSpeechesUsageResult struct {
	Object           string `json:"object"`
	Characters       int    `json:"characters"`
	NumModelRequests int    `json:"num_model_requests"`
	ProjectID        string `json:"project_id"`
	UserID           string `json:"user_id"`
	APIKeyID         string `json:"api_key_id"`
	Model            string `json:"model"`
}

// AudioTranscriptionsUsageResponse represents the response from the audio transcriptions usage endpoint
type AudioTranscriptionsUsageResponse struct {
	Object   string                           `json:"object"`
	Data     []AudioTranscriptionsUsageBucket `json:"data"`
	HasMore  bool                             `json:"has_more"`
	NextPage string                           `json:"next_page"`
}

// AudioTranscriptionsUsageBucket represents a time-based bucket of audio transcriptions usage data
type AudioTranscriptionsUsageBucket struct {
	Object    string                           `json:"object"`
	StartTime int64                            `json:"start_time"`
	EndTime   int64                            `json:"end_time"`
	Results   []AudioTranscriptionsUsageResult `json:"results"`
}

// AudioTranscriptionsUsageResult represents a single audio transcriptions usage record
type AudioTranscriptionsUsageResult struct {
	Object           string `json:"object"`
	Seconds          int    `json:"seconds"`
	NumModelRequests int    `json:"num_model_requests"`
	ProjectID        string `json:"project_id"`
	UserID           string `json:"user_id"`
	APIKeyID         string `json:"api_key_id"`
	Model            string `json:"model"`
}

// VectorStoresUsageResponse represents the response from the vector stores usage endpoint
type VectorStoresUsageResponse struct {
	Object   string                    `json:"object"`
	Data     []VectorStoresUsageBucket `json:"data"`
	HasMore  bool                      `json:"has_more"`
	NextPage string                    `json:"next_page"`
}

// VectorStoresUsageBucket represents a time-based bucket of vector stores usage data
type VectorStoresUsageBucket struct {
	Object    string                    `json:"object"`
	StartTime int64                     `json:"start_time"`
	EndTime   int64                     `json:"end_time"`
	Results   []VectorStoresUsageResult `json:"results"`
}

// VectorStoresUsageResult represents a single vector stores usage record
type VectorStoresUsageResult struct {
	Object     string `json:"object"`
	UsageBytes int    `json:"usage_bytes"`
	ProjectID  string `json:"project_id"`
}

// CodeInterpreterUsageResponse represents the response from the code interpreter sessions usage endpoint
type CodeInterpreterUsageResponse struct {
	Object   string                       `json:"object"`
	Data     []CodeInterpreterUsageBucket `json:"data"`
	HasMore  bool                         `json:"has_more"`
	NextPage string                       `json:"next_page"`
}

// CodeInterpreterUsageBucket represents a time-based bucket of code interpreter sessions usage data
type CodeInterpreterUsageBucket struct {
	Object    string                       `json:"object"`
	StartTime int64                        `json:"start_time"`
	EndTime   int64                        `json:"end_time"`
	Results   []CodeInterpreterUsageResult `json:"results"`
}

// CodeInterpreterUsageResult represents a single code interpreter sessions usage record
type CodeInterpreterUsageResult struct {
	Object      string `json:"object"`
	NumSessions int    `json:"num_sessions"`
	ProjectID   string `json:"project_id"`
}

// CostsUsageResponse represents the response from the costs usage endpoint
type CostsUsageResponse struct {
	Object   string             `json:"object"`
	Data     []CostsUsageBucket `json:"data"`
	HasMore  bool               `json:"has_more"`
	NextPage string             `json:"next_page"`
}

// CostsUsageBucket represents a time-based bucket of costs usage data
type CostsUsageBucket struct {
	Object    string             `json:"object"`
	StartTime int64              `json:"start_time"`
	EndTime   int64              `json:"end_time"`
	Results   []CostsUsageResult `json:"results"`
}

// CostsUsageResult represents a single costs usage record
type CostsUsageResult struct {
	Object    string      `json:"object"`
	Amount    CostAmount  `json:"amount"`
	LineItem  interface{} `json:"line_item"`
	ProjectID string      `json:"project_id"`
}

// CostAmount represents the monetary amount for a cost
type CostAmount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// CompletionsUsage represents usage details for completions
type CompletionsUsage struct {
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	TotalTokens      int    `json:"total_tokens"`
	Model            string `json:"model"`
}

// EmbeddingsUsage represents usage details for embeddings
type EmbeddingsUsage struct {
	PromptTokens int    `json:"prompt_tokens"`
	Model        string `json:"model"`
}

// ModerationsUsage represents usage details for moderations
type ModerationsUsage struct {
	PromptTokens int    `json:"prompt_tokens"`
	Model        string `json:"model"`
}

// ImagesUsage represents usage details for image generation
type ImagesUsage struct {
	Images int    `json:"images"`
	Size   string `json:"size"`
	Model  string `json:"model"`
}

// AudioSpeechesUsage represents usage details for audio speeches
type AudioSpeechesUsage struct {
	Characters int    `json:"characters"`
	Model      string `json:"model"`
}

// AudioTranscriptionsUsage represents usage details for audio transcriptions
type AudioTranscriptionsUsage struct {
	Seconds int    `json:"seconds"`
	Model   string `json:"model"`
}

// VectorStoresUsage represents usage details for vector stores
type VectorStoresUsage struct {
	Vectors int    `json:"vectors"`
	Size    int    `json:"size"`
	Model   string `json:"model"`
}

// CodeInterpreterUsage represents usage details for code interpreter sessions
type CodeInterpreterUsage struct {
	SessionDuration int    `json:"session_duration"`
	Model           string `json:"model"`
}

// CostsUsage represents usage details for costs
type CostsUsage struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Period   string  `json:"period"`
}

// String returns a human-readable string representation of the UsageRecord
func (ur *UsageRecord) String() string {
	userInfo := ""
	if ur.UserID != "" {
		userInfo = fmt.Sprintf(", UserID: %s", ur.UserID)
	}
	return fmt.Sprintf("UsageRecord{ID: %s, Type: %s, Cost: %.2f, ProjectID: %s%s, Time: %s}",
		ur.ID, ur.Type, ur.Cost, ur.ProjectID, userInfo, ur.Timestamp.Format(time.RFC3339))
}

// GetCompletionsUsage retrieves completions usage data
//
// Parameters:
//   - queryParams: Query parameters including:
//   - start_time: Required - Start time as Unix timestamp (seconds)
//   - end_time: Optional - End time as Unix timestamp (seconds)
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
//
// Returns:
//   - *CompletionsUsageResponse: The completions usage data response
//   - error: Any error that occurred
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

// GetEmbeddingsUsage retrieves embeddings usage data
//
// Parameters:
//   - queryParams: Query parameters including:
//   - start_time: Required - Start time as Unix timestamp (seconds)
//   - end_time: Optional - End time as Unix timestamp (seconds)
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
//
// Returns:
//   - *EmbeddingsUsageResponse: The embeddings usage data response
//   - error: Any error that occurred
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

// GetModerationsUsage retrieves moderations usage data
//
// Parameters:
//   - queryParams: Query parameters including:
//   - start_time: Required - Start time as Unix timestamp (seconds)
//   - end_time: Optional - End time as Unix timestamp (seconds)
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
//
// Returns:
//   - *ModerationsUsageResponse: The moderations usage data response
//   - error: Any error that occurred
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

// GetImagesUsage retrieves images usage data
//
// Parameters:
//   - queryParams: Query parameters including:
//   - start_time: Required - Start time as Unix timestamp (seconds)
//   - end_time: Optional - End time as Unix timestamp (seconds)
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
//
// Returns:
//   - *ImagesUsageResponse: The images usage data response
//   - error: Any error that occurred
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

// GetAudioSpeechesUsage retrieves audio speeches usage data
//
// Parameters:
//   - queryParams: Query parameters including:
//   - start_time: Required - Start time as Unix timestamp (seconds)
//   - end_time: Optional - End time as Unix timestamp (seconds)
//   - limit: Maximum number of records to return
//   - page: Pagination cursor
//   - project_id: Filter by project ID
//
// Returns:
//   - *AudioSpeechesUsageResponse: The audio speeches usage data response
//   - error: Any error that occurred
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

// GetAudioTranscriptionsUsage retrieves audio transcriptions usage data
//
// Parameters:
//   - queryParams: Query parameters including:
//   - start_time: Required - Start time as Unix timestamp (seconds)
//   - end_time: Optional - End time as Unix timestamp (seconds)
//   - limit: Maximum number of records to return
//   - page: Pagination cursor
//   - project_id: Filter by project ID
//
// Returns:
//   - *AudioTranscriptionsUsageResponse: The audio transcriptions usage data response
//   - error: Any error that occurred
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

// GetVectorStoresUsage retrieves vector stores usage data
//
// Parameters:
//   - queryParams: Query parameters including:
//   - start_time: Required - Start time as Unix timestamp (seconds)
//   - end_time: Optional - End time as Unix timestamp (seconds)
//   - limit: Maximum number of records to return
//   - page: Pagination cursor
//   - project_id: Filter by project ID
//
// Returns:
//   - *VectorStoresUsageResponse: The vector stores usage data response
//   - error: Any error that occurred
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

// GetCodeInterpreterUsage retrieves code interpreter usage data
//
// Parameters:
//   - queryParams: Query parameters including:
//   - start_time: Required - Start time as Unix timestamp (seconds)
//   - end_time: Optional - End time as Unix timestamp (seconds)
//   - limit: Maximum number of records to return
//   - page: Pagination cursor
//   - project_id: Filter by project ID
//
// Returns:
//   - *CodeInterpreterUsageResponse: The code interpreter usage data response
//   - error: Any error that occurred
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

// GetCostsUsage retrieves costs usage data
//
// Parameters:
//   - queryParams: Query parameters including:
//   - start_time: Required - Start time as Unix timestamp (seconds)
//   - end_time: Optional - End time as Unix timestamp (seconds)
//   - limit: Maximum number of records to return
//   - page: Pagination cursor
//   - project_id: Filter by project ID
//
// Returns:
//   - *CostsUsageResponse: The costs usage data response
//   - error: Any error that occurred
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
