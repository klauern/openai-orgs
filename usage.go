package openaiorgs

import (
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
	usageCostsEndpoint               = "/organization/usage/costs"
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

// GetCompletionsUsage retrieves completions usage data
//
// Parameters:
//   - queryParams: Optional query parameters including:
//   - start_date: Start date in YYYY-MM-DD format
//   - end_date: End date in YYYY-MM-DD format
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
//
// Returns:
//   - *UsageResponse: The usage data response
//   - error: Any error that occurred
func (c *Client) GetCompletionsUsage(queryParams map[string]string) (*UsageResponse, error) {
	// Create a copy of the query parameters to avoid modifying the original
	params := make(map[string]string)
	if queryParams != nil {
		for k, v := range queryParams {
			params[k] = v
		}
	}

	listResp, err := Get[UsageRecord](c.client, usageCompletionsEndpoint, params)
	if err != nil {
		return nil, err
	}

	return &UsageResponse{
		Object:  listResp.Object,
		Data:    listResp.Data,
		FirstID: listResp.FirstID,
		LastID:  listResp.LastID,
		HasMore: listResp.HasMore,
	}, nil
}

// GetEmbeddingsUsage retrieves embeddings usage data
//
// Parameters:
//   - queryParams: Optional query parameters including:
//   - start_date: Start date in YYYY-MM-DD format
//   - end_date: End date in YYYY-MM-DD format
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
func (c *Client) GetEmbeddingsUsage(queryParams map[string]string) (*UsageResponse, error) {
	listResp, err := Get[UsageRecord](c.client, usageEmbeddingsEndpoint, queryParams)
	if err != nil {
		return nil, err
	}

	return &UsageResponse{
		Object:  listResp.Object,
		Data:    listResp.Data,
		FirstID: listResp.FirstID,
		LastID:  listResp.LastID,
		HasMore: listResp.HasMore,
	}, nil
}

// GetModerationsUsage retrieves moderations usage data
//
// Parameters:
//   - queryParams: Optional query parameters including:
//   - start_date: Start date in YYYY-MM-DD format
//   - end_date: End date in YYYY-MM-DD format
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
func (c *Client) GetModerationsUsage(queryParams map[string]string) (*UsageResponse, error) {
	listResp, err := Get[UsageRecord](c.client, usageModerationsEndpoint, queryParams)
	if err != nil {
		return nil, err
	}

	return &UsageResponse{
		Object:  listResp.Object,
		Data:    listResp.Data,
		FirstID: listResp.FirstID,
		LastID:  listResp.LastID,
		HasMore: listResp.HasMore,
	}, nil
}

// GetImagesUsage retrieves images usage data
//
// Parameters:
//   - queryParams: Optional query parameters including:
//   - start_date: Start date in YYYY-MM-DD format
//   - end_date: End date in YYYY-MM-DD format
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
func (c *Client) GetImagesUsage(queryParams map[string]string) (*UsageResponse, error) {
	listResp, err := Get[UsageRecord](c.client, usageImagesEndpoint, queryParams)
	if err != nil {
		return nil, err
	}

	return &UsageResponse{
		Object:  listResp.Object,
		Data:    listResp.Data,
		FirstID: listResp.FirstID,
		LastID:  listResp.LastID,
		HasMore: listResp.HasMore,
	}, nil
}

// GetAudioSpeechesUsage retrieves audio speeches usage data
//
// Parameters:
//   - queryParams: Optional query parameters including:
//   - start_date: Start date in YYYY-MM-DD format
//   - end_date: End date in YYYY-MM-DD format
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
func (c *Client) GetAudioSpeechesUsage(queryParams map[string]string) (*UsageResponse, error) {
	listResp, err := Get[UsageRecord](c.client, usageAudioSpeechesEndpoint, queryParams)
	if err != nil {
		return nil, err
	}

	return &UsageResponse{
		Object:  listResp.Object,
		Data:    listResp.Data,
		FirstID: listResp.FirstID,
		LastID:  listResp.LastID,
		HasMore: listResp.HasMore,
	}, nil
}

// GetAudioTranscriptionsUsage retrieves audio transcriptions usage data
//
// Parameters:
//   - queryParams: Optional query parameters including:
//   - start_date: Start date in YYYY-MM-DD format
//   - end_date: End date in YYYY-MM-DD format
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
func (c *Client) GetAudioTranscriptionsUsage(queryParams map[string]string) (*UsageResponse, error) {
	listResp, err := Get[UsageRecord](c.client, usageAudioTranscriptionsEndpoint, queryParams)
	if err != nil {
		return nil, err
	}

	return &UsageResponse{
		Object:  listResp.Object,
		Data:    listResp.Data,
		FirstID: listResp.FirstID,
		LastID:  listResp.LastID,
		HasMore: listResp.HasMore,
	}, nil
}

// GetVectorStoresUsage retrieves vector stores usage data
//
// Parameters:
//   - queryParams: Optional query parameters including:
//   - start_date: Start date in YYYY-MM-DD format
//   - end_date: End date in YYYY-MM-DD format
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
func (c *Client) GetVectorStoresUsage(queryParams map[string]string) (*UsageResponse, error) {
	listResp, err := Get[UsageRecord](c.client, usageVectorStoresEndpoint, queryParams)
	if err != nil {
		return nil, err
	}

	return &UsageResponse{
		Object:  listResp.Object,
		Data:    listResp.Data,
		FirstID: listResp.FirstID,
		LastID:  listResp.LastID,
		HasMore: listResp.HasMore,
	}, nil
}

// GetCodeInterpreterUsage retrieves code interpreter usage data
//
// Parameters:
//   - queryParams: Optional query parameters including:
//   - start_date: Start date in YYYY-MM-DD format
//   - end_date: End date in YYYY-MM-DD format
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
func (c *Client) GetCodeInterpreterUsage(queryParams map[string]string) (*UsageResponse, error) {
	listResp, err := Get[UsageRecord](c.client, usageCodeInterpreterEndpoint, queryParams)
	if err != nil {
		return nil, err
	}

	return &UsageResponse{
		Object:  listResp.Object,
		Data:    listResp.Data,
		FirstID: listResp.FirstID,
		LastID:  listResp.LastID,
		HasMore: listResp.HasMore,
	}, nil
}

// GetCostsUsage retrieves costs usage data
//
// Parameters:
//   - queryParams: Optional query parameters including:
//   - start_date: Start date in YYYY-MM-DD format
//   - end_date: End date in YYYY-MM-DD format
//   - limit: Maximum number of records to return
//   - after: Pagination cursor
//   - project_id: Filter by project ID
func (c *Client) GetCostsUsage(queryParams map[string]string) (*UsageResponse, error) {
	listResp, err := Get[UsageRecord](c.client, usageCostsEndpoint, queryParams)
	if err != nil {
		return nil, err
	}

	return &UsageResponse{
		Object:  listResp.Object,
		Data:    listResp.Data,
		FirstID: listResp.FirstID,
		LastID:  listResp.LastID,
		HasMore: listResp.HasMore,
	}, nil
}
