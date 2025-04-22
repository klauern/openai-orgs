package openaiorgs

import (
	"fmt"
)

// User represents a user account within an OpenAI organization.
// Users can have different roles and permissions within the organization.
type User struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`
	// ID is the unique identifier for this user.
	ID string `json:"id"`
	// Name is the user's full name.
	Name string `json:"name"`
	// Email is the user's email address, used for authentication and communication.
	Email string `json:"email"`
	// Role defines the user's permissions within the organization.
	Role string `json:"role"`
	// AddedAt is the timestamp when this user was added to the organization.
	AddedAt UnixSeconds `json:"added_at"`
}

// UsersListEndpoint is the base endpoint for user management operations.
const UsersListEndpoint = "/organization/users"

// ListUsers retrieves a paginated list of users in the organization.
//
// Parameters:
//   - limit: Maximum number of users to return (0 for default)
//   - after: Pagination token for fetching next page (empty string for first page)
//
// Returns a ListResponse containing the users and pagination metadata.
// Returns an error if the API request fails.
func (c *Client) ListUsers(limit int, after string) (*ListResponse[User], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}

	return Get[User](c.client, UsersListEndpoint, queryParams)
}

// RetrieveUser fetches details of a specific user in the organization.
//
// Parameters:
//   - id: The unique identifier of the user to retrieve
//
// Returns the User details or an error if retrieval fails.
func (c *Client) RetrieveUser(id string) (*User, error) {
	return GetSingle[User](c.client, UsersListEndpoint+"/"+id)
}

// DeleteUser removes a user from the organization.
//
// Parameters:
//   - id: The unique identifier of the user to delete
//
// Returns an error if deletion fails or nil on success.
func (c *Client) DeleteUser(id string) error {
	err := Delete[User](c.client, UsersListEndpoint+"/"+id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ModifyUserRole updates a user's role within the organization.
//
// Parameters:
//   - id: The unique identifier of the user to modify
//   - role: The new role to assign to the user
//
// Returns an error if the role modification fails or nil on success.
func (c *Client) ModifyUserRole(id string, role string) error {
	body := map[string]string{
		"role": role,
	}

	_, err := Post[User](c.client, UsersListEndpoint+"/"+id, body)
	if err != nil {
		return fmt.Errorf("failed to modify user role: %w", err)
	}

	return nil
}

// String returns a human-readable string representation of the User.
// It includes the user's ID, name, email, and role within the organization.
func (u *User) String() string {
	return fmt.Sprintf("User{ID: %s, Name: %s, Email: %s, Role: %s}",
		u.ID, u.Name, u.Email, u.Role)
}
