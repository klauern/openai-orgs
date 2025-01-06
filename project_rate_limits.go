package openaiorgs

import (
	"fmt"
)

type ProjectRateLimit struct {
	Object                      string `json:"object"`
	ID                          string `json:"id"`
	Model                       string `json:"model"`
	MaxRequestsPer1Minute       int    `json:"max_requests_per_1_minute"`
	MaxTokensPer1Minute         int    `json:"max_tokens_per_1_minute"`
	MaxImagesPer1Minute         int    `json:"max_images_per_1_minute"`
	MaxAudioMegabytesPer1Minute int    `json:"max_audio_megabytes_per_1_minute"`
	MaxRequestsPer1Day          int    `json:"max_requests_per_1_day"`
	Batch1DayMaxInputTokens     int    `json:"batch_1_day_max_input_tokens"`
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
