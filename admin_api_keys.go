package openaiorgs

import (
	"fmt"
	"strings"
)

// AdminAPIKeysEndpoint is the base endpoint for organization API key operations.
const AdminAPIKeysEndpoint = "/organization/admin_api_keys"

// AdminAPIKey represents an organization-level API key with its associated metadata.
// These keys provide administrative access to organization resources and operations.
type AdminAPIKey struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`
	// ID is the unique identifier for this API key.
	ID string `json:"id"`
	// Name is a human-readable identifier for the API key.
	Name string `json:"name"`
	// RedactedValue is a partially hidden version of the API key for display purposes.
	RedactedValue string `json:"redacted_value"`
	// CreatedAt is the timestamp when this API key was created.
	CreatedAt UnixSeconds `json:"created_at"`
	// LastUsedAt is the timestamp when this API key was last used.
	LastUsedAt UnixSeconds `json:"last_used_at"`
	// Scopes define the permissions granted to this API key.
	Scopes []string `json:"scopes"`
}

// ListAdminAPIKeys retrieves a paginated list of organization API keys.
//
// Parameters:
//   - limit: Maximum number of keys to return (0 for default)
//   - after: Pagination token for fetching next page (empty string for first page)
//
// Returns a ListResponse containing the API keys and pagination metadata.
// Returns an error if the API request fails.
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

// CreateAdminAPIKey creates a new organization API key.
//
// Parameters:
//   - name: A human-readable name for the API key
//   - scopes: List of permission scopes to grant to the key
//
// Returns the newly created AdminAPIKey or an error if creation fails.
// Note: The full API key value is only returned once upon creation.
func (c *Client) CreateAdminAPIKey(name string, scopes []string) (*AdminAPIKey, error) {
	body := map[string]any{
		"name":   name,
		"scopes": scopes,
	}
	return Post[AdminAPIKey](c.client, AdminAPIKeysEndpoint, body)
}

// RetrieveAdminAPIKey fetches details of a specific organization API key.
//
// Parameters:
//   - apiKeyID: The unique identifier of the API key to retrieve
//
// Returns the AdminAPIKey details or an error if retrieval fails.
func (c *Client) RetrieveAdminAPIKey(apiKeyID string) (*AdminAPIKey, error) {
	return GetSingle[AdminAPIKey](c.client, fmt.Sprintf("%s/%s", AdminAPIKeysEndpoint, apiKeyID))
}

// DeleteAdminAPIKey permanently removes an organization API key.
//
// Parameters:
//   - apiKeyID: The unique identifier of the API key to delete
//
// Returns an error if deletion fails or nil on success.
func (c *Client) DeleteAdminAPIKey(apiKeyID string) error {
	return Delete[AdminAPIKey](c.client, fmt.Sprintf("%s/%s", AdminAPIKeysEndpoint, apiKeyID))
}

// String returns a human-readable string representation of the AdminAPIKey.
// It includes the key's ID, name, and a comma-separated list of scopes if any exist.
func (ak *AdminAPIKey) String() string {
	scopeInfo := "no scopes"
	if len(ak.Scopes) > 0 {
		scopeInfo = fmt.Sprintf("scopes:%s", strings.Join(ak.Scopes, ","))
	}
	return fmt.Sprintf("AdminAPIKey{ID: %s, Name: %s, %s}",
		ak.ID, ak.Name, scopeInfo)
}
