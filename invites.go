package openaiorgs

import (
	"fmt"
)

type Invite struct {
	ObjectType string       `json:"object"`
	ID         string       `json:"id"`
	Email      string       `json:"email"`
	Role       string       `json:"role"`
	Status     string       `json:"status"`
	CreatedAt  UnixSeconds  `json:"created_at"`
	ExpiresAt  UnixSeconds  `json:"expires_at"`
	AcceptedAt *UnixSeconds `json:"accepted_at,omitempty"`
}

const InviteListEndpoint = "/organization/invites"

func (c *Client) ListInvites() ([]Invite, error) {
	resp, err := Get[Invite](c.client, InviteListEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get invites: %w", err)
	}

	return resp.Data, nil
}

func (c *Client) CreateInvite(email string, role string) (*Invite, error) {
	body := map[string]string{
		"email": email,
		"role":  role,
	}

	invite, err := Post[Invite](c.client, InviteListEndpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create invite: %w", err)
	}

	return invite, nil
}

func (c *Client) RetrieveInvite(id string) (*Invite, error) {
	resp, err := GetSingle[Invite](c.client, InviteListEndpoint+"/"+id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve invite: %w", err)
	}

	return resp, nil
}

func (c *Client) DeleteInvite(id string) error {
	err := Delete[Invite](c.client, InviteListEndpoint+"/"+id)
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}
