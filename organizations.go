package openaiorgs

import (
	"fmt"
)

// Organization represents an OpenAI organization
type Organization struct {
	ID          string      `json:"id"`
	Object      string      `json:"object"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CreatedAt   UnixSeconds `json:"created"`
	IsBlocked   bool        `json:"is_blocked"`
	Settings    struct {
		DefaultBillingAddress *BillingAddress `json:"default_billing_address,omitempty"`
	} `json:"settings"`
}

// BillingAddress represents a billing address for an organization
type BillingAddress struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

// UpdateOrganizationRequest represents the request body for updating an organization
type UpdateOrganizationRequest struct {
	Name        *string      `json:"name,omitempty"`
	Description *string      `json:"description,omitempty"`
	Settings    *OrgSettings `json:"settings,omitempty"`
}

type OrgSettings struct {
	DefaultBillingAddress *BillingAddress `json:"default_billing_address,omitempty"`
}

// ListOrganizations retrieves a list of organizations
func (c *Client) ListOrganizations(limit int, after string) (*ListResponse[Organization], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}

	return Get[Organization](c.client, "/organizations", queryParams)
}

// GetOrganization retrieves a specific organization by ID
func (c *Client) GetOrganization(id string) (*Organization, error) {
	return GetSingle[Organization](c.client, fmt.Sprintf("/organizations/%s", id))
}

// UpdateOrganization updates an organization's settings
func (c *Client) UpdateOrganization(id string, req *UpdateOrganizationRequest) (*Organization, error) {
	return Post[Organization](c.client, fmt.Sprintf("/organizations/%s", id), req)
}

// DeleteOrganization deletes an organization
func (c *Client) DeleteOrganization(id string) error {
	return Delete[Organization](c.client, fmt.Sprintf("/organizations/%s", id))
}
