package oaiprom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type CustomTime time.Time

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		// Handle string format
		t, err := time.Parse(time.RFC3339, string(b[1:len(b)-1]))
		if err != nil {
			return err
		}
		*ct = CustomTime(t)
	} else {
		// Handle numeric format (assume Unix timestamp in seconds)
		var timestamp int64
		err := json.Unmarshal(b, &timestamp)
		if err != nil {
			return err
		}
		*ct = CustomTime(time.Unix(timestamp, 0))
	}
	return nil
}

// Add this method to the CustomTime type
func (ct CustomTime) String() string {
	return time.Time(ct).Format(time.RFC3339)
}

type Invite struct {
	ObjectType string      `json:"object"`
	ID         string      `json:"id"`
	Email      string      `json:"email"`
	Role       string      `json:"role"`
	Status     string      `json:"status"`
	CreatedAt  CustomTime  `json:"created_at"`
	ExpiresAt  CustomTime  `json:"expires_at"`
	AcceptedAt *CustomTime `json:"accepted_at,omitempty"`
}

const InviteListEndpoint = "/organization/invites"

func (c *Client) ListInvites() ([]Invite, error) {
	// Get the raw response
	rawResp, err := c.client.R().Get(InviteListEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get invites: %w", err)
	}

	// Read and log the raw response body
	body, err := io.ReadAll(bytes.NewReader(rawResp.Body()))
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Attempt to parse the response
	var resp ListResponse[Invite]
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp.Data, nil
}

func (c *Client) CreateInvite(email string, role RoleType) (*Invite, error) {
	body := map[string]string{
		"email": email,
		"role":  string(role),
	}

	invite, err := Post[Invite](c.client, InviteListEndpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create invite: %w", err)
	}

	return invite, nil
}

func (c *Client) RetrieveInvite(id string) (*Invite, error) {
	resp, err := Get[Invite](c.client, InviteListEndpoint+"/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve invite: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no invite found with id: %s", id)
	}

	return &resp.Data[0], nil
}

func (c *Client) DeleteInvite(id string) error {
	err := Delete[Invite](c.client, InviteListEndpoint+"/"+id)
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}
