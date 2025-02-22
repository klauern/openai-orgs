package openaiorgs

import "fmt"

const AdminAPIKeysEndpoint = "/organization/admin_api_keys"

// AdminAPIKey represents an organization-level API key
type AdminAPIKey struct {
	Object        string      `json:"object"`
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	RedactedValue string      `json:"redacted_value"`
	CreatedAt     UnixSeconds `json:"created_at"`
	LastUsedAt    UnixSeconds `json:"last_used_at"`
	Scopes        []string    `json:"scopes"`
}

// ListAdminAPIKeys retrieves a list of organization API keys
func (c *Client) ListAdminAPIKeys(limit int, after string) (*ListResponse[AdminAPIKey], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}

	return Get[AdminAPIKey](c.client, AdminAPIKeysEndpoint, queryParams)
}

// CreateAdminAPIKey creates a new organization API key
func (c *Client) CreateAdminAPIKey(name string, scopes []string) (*AdminAPIKey, error) {
	body := map[string]interface{}{
		"name":   name,
		"scopes": scopes,
	}
	return Post[AdminAPIKey](c.client, AdminAPIKeysEndpoint, body)
}

// RetrieveAdminAPIKey retrieves details of a specific organization API key
func (c *Client) RetrieveAdminAPIKey(apiKeyID string) (*AdminAPIKey, error) {
	return GetSingle[AdminAPIKey](c.client, fmt.Sprintf("%s/%s", AdminAPIKeysEndpoint, apiKeyID))
}

// DeleteAdminAPIKey deletes an organization API key
func (c *Client) DeleteAdminAPIKey(apiKeyID string) error {
	return Delete[AdminAPIKey](c.client, fmt.Sprintf("%s/%s", AdminAPIKeysEndpoint, apiKeyID))
}
