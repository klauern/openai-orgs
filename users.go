package openaiorgs

import (
	"fmt"
)

// User represents a user in the OpenAI organization
type User struct {
	Object  string      `json:"object"`
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Email   string      `json:"email"`
	Role    string      `json:"role"`
	AddedAt UnixSeconds `json:"added_at"`
}

const UsersListEndpoint = "/organization/users"

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

func (c *Client) RetrieveUser(id string) (*User, error) {
	return GetSingle[User](c.client, UsersListEndpoint+"/"+id)
}

func (c *Client) DeleteUser(id string) error {
	err := Delete[User](c.client, UsersListEndpoint+"/"+id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

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

// String returns a human-readable string representation of the User
func (u *User) String() string {
	return fmt.Sprintf("User{ID: %s, Name: %s, Email: %s, Role: %s}",
		u.ID, u.Name, u.Email, u.Role)
}
