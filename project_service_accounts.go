package openaiorgs

import "fmt"

// ProjectServiceAccountsListEndpoint is the base endpoint template for service account management.
// The %s placeholder must be filled with the project ID for all requests.
const ProjectServiceAccountsListEndpoint = "/organization/projects/%s/service_accounts"

// ProjectServiceAccount represents a service account within a project.
// Service accounts are non-human identities that can be used for automated access to
// project resources. Each service account can have its own API key and role.
type ProjectServiceAccount struct {
	// Object identifies the type of this resource.
	// This will always be "project_service_account" for ProjectServiceAccount objects.
	Object string `json:"object"`

	// ID is the unique identifier for this service account.
	ID string `json:"id"`

	// Name is the human-readable identifier for the service account.
	// This should be descriptive of the account's purpose.
	Name string `json:"name"`

	// Role defines the service account's access level within the project.
	// Valid values are defined by RoleType: "admin", "developer", "viewer".
	Role string `json:"role"`

	// CreatedAt is the timestamp when this service account was created.
	// The timestamp is in Unix epoch seconds.
	CreatedAt UnixSeconds `json:"created_at"`

	// APIKey contains the API key details if one was generated during creation.
	// This field is only populated in the response of CreateProjectServiceAccount.
	// Subsequent retrievals will not include the API key value.
	APIKey *ProjectServiceAccountAPIKey `json:"api_key,omitempty"`
}

// ProjectServiceAccountAPIKey represents an API key associated with a service account.
// The API key value is only returned once during creation and cannot be retrieved later.
type ProjectServiceAccountAPIKey struct {
	// Object identifies the type of this resource.
	// This will always be "project_service_account_api_key" for ProjectServiceAccountAPIKey objects.
	Object string `json:"object"`

	// Value is the actual API key string to be used for authentication.
	// This value is only returned once during creation and cannot be retrieved later.
	Value string `json:"value"`

	// Name is an optional descriptive name for the API key.
	// This can help identify the key's purpose or usage context.
	Name *string `json:"name"`

	// CreatedAt is the timestamp when this API key was created.
	// The timestamp is in Unix epoch seconds.
	CreatedAt UnixSeconds `json:"created_at"`

	// ID is the unique identifier for this API key.
	// This can be used to identify and revoke the key later.
	ID string `json:"id"`
}

// String returns a human-readable string representation of the ProjectServiceAccount.
// This is useful for logging and debugging purposes. The API key value is never included
// in the string representation for security reasons.
func (psa *ProjectServiceAccount) String() string {
	apiKeyInfo := "no key"
	if psa.APIKey != nil {
		apiKeyInfo = fmt.Sprintf("key:%s", psa.APIKey.ID)
	}
	return fmt.Sprintf("ProjectServiceAccount{ID: %s, Name: %s, Role: %s, APIKey: %s}",
		psa.ID, psa.Name, psa.Role, apiKeyInfo)
}

// ListProjectServiceAccounts retrieves a paginated list of service accounts in a project.
// Service accounts are ordered by creation date, with newest accounts first.
//
// Parameters:
//   - projectID: The unique identifier of the project to list service accounts from
//   - limit: Maximum number of accounts to return (0 for default, which is typically 20)
//   - after: Pagination token for fetching next page (empty string for first page)
//
// Returns a ListResponse containing the service accounts and pagination metadata.
// The ListResponse includes the next pagination token if more results are available.
func (c *Client) ListProjectServiceAccounts(projectID string, limit int, after string) (*ListResponse[ProjectServiceAccount], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}

	return Get[ProjectServiceAccount](c.client, fmt.Sprintf(ProjectServiceAccountsListEndpoint, projectID), queryParams)
}

// CreateProjectServiceAccount creates a new service account in a project.
// The service account is created with a default role and an optional API key.
//
// Parameters:
//   - projectID: The unique identifier of the project to create the service account in
//   - name: The human-readable name for the service account
//
// Returns the created ProjectServiceAccount object or an error if creation fails.
// If successful, the response includes the API key value which should be stored securely
// as it cannot be retrieved later.
func (c *Client) CreateProjectServiceAccount(projectID string, name string) (*ProjectServiceAccount, error) {
	body := map[string]string{"name": name}
	return Post[ProjectServiceAccount](c.client, fmt.Sprintf(ProjectServiceAccountsListEndpoint, projectID), body)
}

// RetrieveProjectServiceAccount fetches details about a specific service account.
// The API key value is never included in the response, even if one exists.
//
// Parameters:
//   - projectID: The unique identifier of the project
//   - serviceAccountID: The unique identifier of the service account to retrieve
//
// Returns the ProjectServiceAccount details or an error if retrieval fails.
// Returns an error if the service account does not exist or if the caller lacks permission.
func (c *Client) RetrieveProjectServiceAccount(projectID string, serviceAccountID string) (*ProjectServiceAccount, error) {
	return GetSingle[ProjectServiceAccount](c.client, fmt.Sprintf(ProjectServiceAccountsListEndpoint+"/%s", projectID, serviceAccountID))
}

// DeleteProjectServiceAccount removes a service account from a project.
// This also invalidates any API keys associated with the service account.
//
// Parameters:
//   - projectID: The unique identifier of the project
//   - serviceAccountID: The unique identifier of the service account to delete
//
// Returns an error if the deletion fails or if the caller lacks permission.
// This operation cannot be undone, and any applications using the service account's
// API key will lose access immediately.
func (c *Client) DeleteProjectServiceAccount(projectID string, serviceAccountID string) error {
	return Delete[ProjectServiceAccount](c.client, fmt.Sprintf(ProjectServiceAccountsListEndpoint+"/%s", projectID, serviceAccountID))
}
