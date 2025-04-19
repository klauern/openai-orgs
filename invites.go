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
	var allInvites []Invite
	queryParams := map[string]string{
		"limit": "100",
	}

	for {
		resp, err := Get[Invite](c.client, InviteListEndpoint, queryParams)
		if err != nil {
			return nil, fmt.Errorf("failed to get invites: %w", err)
		}

		allInvites = append(allInvites, resp.Data...)

		if !resp.HasMore {
			break
		}

		fmt.Println("Getting more invites after", resp.LastID)
		queryParams["after"] = resp.LastID
	}

	return allInvites, nil
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

// String returns a human-readable string representation of the Invite
func (i *Invite) String() string {
	acceptedInfo := ""
	if i.AcceptedAt != nil {
		acceptedInfo = fmt.Sprintf(", Accepted: %s", i.AcceptedAt.String())
	}
	return fmt.Sprintf("Invite{ID: %s, Email: %s, Role: %s, Status: %s%s}",
		i.ID, i.Email, i.Role, i.Status, acceptedInfo)
}
