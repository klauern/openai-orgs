package openaiorgs

import (
	"fmt"
)

// ProjectRateLimit represents the rate limiting configuration for a specific model
// within a project. Rate limits control the maximum usage of OpenAI's APIs across
// different time windows and resource types (requests, tokens, images, audio).
type ProjectRateLimit struct {
	// Object identifies the type of this resource.
	// This will always be "project_rate_limit" for ProjectRateLimit objects.
	Object string `json:"object"`

	// ID is the unique identifier for this rate limit configuration.
	ID string `json:"id"`

	// Model specifies which OpenAI model this rate limit applies to.
	// Examples include "gpt-4", "gpt-3.5-turbo", "dall-e-3", etc.
	Model string `json:"model"`

	// MaxRequestsPer1Minute defines the maximum number of API requests
	// allowed within a 1-minute rolling window.
	MaxRequestsPer1Minute int64 `json:"max_requests_per_1_minute"`

	// MaxTokensPer1Minute defines the maximum number of tokens (input + output)
	// that can be processed within a 1-minute rolling window.
	MaxTokensPer1Minute int64 `json:"max_tokens_per_1_minute"`

	// MaxImagesPer1Minute defines the maximum number of images that can be
	// generated within a 1-minute rolling window.
	MaxImagesPer1Minute int64 `json:"max_images_per_1_minute"`

	// MaxAudioMegabytesPer1Minute defines the maximum audio data size in MB
	// that can be processed within a 1-minute rolling window.
	MaxAudioMegabytesPer1Minute int64 `json:"max_audio_megabytes_per_1_minute"`

	// MaxRequestsPer1Day defines the maximum number of API requests
	// allowed within a 24-hour rolling window.
	MaxRequestsPer1Day int64 `json:"max_requests_per_1_day"`

	// Batch1DayMaxInputTokens defines the maximum number of input tokens
	// allowed for batch processing within a 24-hour rolling window.
	Batch1DayMaxInputTokens int64 `json:"batch_1_day_max_input_tokens"`
}

// String returns a human-readable string representation of the ProjectRateLimit.
// This is useful for logging and debugging purposes, showing key rate limit values.
func (prl *ProjectRateLimit) String() string {
	return fmt.Sprintf("ProjectRateLimit{ID: %s, Model: %s, MaxReq/Min: %d, MaxTokens/Min: %d}",
		prl.ID, prl.Model, prl.MaxRequestsPer1Minute, prl.MaxTokensPer1Minute)
}

// ListProjectRateLimits retrieves a paginated list of rate limits for a specific project.
// Each rate limit configuration applies to a different model or API endpoint.
//
// Parameters:
//   - limit: Maximum number of rate limits to return (0 for default, which is typically 20)
//   - after: Pagination token for fetching next page (empty string for first page)
//   - projectId: The unique identifier of the project to list rate limits from
//
// Returns a ListResponse containing the rate limits and pagination metadata.
// The ListResponse includes the next pagination token if more results are available.
func (c *Client) ListProjectRateLimits(limit int, after string, projectId string) (*ListResponse[ProjectRateLimit], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}

	path := fmt.Sprintf("%s/%s/rate_limits", ProjectsListEndpoint, projectId)
	return Get[ProjectRateLimit](c.client, path, queryParams)
}

// ProjectRateLimitRequestFields defines the modifiable fields when updating a rate limit.
// All fields are optional - only non-zero values will be included in the update request.
// This allows for partial updates of rate limit configurations.
type ProjectRateLimitRequestFields struct {
	// MaxRequestsPer1Minute defines the new maximum requests per minute.
	// Set to 0 to leave unchanged.
	MaxRequestsPer1Minute int64

	// MaxTokensPer1Minute defines the new maximum tokens per minute.
	// Set to 0 to leave unchanged.
	MaxTokensPer1Minute int64

	// MaxImagesPer1Minute defines the new maximum images per minute.
	// Set to 0 to leave unchanged.
	MaxImagesPer1Minute int64

	// MaxAudioMegabytesPer1Minute defines the new maximum audio MB per minute.
	// Set to 0 to leave unchanged.
	MaxAudioMegabytesPer1Minute int64

	// MaxRequestsPer1Day defines the new maximum requests per day.
	// Set to 0 to leave unchanged.
	MaxRequestsPer1Day int64

	// Batch1DayMaxInputTokens defines the new maximum batch input tokens per day.
	// Set to 0 to leave unchanged.
	Batch1DayMaxInputTokens int64
}

// ModifyProjectRateLimit updates the rate limit configuration for a specific model
// within a project. Only the non-zero fields in the request will be updated.
//
// Parameters:
//   - projectId: The unique identifier of the project containing the rate limit
//   - rateLimitId: The unique identifier of the rate limit configuration to modify
//   - fields: The new values to set for the rate limit fields
//
// Returns the updated ProjectRateLimit object or an error if modification fails.
// Common errors include invalid rate limit values or insufficient permissions.
func (c *Client) ModifyProjectRateLimit(projectId, rateLimitId string, fields ProjectRateLimitRequestFields) (*ProjectRateLimit, error) {
	body := map[string]int64{}
	if fields.MaxRequestsPer1Minute > 0 {
		body["max_requests_per_1_minute"] = fields.MaxRequestsPer1Minute
	}

	if fields.MaxTokensPer1Minute > 0 {
		body["max_tokens_per_1_minute"] = fields.MaxTokensPer1Minute
	}

	if fields.MaxImagesPer1Minute > 0 {
		body["max_images_per_1_minute"] = fields.MaxImagesPer1Minute
	}

	if fields.MaxAudioMegabytesPer1Minute > 0 {
		body["max_audio_megabytes_per_1_minute"] = fields.MaxAudioMegabytesPer1Minute
	}

	if fields.MaxRequestsPer1Day > 0 {
		body["max_requests_per_1_day"] = fields.MaxRequestsPer1Day
	}

	if fields.Batch1DayMaxInputTokens > 0 {
		body["batch_1_day_max_input_tokens"] = fields.Batch1DayMaxInputTokens
	}

	path := fmt.Sprintf("%s/%s/rate_limits/%s", ProjectsListEndpoint, projectId, rateLimitId)
	return Post[ProjectRateLimit](c.client, path, body)
}
