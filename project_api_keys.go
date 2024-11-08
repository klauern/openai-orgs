package openaiorgs

import "fmt"

const ProjectApiKeysListEndpoint = "/organization/projects/%s/api_keys"

type ProjectApiKey struct {
	Object        string      `json:"object"`
	RedactedValue string      `json:"redacted_value"`
	Name          string      `json:"name"`
	CreatedAt     UnixSeconds `json:"created_at"`
	ID            string      `json:"id"`
	Owner         Owner       `json:"owner"`
}

func (c *Client) ListProjectApiKeys(projectID string, limit int, after string) (*ListResponse[ProjectApiKey], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}

	return Get[ProjectApiKey](c.client, fmt.Sprintf(ProjectApiKeysListEndpoint, projectID), queryParams)
}

func (c *Client) RetrieveProjectApiKey(projectID string, apiKeyID string) (*ProjectApiKey, error) {
	return GetSingle[ProjectApiKey](c.client, fmt.Sprintf(ProjectApiKeysListEndpoint+"/%s", projectID, apiKeyID))
}

func (c *Client) DeleteProjectApiKey(projectID string, apiKeyID string) error {
	return Delete[ProjectApiKey](c.client, fmt.Sprintf(ProjectApiKeysListEndpoint+"/%s", projectID, apiKeyID))
}
