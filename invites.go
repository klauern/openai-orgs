package openaiorgs

import (
	"fmt"
)

// Invite represents an invitation for a user to join an OpenAI organization.
// Invites have a limited lifetime and can be in various states (pending, accepted, expired).
type Invite struct {
	// ObjectType identifies the type of this resource.
	ObjectType string `json:"object"`
	// ID is the unique identifier for this invitation.
	ID string `json:"id"`
	// Email is the address of the invited user.
	Email string `json:"email"`
	// Role specifies the permissions the user will have upon accepting.
	Role string `json:"role"`
	// Status indicates the current state of the invitation (e.g., pending, accepted).
	Status string `json:"status"`
	// CreatedAt is the timestamp when this invitation was created.
	CreatedAt UnixSeconds `json:"created_at"`
	// ExpiresAt is the timestamp when this invitation will expire if not accepted.
	ExpiresAt UnixSeconds `json:"expires_at"`
	// AcceptedAt is the timestamp when the invitation was accepted, if applicable.
	AcceptedAt *UnixSeconds `json:"accepted_at,omitempty"`
}

// InviteListEndpoint is the base endpoint for invitation management operations.
const InviteListEndpoint = "/organization/invites"

// ListInvites retrieves all invitations in the organization.
// It automatically handles pagination to fetch all available invites.
//
// Returns a slice of all Invite objects or an error if the retrieval fails.
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

// CreateInvite sends a new invitation to join the organization.
//
// Parameters:
//   - email: The email address of the user to invite
//   - role: The role to assign to the user upon acceptance
//
// Returns the created Invite object or an error if creation fails.
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

// RetrieveInvite fetches details of a specific invitation.
//
// Parameters:
//   - id: The unique identifier of the invitation to retrieve
//
// Returns the Invite details or an error if retrieval fails.
func (c *Client) RetrieveInvite(id string) (*Invite, error) {
	resp, err := GetSingle[Invite](c.client, InviteListEndpoint+"/"+id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve invite: %w", err)
	}

	return resp, nil
}

// DeleteInvite cancels a pending invitation.
//
// Parameters:
//   - id: The unique identifier of the invitation to delete
//
// Returns an error if deletion fails or nil on success.
func (c *Client) DeleteInvite(id string) error {
	err := Delete[Invite](c.client, InviteListEndpoint+"/"+id)
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}

// String returns a human-readable string representation of the Invite.
// It includes the invitation's ID, email, role, status, and acceptance time if applicable.
func (i *Invite) String() string {
	acceptedInfo := ""
	if i.AcceptedAt != nil {
		acceptedInfo = fmt.Sprintf(", Accepted: %s", i.AcceptedAt.String())
	}
	return fmt.Sprintf("Invite{ID: %s, Email: %s, Role: %s, Status: %s%s}",
		i.ID, i.Email, i.Role, i.Status, acceptedInfo)
}
