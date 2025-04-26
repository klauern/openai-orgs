// Package openaiorgs implements a client for managing OpenAI organization resources.
// It provides methods for managing projects, users, API keys, and other organizational
// aspects of OpenAI accounts.
package openaiorgs

import "fmt"

// ProjectApiKeysListEndpoint specifies the API endpoint path for project API key operations.
const ProjectApiKeysListEndpoint = "/organization/projects/%s/api_keys"

// ProjectApiKey represents an API key for project authentication.
type ProjectApiKey struct {
	// Object is the type identifier of this resource.
	Object string `json:"object"`
	// RedactedValue is the partially hidden API key string.
	RedactedValue string `json:"redacted_value"`
	// Name is the user-assigned identifier.
	Name string `json:"name"`
	// CreatedAt records when this API key was created.
	CreatedAt UnixSeconds `json:"created_at"`
	// ID uniquely identifies this API key.
	ID string `json:"id"`
	// Owner identifies the entity that controls this API key.
	Owner Owner `json:"owner"`
}

// String implements the Stringer interface, returning a concise representation of the API key.
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

// ListProjectApiKeys lists all API keys associated with the specified project.
// Parameters:
//   - projectID: The ID of the project to list API keys for
//   - limit: Maximum number of keys to return (0 for API default)
//   - after: Pagination cursor for fetching next page of results
//
// Returns a ListResponse containing the API keys and pagination information.
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

// RetrieveProjectApiKey gets a specific API key by its ID.
// Parameters:
//   - projectID: The ID of the project the API key belongs to
//   - apiKeyID: The ID of the specific API key to retrieve
//
// Returns the ProjectApiKey if found, or an error if not found or on API failure.
func (c *Client) RetrieveProjectApiKey(projectID string, apiKeyID string) (*ProjectApiKey, error) {
	return GetSingle[ProjectApiKey](c.client, fmt.Sprintf(ProjectApiKeysListEndpoint+"/%s", projectID, apiKeyID))
}

// DeleteProjectApiKey permanently removes an API key.
// Parameters:
//   - projectID: The ID of the project the API key belongs to
//   - apiKeyID: The ID of the API key to delete
//
// Returns an error if the deletion fails or the key doesn't exist.
func (c *Client) DeleteProjectApiKey(projectID string, apiKeyID string) error {
	return Delete[ProjectApiKey](c.client, fmt.Sprintf(ProjectApiKeysListEndpoint+"/%s", projectID, apiKeyID))
}
