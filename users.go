package openaiorgs

import (
	"fmt"
)

// Users represents a user in the OpenAI organization
type Users struct {
	Object  string     `json:"object"`
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Email   string     `json:"email"`
	Role    string     `json:"role"`
	AddedAt CustomTime `json:"added_at"`
}

const UsersListEndpoint = "/organization/users"

func (c *Client) ListUsers(limit int, after string) (*ListResponse[Users], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}

	return Get[Users](c.client, UsersListEndpoint, queryParams)
}

func (c *Client) RetrieveUser(id string) (*Users, error) {
	resp, err := Get[Users](c.client, UsersListEndpoint+"/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return &resp.Data[0], nil
}

func (c *Client) DeleteUser(id string) error {
	err := Delete[Users](c.client, UsersListEndpoint+"/"+id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (c *Client) ModifyUserRole(id string, role RoleType) error {
	body := map[string]string{
		"role": string(role),
	}

	_, err := Post[Users](c.client, UsersListEndpoint+"/"+id, body)
	if err != nil {
		return fmt.Errorf("failed to modify user role: %w", err)
	}

	return nil
}
