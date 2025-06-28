package openaiorgs

import (
	"encoding/json"
	"fmt"
)

// ListOrganizationCertificates retrieves a paginated list of certificates in the organization.
// Certificates are used for mutual TLS authentication and can be managed at the organization level.
// Results are ordered by creation date, with newest certificates first.
//
// Parameters:
//   - limit: Maximum number of certificates to return (0 for default, which is typically 20)
//   - after: Pagination token for fetching next page (empty string for first page)
//   - order: Sort order for results ("asc" or "desc", defaults to "desc")
//
// Returns a ListResponse containing the certificates and pagination metadata.
// The ListResponse includes the next pagination token if more results are available.
func (c *Client) ListOrganizationCertificates(limit int, after string, order string) (*ListResponse[Certificate], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}
	if order != "" {
		queryParams["order"] = order
	}

	return Get[Certificate](c.client, OrganizationCertificatesEndpoint, queryParams)
}

// UploadCertificate uploads a new certificate to the organization.
// The certificate content must be provided in PEM format.
// New certificates are created in an inactive state and must be explicitly activated.
//
// Parameters:
//   - content: The PEM-encoded certificate content
//   - name: The human-readable name for the certificate (must be unique)
//
// Returns the created Certificate object or an error if upload fails.
// Common errors include invalid certificate format or duplicate names.
func (c *Client) UploadCertificate(content string, name string) (*Certificate, error) {
	body := map[string]string{
		"content": content,
		"name":    name,
	}
	return Post[Certificate](c.client, OrganizationCertificatesEndpoint, body)
}

// GetCertificate fetches details of a specific certificate.
// This can be used to get the current state of any certificate in the organization.
//
// Parameters:
//   - certificateID: The unique identifier of the certificate to retrieve
//   - includeContent: Whether to include the PEM certificate content in the response
//
// Returns the Certificate details or an error if retrieval fails.
// Returns an error if the certificate ID does not exist or if the caller lacks permission.
func (c *Client) GetCertificate(certificateID string, includeContent bool) (*Certificate, error) {
	endpoint := OrganizationCertificatesEndpoint + "/" + certificateID
	if includeContent {
		queryParams := map[string]string{"include": "content"}
		resp, err := c.client.R().
			SetQueryParams(queryParams).
			ExpectContentType("application/json").
			Get(endpoint)
		if err != nil {
			return nil, fmt.Errorf("error making GET request: %v", err)
		}

		if resp.IsError() {
			return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
		}

		var result Certificate
		err = json.Unmarshal(resp.Body(), &result)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %v", err)
		}

		return &result, nil
	}
	return GetSingle[Certificate](c.client, endpoint)
}

// ModifyCertificate updates the properties of an existing certificate.
// Currently, only the certificate name can be modified.
// This operation cannot be performed on deleted certificates.
//
// Parameters:
//   - certificateID: The unique identifier of the certificate to modify
//   - name: The new name for the certificate (must be unique within the organization)
//
// Returns the updated Certificate object or an error if modification fails.
// Common errors include duplicate names or attempting to modify a deleted certificate.
func (c *Client) ModifyCertificate(certificateID string, name string) (*Certificate, error) {
	body := map[string]string{"name": name}
	endpoint := OrganizationCertificatesEndpoint + "/" + certificateID
	return Post[Certificate](c.client, endpoint, body)
}

