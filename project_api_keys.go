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

// String returns a human-readable string representation of the ProjectApiKey
func (pak *ProjectApiKey) String() string {
	ownerInfo := "no owner"
	if pak.Owner.Name != "" {
		ownerInfo = fmt.Sprintf("%s(%s)", pak.Owner.Name, pak.Owner.Type)
	}
	name := pak.Name
	if name == "" {
		name = "unnamed"
	}
	return fmt.Sprintf("ProjectApiKey{ID: %s, Name: %s, Owner: %s}",
		pak.ID, name, ownerInfo)
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
