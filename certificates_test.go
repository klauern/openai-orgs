package openaiorgs

import (
	"testing"
	"time"
)

func TestListOrganizationCertificates(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	activeTrue := true
	mockCertificates := []Certificate{
		{
			Object:    "certificate",
			ID:        "cert_123",
			Name:      "Test Certificate",
			Active:    &activeTrue,
			CreatedAt: UnixSeconds(now),
			CertificateDetails: CertificateDetails{
				ValidAt:   UnixSeconds(now),
				ExpiresAt: UnixSeconds(now.Add(365 * 24 * time.Hour)),
			},
		},
	}

	// Register mock response
	response := ListResponse[Certificate]{
		Object:  "list",
		Data:    mockCertificates,
		FirstID: "cert_123",
		LastID:  "cert_123",
		HasMore: false,
	}
	h.mockResponse("GET", OrganizationCertificatesEndpoint, 200, response)

	// Make the API call
	certificates, err := h.client.ListOrganizationCertificates(10, "", "desc")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(certificates.Data) != 1 {
		t.Errorf("Expected 1 certificate, got %d", len(certificates.Data))
		return
	}
	if mockCertificates[0].ID != certificates.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockCertificates[0].ID, certificates.Data[0].ID)
	}
	if mockCertificates[0].Name != certificates.Data[0].Name {
		t.Errorf("Expected Name %s, got %s", mockCertificates[0].Name, certificates.Data[0].Name)
	}
	if *mockCertificates[0].Active != *certificates.Data[0].Active {
		t.Errorf("Expected Active %v, got %v", *mockCertificates[0].Active, *certificates.Data[0].Active)
	}
}

func TestUploadCertificate(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockCertificate := Certificate{
		Object:    "certificate",
		ID:        "cert_123",
		Name:      "Test Certificate",
		CreatedAt: UnixSeconds(now),
		CertificateDetails: CertificateDetails{
			ValidAt:   UnixSeconds(now),
			ExpiresAt: UnixSeconds(now.Add(365 * 24 * time.Hour)),
		},
	}

	// Register mock response
	h.mockResponse("POST", OrganizationCertificatesEndpoint, 200, mockCertificate)

	// Make the API call
	certificate, err := h.client.UploadCertificate("-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----", "Test Certificate")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if mockCertificate.ID != certificate.ID {
		t.Errorf("Expected ID %s, got %s", mockCertificate.ID, certificate.ID)
	}
	if mockCertificate.Name != certificate.Name {
		t.Errorf("Expected Name %s, got %s", mockCertificate.Name, certificate.Name)
	}
}

func TestGetCertificate(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockCertificate := Certificate{
		Object:    "certificate",
		ID:        "cert_123",
		Name:      "Test Certificate",
		CreatedAt: UnixSeconds(now),
		CertificateDetails: CertificateDetails{
			ValidAt:   UnixSeconds(now),
			ExpiresAt: UnixSeconds(now.Add(365 * 24 * time.Hour)),
		},
	}

	// Register mock response
	h.mockResponse("GET", OrganizationCertificatesEndpoint+"/cert_123", 200, mockCertificate)

	// Make the API call
	certificate, err := h.client.GetCertificate("cert_123", false)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if mockCertificate.ID != certificate.ID {
		t.Errorf("Expected ID %s, got %s", mockCertificate.ID, certificate.ID)
	}
	if mockCertificate.Name != certificate.Name {
		t.Errorf("Expected Name %s, got %s", mockCertificate.Name, certificate.Name)
	}
}

func TestGetCertificateWithContent(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	certContent := "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----"
	mockCertificate := Certificate{
		Object:    "certificate",
		ID:        "cert_123",
		Name:      "Test Certificate",
		CreatedAt: UnixSeconds(now),
		CertificateDetails: CertificateDetails{
			ValidAt:   UnixSeconds(now),
			ExpiresAt: UnixSeconds(now.Add(365 * 24 * time.Hour)),
			Content:   &certContent,
		},
	}

	// Register mock response
	h.mockResponse("GET", OrganizationCertificatesEndpoint+"/cert_123", 200, mockCertificate)

	// Make the API call
	certificate, err := h.client.GetCertificate("cert_123", true)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if mockCertificate.ID != certificate.ID {
		t.Errorf("Expected ID %s, got %s", mockCertificate.ID, certificate.ID)
	}
	if certificate.CertificateDetails.Content == nil {
		t.Error("Expected certificate content to be present")
		return
	}
	if *certificate.CertificateDetails.Content != certContent {
		t.Errorf("Expected content %s, got %s", certContent, *certificate.CertificateDetails.Content)
	}
}