// DeleteCertificate removes a certificate from the organization.
// Deleted certificates cannot be recovered and will be immediately deactivated if they were active.
// This operation cannot be undone.
//
// Parameters:
//   - certificateID: The unique identifier of the certificate to delete
//
// Returns a CertificateDeletedResponse confirming the deletion or an error if deletion fails.
// Returns an error if the certificate doesn't exist or if the caller lacks permission.
func (c *Client) DeleteCertificate(certificateID string) (*CertificateDeletedResponse, error) {
	endpoint := OrganizationCertificatesEndpoint + "/" + certificateID
	resp, err := c.client.R().
		ExpectContentType("application/json").
		Delete(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making DELETE request: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	var result CertificateDeletedResponse
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return &result, nil
}

// ActivateOrganizationCertificates activates multiple certificates at the organization level.
// This is a bulk operation that atomically activates all specified certificates.
// All certificates must exist and be owned by the organization.
//
// Parameters:
//   - certificateIDs: List of certificate IDs to activate
//
// Returns a CertificateActivationResponse indicating success or failure.
// If any certificate fails to activate, the entire operation is rolled back.
func (c *Client) ActivateOrganizationCertificates(certificateIDs []string) (*CertificateActivationResponse, error) {
	body := map[string][]string{"certificate_ids": certificateIDs}
	return Post[CertificateActivationResponse](c.client, OrganizationCertificateActivateEndpoint, body)
}

// DeactivateOrganizationCertificates deactivates multiple certificates at the organization level.
// This is a bulk operation that atomically deactivates all specified certificates.
// All certificates must exist and be owned by the organization.
//
// Parameters:
//   - certificateIDs: List of certificate IDs to deactivate
//
// Returns a CertificateActivationResponse indicating success or failure.
// If any certificate fails to deactivate, the entire operation is rolled back.
func (c *Client) DeactivateOrganizationCertificates(certificateIDs []string) (*CertificateActivationResponse, error) {
	body := map[string][]string{"certificate_ids": certificateIDs}
	return Post[CertificateActivationResponse](c.client, OrganizationCertificateDeactivateEndpoint, body)
}

// ListProjectCertificates retrieves a paginated list of certificates available to a specific project.
// This includes both organization-level certificates and project-specific certificates.
// Results are ordered by creation date, with newest certificates first.
//
// Parameters:
//   - projectID: The unique identifier of the project
//   - limit: Maximum number of certificates to return (0 for default, which is typically 20)
//   - after: Pagination token for fetching next page (empty string for first page)
//   - order: Sort order for results ("asc" or "desc", defaults to "desc")
//
// Returns a ListResponse containing the certificates and pagination metadata.
// The ListResponse includes the next pagination token if more results are available.
func (c *Client) ListProjectCertificates(projectID string, limit int, after string, order string) (*ListResponse[Certificate], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}
	if order != "" {
		queryParams["order"] = order
	}

	endpoint := fmt.Sprintf(ProjectCertificatesEndpoint, projectID)
	return Get[Certificate](c.client, endpoint, queryParams)
}

// ActivateProjectCertificates activates multiple certificates for a specific project.
// This is a bulk operation that atomically activates all specified certificates for the project.
// All certificates must exist and be accessible to the project.
//
// Parameters:
//   - projectID: The unique identifier of the project
//   - certificateIDs: List of certificate IDs to activate for the project
//
// Returns a CertificateActivationResponse indicating success or failure.
// If any certificate fails to activate, the entire operation is rolled back.
func (c *Client) ActivateProjectCertificates(projectID string, certificateIDs []string) (*CertificateActivationResponse, error) {
	body := map[string][]string{"certificate_ids": certificateIDs}
	endpoint := fmt.Sprintf(ProjectCertificateActivateEndpoint, projectID)
	return Post[CertificateActivationResponse](c.client, endpoint, body)
}

// DeactivateProjectCertificates deactivates multiple certificates for a specific project.
// This is a bulk operation that atomically deactivates all specified certificates for the project.
// All certificates must exist and be accessible to the project.
//
// Parameters:
//   - projectID: The unique identifier of the project
//   - certificateIDs: List of certificate IDs to deactivate for the project
//
// Returns a CertificateActivationResponse indicating success or failure.
// If any certificate fails to deactivate, the entire operation is rolled back.
func (c *Client) DeactivateProjectCertificates(projectID string, certificateIDs []string) (*CertificateActivationResponse, error) {
	body := map[string][]string{"certificate_ids": certificateIDs}
	endpoint := fmt.Sprintf(ProjectCertificateDeactivateEndpoint, projectID)
	return Post[CertificateActivationResponse](c.client, endpoint, body)
}