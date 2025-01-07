package openaiorgs

import (
	"fmt"
)

type ProjectRateLimit struct {
	Object                      string `json:"object"`
	ID                          string `json:"id"`
	Model                       string `json:"model"`
	MaxRequestsPer1Minute       int64  `json:"max_requests_per_1_minute"`
	MaxTokensPer1Minute         int64  `json:"max_tokens_per_1_minute"`
	MaxImagesPer1Minute         int64  `json:"max_images_per_1_minute"`
	MaxAudioMegabytesPer1Minute int64  `json:"max_audio_megabytes_per_1_minute"`
	MaxRequestsPer1Day          int64  `json:"max_requests_per_1_day"`
	Batch1DayMaxInputTokens     int64  `json:"batch_1_day_max_input_tokens"`
}

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

type ProjectRateLimitRequestFields struct {
	MaxRequestsPer1Minute       int64
	MaxTokensPer1Minute         int64
	MaxImagesPer1Minute         int64
	MaxAudioMegabytesPer1Minute int64
	MaxRequestsPer1Day          int64
	Batch1DayMaxInputTokens     int64
}

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
