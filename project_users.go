package openaiorgs

import "fmt"

// ProjectUser represents a user's membership and role within a specific project.
// Each ProjectUser entry defines the access level and permissions a user has within
// the project context. Users can have different roles across different projects
// within the same organization.
type ProjectUser struct {
	// Object identifies the type of this resource.
	// This will always be "project_user" for ProjectUser objects.
	Object string `json:"object"`

	// ID is the unique identifier of the user.
	// This matches the user's organization-level ID.
	ID string `json:"id"`

	// Name is the user's full name.
	// This is for display purposes and matches their organization profile.
	Name string `json:"name"`

	// Email is the user's email address.
	// This is the primary identifier for the user in the UI.
	Email string `json:"email"`

	// Role defines the user's access level within the project.
	// Valid values are defined by RoleType: "owner", "admin", "developer", "viewer".
	Role string `json:"role"`

	// AddedAt is the timestamp when the user was added to the project.
	// The timestamp is in Unix epoch seconds.
	AddedAt UnixSeconds `json:"added_at"`
}

// ProjectUsersListEndpoint is the base endpoint template for project user management.
// The %s placeholder must be filled with the project ID for all requests.
const ProjectUsersListEndpoint = "/organization/projects/%s/users"

// ListProjectUsers retrieves a paginated list of users in a specific project.
// Results are ordered by when users were added to the project, with most recent first.
//
// Parameters:
//   - projectID: The unique identifier of the project to list users from
//   - limit: Maximum number of users to return (0 for default, which is typically 20)
//   - after: Pagination token for fetching next page (empty string for first page)
//
// Returns a ListResponse containing the project users and pagination metadata.
// The ListResponse includes the next pagination token if more results are available.
func (c *Client) ListProjectUsers(projectID string, limit int, after string) (*ListResponse[ProjectUser], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}

	return Get[ProjectUser](c.client, fmt.Sprintf(ProjectUsersListEndpoint, projectID), queryParams)
}

// CreateProjectUser adds a user to a project with a specified role.
// The user must already be a member of the organization.
//
// Parameters:
//   - projectID: The unique identifier of the project to add the user to
//   - userID: The unique identifier of the user to add
//   - role: The role to assign to the user (must be a valid RoleType)
//
// Returns the created ProjectUser object or an error if the operation fails.
// Common errors include invalid roles, duplicate users, or insufficient permissions.
func (c *Client) CreateProjectUser(projectID string, userID string, role string) (*ProjectUser, error) {
	roleType := ParseRoleType(role)
	if roleType == "" {
		return nil, fmt.Errorf("invalid role: %s", role)
	}
	body := map[string]string{"user_id": userID, "role": string(roleType)}
	return Post[ProjectUser](c.client, fmt.Sprintf(ProjectUsersListEndpoint, projectID), body)
}

// RetrieveProjectUser fetches details about a specific user's membership in a project.
//
// Parameters:
//   - projectID: The unique identifier of the project
//   - userID: The unique identifier of the user to retrieve
//
// Returns the ProjectUser details or an error if retrieval fails.
// Returns an error if the user is not a member of the project or if the caller lacks permission.
func (c *Client) RetrieveProjectUser(projectID string, userID string) (*ProjectUser, error) {
	return GetSingle[ProjectUser](c.client, fmt.Sprintf(ProjectUsersListEndpoint+"/%s", projectID, userID))
}

// ModifyProjectUser updates a user's role within a project.
// This can be used to promote or demote a user's access level.
//
// Parameters:
//   - projectID: The unique identifier of the project
//   - userID: The unique identifier of the user to modify
//   - role: The new role to assign to the user (must be a valid RoleType)
//
// Returns the updated ProjectUser object or an error if modification fails.
// Common errors include invalid roles or insufficient permissions.
func (c *Client) ModifyProjectUser(projectID string, userID string, role string) (*ProjectUser, error) {
	roleType := ParseRoleType(role)
	if roleType == "" {
		return nil, fmt.Errorf("invalid role: %s", role)
	}
	body := map[string]string{"role": string(roleType)}
	return Post[ProjectUser](c.client, fmt.Sprintf(ProjectUsersListEndpoint+"/%s", projectID, userID), body)
}

// DeleteProjectUser removes a user from a project.
// This revokes all of the user's access to the project resources.
//
// Parameters:
//   - projectID: The unique identifier of the project
//   - userID: The unique identifier of the user to remove
//
// Returns an error if the deletion fails or if the caller lacks permission.
// The last owner of a project cannot be removed.
func (c *Client) DeleteProjectUser(projectID string, userID string) error {
	return Delete[ProjectUser](c.client, fmt.Sprintf(ProjectUsersListEndpoint+"/%s", projectID, userID))
}

// String returns a human-readable string representation of the ProjectUser.
// This is useful for logging and debugging purposes.
func (pu *ProjectUser) String() string {
	return fmt.Sprintf("ProjectUser{ID: %s, Name: %s, Email: %s, Role: %s}",
		pu.ID, pu.Name, pu.Email, pu.Role)
}