func TestModifyCertificate(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockCertificate := Certificate{
		Object:    "certificate",
		ID:        "cert_123",
		Name:      "Updated Certificate",
		CreatedAt: UnixSeconds(now),
		CertificateDetails: CertificateDetails{
			ValidAt:   UnixSeconds(now),
			ExpiresAt: UnixSeconds(now.Add(365 * 24 * time.Hour)),
		},
	}

	// Register mock response
	h.mockResponse("POST", OrganizationCertificatesEndpoint+"/cert_123", 200, mockCertificate)

	// Make the API call
	certificate, err := h.client.ModifyCertificate("cert_123", "Updated Certificate")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if mockCertificate.ID != certificate.ID {
		t.Errorf("Expected ID %s, got %s", mockCertificate.ID, certificate.ID)
	}
	if mockCertificate.Name != certificate.Name {
		t.Errorf("Expected Name %s, got %s", mockCertificate.Name, certificate.Name)
	}
}

func TestDeleteCertificate(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResponse := CertificateDeletedResponse{
		Object:  "certificate.deleted",
		ID:      "cert_123",
		Deleted: true,
	}

	// Register mock response
	h.mockResponse("DELETE", OrganizationCertificatesEndpoint+"/cert_123", 200, mockResponse)

	// Make the API call
	response, err := h.client.DeleteCertificate("cert_123")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if mockResponse.ID != response.ID {
		t.Errorf("Expected ID %s, got %s", mockResponse.ID, response.ID)
	}
	if mockResponse.Deleted != response.Deleted {
		t.Errorf("Expected Deleted %v, got %v", mockResponse.Deleted, response.Deleted)
	}
}

func TestActivateOrganizationCertificates(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResponse := CertificateActivationResponse{
		Object:  "certificate.activation",
		Success: true,
	}

	// Register mock response
	h.mockResponse("POST", OrganizationCertificateActivateEndpoint, 200, mockResponse)

	// Make the API call
	response, err := h.client.ActivateOrganizationCertificates([]string{"cert_123", "cert_456"})
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if mockResponse.Success != response.Success {
		t.Errorf("Expected Success %v, got %v", mockResponse.Success, response.Success)
	}
}

func TestDeactivateOrganizationCertificates(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResponse := CertificateActivationResponse{
		Object:  "certificate.activation",
		Success: true,
	}

	// Register mock response
	h.mockResponse("POST", OrganizationCertificateDeactivateEndpoint, 200, mockResponse)

	// Make the API call
	response, err := h.client.DeactivateOrganizationCertificates([]string{"cert_123", "cert_456"})
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if mockResponse.Success != response.Success {
		t.Errorf("Expected Success %v, got %v", mockResponse.Success, response.Success)
	}
}

func TestListProjectCertificates(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	activeTrue := true
	mockCertificates := []Certificate{
		{
			Object:    "organization.project.certificate",
			ID:        "cert_123",
			Name:      "Test Certificate",
			Active:    &activeTrue,
			CreatedAt: UnixSeconds(now),
			CertificateDetails: CertificateDetails{
				ValidAt:   UnixSeconds(now),
				ExpiresAt: UnixSeconds(now.Add(365 * 24 * time.Hour)),
			},
		},
	}

	// Register mock response
	response := ListResponse[Certificate]{
		Object:  "list",
		Data:    mockCertificates,
		FirstID: "cert_123",
		LastID:  "cert_123",
		HasMore: false,
	}
	h.mockResponse("GET", "/organization/projects/proj_123/certificates", 200, response)

	// Make the API call
	certificates, err := h.client.ListProjectCertificates("proj_123", 10, "", "desc")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(certificates.Data) != 1 {
		t.Errorf("Expected 1 certificate, got %d", len(certificates.Data))
		return
	}
	if mockCertificates[0].ID != certificates.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockCertificates[0].ID, certificates.Data[0].ID)
	}
	if mockCertificates[0].Object != certificates.Data[0].Object {
		t.Errorf("Expected Object %s, got %s", mockCertificates[0].Object, certificates.Data[0].Object)
	}
}

func TestActivateProjectCertificates(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResponse := CertificateActivationResponse{
		Object:  "certificate.activation",
		Success: true,
	}

	// Register mock response
	h.mockResponse("POST", "/organization/projects/proj_123/certificates/activate", 200, mockResponse)

	// Make the API call
	response, err := h.client.ActivateProjectCertificates("proj_123", []string{"cert_123", "cert_456"})
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if mockResponse.Success != response.Success {
		t.Errorf("Expected Success %v, got %v", mockResponse.Success, response.Success)
	}
}

func TestDeactivateProjectCertificates(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	mockResponse := CertificateActivationResponse{
		Object:  "certificate.activation",
		Success: true,
	}

	// Register mock response
	h.mockResponse("POST", "/organization/projects/proj_123/certificates/deactivate", 200, mockResponse)

	// Make the API call
	response, err := h.client.DeactivateProjectCertificates("proj_123", []string{"cert_123", "cert_456"})
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if mockResponse.Success != response.Success {
		t.Errorf("Expected Success %v, got %v", mockResponse.Success, response.Success)
	}
}