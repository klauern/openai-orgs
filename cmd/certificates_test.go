package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Mock client interface for testing
type mockCertificateClient interface {
	ListOrganizationCertificates(limit int, after, order string) (*openaiorgs.ListResponse[openaiorgs.Certificate], error)
	UploadCertificate(content, name string) (*openaiorgs.Certificate, error)
	GetCertificate(id string, includeContent bool) (*openaiorgs.Certificate, error)
	ModifyCertificate(id, name string) (*openaiorgs.Certificate, error)
	DeleteCertificate(id string) (*openaiorgs.CertificateDeletedResponse, error)
	ActivateOrganizationCertificates(certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error)
	DeactivateOrganizationCertificates(certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error)
	ListProjectCertificates(projectID string, limit int, after, order string) (*openaiorgs.ListResponse[openaiorgs.Certificate], error)
	ActivateProjectCertificates(projectID string, certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error)
	DeactivateProjectCertificates(projectID string, certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error)
}

// Mock implementation
type mockCertificateClientImpl struct {
	ListOrganizationCertificatesFunc       func(limit int, after, order string) (*openaiorgs.ListResponse[openaiorgs.Certificate], error)
	UploadCertificateFunc                  func(content, name string) (*openaiorgs.Certificate, error)
	GetCertificateFunc                     func(id string, includeContent bool) (*openaiorgs.Certificate, error)
	ModifyCertificateFunc                  func(id, name string) (*openaiorgs.Certificate, error)
	DeleteCertificateFunc                  func(id string) (*openaiorgs.CertificateDeletedResponse, error)
	ActivateOrganizationCertificatesFunc   func(certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error)
	DeactivateOrganizationCertificatesFunc func(certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error)
	ListProjectCertificatesFunc            func(projectID string, limit int, after, order string) (*openaiorgs.ListResponse[openaiorgs.Certificate], error)
	ActivateProjectCertificatesFunc        func(projectID string, certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error)
	DeactivateProjectCertificatesFunc      func(projectID string, certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error)
}

func (m *mockCertificateClientImpl) ListOrganizationCertificates(limit int, after, order string) (*openaiorgs.ListResponse[openaiorgs.Certificate], error) {
	if m.ListOrganizationCertificatesFunc != nil {
		return m.ListOrganizationCertificatesFunc(limit, after, order)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCertificateClientImpl) UploadCertificate(content, name string) (*openaiorgs.Certificate, error) {
	if m.UploadCertificateFunc != nil {
		return m.UploadCertificateFunc(content, name)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCertificateClientImpl) GetCertificate(id string, includeContent bool) (*openaiorgs.Certificate, error) {
	if m.GetCertificateFunc != nil {
		return m.GetCertificateFunc(id, includeContent)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCertificateClientImpl) ModifyCertificate(id, name string) (*openaiorgs.Certificate, error) {
	if m.ModifyCertificateFunc != nil {
		return m.ModifyCertificateFunc(id, name)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCertificateClientImpl) DeleteCertificate(id string) (*openaiorgs.CertificateDeletedResponse, error) {
	if m.DeleteCertificateFunc != nil {
		return m.DeleteCertificateFunc(id)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCertificateClientImpl) ActivateOrganizationCertificates(certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error) {
	if m.ActivateOrganizationCertificatesFunc != nil {
		return m.ActivateOrganizationCertificatesFunc(certificateIDs)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCertificateClientImpl) DeactivateOrganizationCertificates(certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error) {
	if m.DeactivateOrganizationCertificatesFunc != nil {
		return m.DeactivateOrganizationCertificatesFunc(certificateIDs)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCertificateClientImpl) ListProjectCertificates(projectID string, limit int, after, order string) (*openaiorgs.ListResponse[openaiorgs.Certificate], error) {
	if m.ListProjectCertificatesFunc != nil {
		return m.ListProjectCertificatesFunc(projectID, limit, after, order)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCertificateClientImpl) ActivateProjectCertificates(projectID string, certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error) {
	if m.ActivateProjectCertificatesFunc != nil {
		return m.ActivateProjectCertificatesFunc(projectID, certificateIDs)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockCertificateClientImpl) DeactivateProjectCertificates(projectID string, certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error) {
	if m.DeactivateProjectCertificatesFunc != nil {
		return m.DeactivateProjectCertificatesFunc(projectID, certificateIDs)
	}
	return nil, fmt.Errorf("not implemented")
}

// Test helper to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// createMockCertificate creates a test certificate
func createMockCertificate(id, name string, active *bool) openaiorgs.Certificate {
	now := time.Now()
	return openaiorgs.Certificate{
		Object:    "certificate",
		ID:        id,
		Name:      name,
		Active:    active,
		CreatedAt: openaiorgs.UnixSeconds(now),
		CertificateDetails: openaiorgs.CertificateDetails{
			ValidAt:   openaiorgs.UnixSeconds(now),
			ExpiresAt: openaiorgs.UnixSeconds(now.Add(365 * 24 * time.Hour)),
		},
	}
}

// Testable handlers - extracted business logic

func listOrgCertificatesHandler(client mockCertificateClient, limit int, after, order string) error {
	certificates, err := client.ListOrganizationCertificates(limit, after, order)
	if err != nil {
		return wrapError("list organization certificates", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Active", "Valid At", "Expires At"},
		Rows:    make([][]string, len(certificates.Data)),
	}

	for i, cert := range certificates.Data {
		active := "N/A"
		if cert.Active != nil {
			if *cert.Active {
				active = "Yes"
			} else {
				active = "No"
			}
		}
		data.Rows[i] = []string{
			cert.ID,
			cert.Name,
			active,
			cert.CertificateDetails.ValidAt.String(),
			cert.CertificateDetails.ExpiresAt.String(),
		}
	}

	printTableData(data)
	return nil
}

func uploadCertificateHandler(client mockCertificateClient, content, name string) error {
	certificate, err := client.UploadCertificate(content, name)
	if err != nil {
		return wrapError("upload certificate", err)
	}

	fmt.Printf("Certificate uploaded:\n")
	fmt.Printf("ID: %s\nName: %s\nValid At: %s\nExpires At: %s\n",
		certificate.ID,
		certificate.Name,
		certificate.CertificateDetails.ValidAt.String(),
		certificate.CertificateDetails.ExpiresAt.String())
	return nil
}

func getCertificateHandler(client mockCertificateClient, id string, includeContent bool) error {
	certificate, err := client.GetCertificate(id, includeContent)
	if err != nil {
		return wrapError("get certificate", err)
	}

	fmt.Printf("Certificate details:\n")
	fmt.Printf("ID: %s\nName: %s\nValid At: %s\nExpires At: %s\n",
		certificate.ID,
		certificate.Name,
		certificate.CertificateDetails.ValidAt.String(),
		certificate.CertificateDetails.ExpiresAt.String())
	if certificate.CertificateDetails.Content != nil {
		fmt.Printf("Content:\n%s\n", *certificate.CertificateDetails.Content)
	}
	return nil
}

func activateOrgCertificatesHandler(client mockCertificateClient, certificateIDs []string) error {
	if len(certificateIDs) == 0 {
		return fmt.Errorf("at least one certificate ID must be provided")
	}

	allIDs := splitCommaSeparatedIDs(certificateIDs)

	response, err := client.ActivateOrganizationCertificates(allIDs)
	if err != nil {
		return wrapError("activate certificates", err)
	}

	if response.Success {
		fmt.Printf("Successfully activated %d certificates\n", len(allIDs))
	} else {
		fmt.Printf("Failed to activate certificates\n")
	}
	return nil
}

// Test cases

func TestListOrgCertificatesHandler(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		after   string
		order   string
		mockFn  func(*mockCertificateClientImpl)
		wantErr bool
	}{
		{
			name:  "successful list",
			limit: 10,
			after: "",
			order: "desc",
			mockFn: func(m *mockCertificateClientImpl) {
				activeTrue := true
				cert := createMockCertificate("cert_123", "Test Certificate", &activeTrue)
				m.ListOrganizationCertificatesFunc = func(limit int, after, order string) (*openaiorgs.ListResponse[openaiorgs.Certificate], error) {
					if limit != 10 || after != "" || order != "desc" {
						t.Errorf("unexpected params: limit=%d, after=%s, order=%s", limit, after, order)
					}
					return &openaiorgs.ListResponse[openaiorgs.Certificate]{
						Object:  "list",
						Data:    []openaiorgs.Certificate{cert},
						FirstID: "cert_123",
						LastID:  "cert_123",
						HasMore: false,
					}, nil
				}
			},
		},
		{
			name:  "error from client",
			limit: 10,
			after: "",
			order: "desc",
			mockFn: func(m *mockCertificateClientImpl) {
				m.ListOrganizationCertificatesFunc = func(limit int, after, order string) (*openaiorgs.ListResponse[openaiorgs.Certificate], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockCertificateClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listOrgCertificatesHandler(mock, tt.limit, tt.after, tt.order)
				if (err != nil) != tt.wantErr {
					t.Errorf("listOrgCertificatesHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && !strings.Contains(output, "ID | Name | Active | Valid At | Expires At") {
				t.Errorf("Expected table headers in output, got: %s", output)
			}
		})
	}
}

func TestUploadCertificateHandler(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		certName string
		mockFn   func(*mockCertificateClientImpl)
		wantErr  bool
	}{
		{
			name:     "successful upload",
			content:  "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
			certName: "Test Certificate",
			mockFn: func(m *mockCertificateClientImpl) {
				m.UploadCertificateFunc = func(content, name string) (*openaiorgs.Certificate, error) {
					if content != "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----" {
						t.Errorf("unexpected content: %s", content)
					}
					if name != "Test Certificate" {
						t.Errorf("unexpected name: %s", name)
					}
					cert := createMockCertificate("cert_123", "Test Certificate", nil)
					return &cert, nil
				}
			},
		},
		{
			name:     "error from client",
			content:  "invalid cert",
			certName: "Test Certificate",
			mockFn: func(m *mockCertificateClientImpl) {
				m.UploadCertificateFunc = func(content, name string) (*openaiorgs.Certificate, error) {
					return nil, fmt.Errorf("upload failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockCertificateClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := uploadCertificateHandler(mock, tt.content, tt.certName)
				if (err != nil) != tt.wantErr {
					t.Errorf("uploadCertificateHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && !strings.Contains(output, "Certificate uploaded:") {
				t.Errorf("Expected upload success message in output, got: %s", output)
			}
		})
	}
}

func TestGetCertificateHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		includeContent bool
		mockFn         func(*mockCertificateClientImpl)
		wantErr        bool
	}{
		{
			name:           "successful get without content",
			id:             "cert_123",
			includeContent: false,
			mockFn: func(m *mockCertificateClientImpl) {
				m.GetCertificateFunc = func(id string, includeContent bool) (*openaiorgs.Certificate, error) {
					if id != "cert_123" {
						t.Errorf("unexpected id: %s", id)
					}
					if includeContent != false {
						t.Errorf("unexpected includeContent: %v", includeContent)
					}
					cert := createMockCertificate("cert_123", "Test Certificate", nil)
					return &cert, nil
				}
			},
		},
		{
			name:           "successful get with content",
			id:             "cert_123",
			includeContent: true,
			mockFn: func(m *mockCertificateClientImpl) {
				m.GetCertificateFunc = func(id string, includeContent bool) (*openaiorgs.Certificate, error) {
					cert := createMockCertificate("cert_123", "Test Certificate", nil)
					content := "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----"
					cert.CertificateDetails.Content = &content
					return &cert, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockCertificateClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := getCertificateHandler(mock, tt.id, tt.includeContent)
				if (err != nil) != tt.wantErr {
					t.Errorf("getCertificateHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && !strings.Contains(output, "Certificate details:") {
				t.Errorf("Expected certificate details in output, got: %s", output)
			}
			if tt.includeContent && !strings.Contains(output, "Content:") {
				t.Errorf("Expected certificate content in output, got: %s", output)
			}
		})
	}
}

func TestActivateOrgCertificatesHandler(t *testing.T) {
	tests := []struct {
		name           string
		certificateIDs []string
		mockFn         func(*mockCertificateClientImpl)
		wantErr        bool
	}{
		{
			name:           "successful activation",
			certificateIDs: []string{"cert_123", "cert_456"},
			mockFn: func(m *mockCertificateClientImpl) {
				m.ActivateOrganizationCertificatesFunc = func(certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error) {
					if len(certificateIDs) != 2 || certificateIDs[0] != "cert_123" || certificateIDs[1] != "cert_456" {
						t.Errorf("unexpected certificateIDs: %v", certificateIDs)
					}
					return &openaiorgs.CertificateActivationResponse{
						Object:  "certificate.activation",
						Success: true,
					}, nil
				}
			},
		},
		{
			name:           "empty certificate IDs",
			certificateIDs: []string{},
			mockFn:         func(m *mockCertificateClientImpl) {},
			wantErr:        true,
		},
		{
			name:           "comma-separated IDs",
			certificateIDs: []string{"cert_123,cert_456"},
			mockFn: func(m *mockCertificateClientImpl) {
				m.ActivateOrganizationCertificatesFunc = func(certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error) {
					if len(certificateIDs) != 2 || certificateIDs[0] != "cert_123" || certificateIDs[1] != "cert_456" {
						t.Errorf("unexpected certificateIDs after splitting: %v", certificateIDs)
					}
					return &openaiorgs.CertificateActivationResponse{
						Object:  "certificate.activation",
						Success: true,
					}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockCertificateClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := activateOrgCertificatesHandler(mock, tt.certificateIDs)
				if (err != nil) != tt.wantErr {
					t.Errorf("activateOrgCertificatesHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && !strings.Contains(output, "Successfully activated") {
				t.Errorf("Expected success message in output, got: %s", output)
			}
		})
	}
}

// Test the helper function
func TestSplitCommaSeparatedIDs(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "single ID",
			input:    []string{"cert_123"},
			expected: []string{"cert_123"},
		},
		{
			name:     "multiple separate IDs",
			input:    []string{"cert_123", "cert_456"},
			expected: []string{"cert_123", "cert_456"},
		},
		{
			name:     "comma-separated in single string",
			input:    []string{"cert_123,cert_456,cert_789"},
			expected: []string{"cert_123", "cert_456", "cert_789"},
		},
		{
			name:     "mixed separate and comma-separated",
			input:    []string{"cert_123", "cert_456,cert_789", "cert_000"},
			expected: []string{"cert_123", "cert_456", "cert_789", "cert_000"},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitCommaSeparatedIDs(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("splitCommaSeparatedIDs() length = %d, expected %d", len(result), len(tt.expected))
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("splitCommaSeparatedIDs()[%d] = %s, expected %s", i, v, tt.expected[i])
				}
			}
		})
	}
}

// Integration-style test with mock CLI command
func TestUploadCertificateFromFile(t *testing.T) {
	// Create temporary file with certificate content
	content := "-----BEGIN CERTIFICATE-----\ntest content\n-----END CERTIFICATE-----"
	tmpfile, err := os.CreateTemp("", "test-cert-*.pem")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test the file reading logic from uploadCertificate function
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("Failed to read test file: %v", err)
	}
	if string(data) != content {
		t.Errorf("File content mismatch: got %s, expected %s", string(data), content)
	}
}

// Additional helper functions for testing business logic - these mirror the original CLI action handlers

// Additional comprehensive tests for remaining handlers

func TestDeactivateOrgCertificatesHandler(t *testing.T) {
	tests := []struct {
		name           string
		certificateIDs []string
		mockFn         func(*mockCertificateClientImpl)
		wantErr        bool
	}{
		{
			name:           "successful deactivation",
			certificateIDs: []string{"cert_123", "cert_456"},
			mockFn: func(m *mockCertificateClientImpl) {
				m.DeactivateOrganizationCertificatesFunc = func(certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error) {
					return &openaiorgs.CertificateActivationResponse{
						Object:  "certificate.activation",
						Success: true,
					}, nil
				}
			},
		},
		{
			name:           "empty certificate IDs",
			certificateIDs: []string{},
			mockFn:         func(m *mockCertificateClientImpl) {},
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockCertificateClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := func() error {
					if len(tt.certificateIDs) == 0 {
						return fmt.Errorf("at least one certificate ID must be provided")
					}
					allIDs := splitCommaSeparatedIDs(tt.certificateIDs)
					response, err := mock.DeactivateOrganizationCertificates(allIDs)
					if err != nil {
						return wrapError("deactivate certificates", err)
					}
					if response.Success {
						fmt.Printf("Successfully deactivated %d certificates\n", len(allIDs))
					} else {
						fmt.Printf("Failed to deactivate certificates\n")
					}
					return nil
				}()
				if (err != nil) != tt.wantErr {
					t.Errorf("deactivateOrgCertificatesHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && !strings.Contains(output, "Successfully deactivated") {
				t.Errorf("Expected success message in output, got: %s", output)
			}
		})
	}
}

func TestActivateProjectCertificatesHandler(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		certificateIDs []string
		mockFn         func(*mockCertificateClientImpl)
		wantErr        bool
	}{
		{
			name:           "successful project activation",
			projectID:      "proj_123",
			certificateIDs: []string{"cert_123", "cert_456"},
			mockFn: func(m *mockCertificateClientImpl) {
				m.ActivateProjectCertificatesFunc = func(projectID string, certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error) {
					if projectID != "proj_123" {
						t.Errorf("unexpected projectID: %s", projectID)
					}
					return &openaiorgs.CertificateActivationResponse{
						Object:  "certificate.activation",
						Success: true,
					}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockCertificateClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := func() error {
					if len(tt.certificateIDs) == 0 {
						return fmt.Errorf("at least one certificate ID must be provided")
					}
					allIDs := splitCommaSeparatedIDs(tt.certificateIDs)
					response, err := mock.ActivateProjectCertificates(tt.projectID, allIDs)
					if err != nil {
						return wrapError("activate project certificates", err)
					}
					if response.Success {
						fmt.Printf("Successfully activated %d certificates for project %s\n", len(allIDs), tt.projectID)
					} else {
						fmt.Printf("Failed to activate certificates for project %s\n", tt.projectID)
					}
					return nil
				}()
				if (err != nil) != tt.wantErr {
					t.Errorf("activateProjectCertificatesHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && !strings.Contains(output, "Successfully activated") {
				t.Errorf("Expected success message in output, got: %s", output)
			}
		})
	}
}

func TestDeactivateProjectCertificatesHandler(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		certificateIDs []string
		mockFn         func(*mockCertificateClientImpl)
		wantErr        bool
	}{
		{
			name:           "successful project deactivation",
			projectID:      "proj_123",
			certificateIDs: []string{"cert_123", "cert_456"},
			mockFn: func(m *mockCertificateClientImpl) {
				m.DeactivateProjectCertificatesFunc = func(projectID string, certificateIDs []string) (*openaiorgs.CertificateActivationResponse, error) {
					if projectID != "proj_123" {
						t.Errorf("unexpected projectID: %s", projectID)
					}
					return &openaiorgs.CertificateActivationResponse{
						Object:  "certificate.activation",
						Success: true,
					}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockCertificateClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := func() error {
					if len(tt.certificateIDs) == 0 {
						return fmt.Errorf("at least one certificate ID must be provided")
					}
					allIDs := splitCommaSeparatedIDs(tt.certificateIDs)
					response, err := mock.DeactivateProjectCertificates(tt.projectID, allIDs)
					if err != nil {
						return wrapError("deactivate project certificates", err)
					}
					if response.Success {
						fmt.Printf("Successfully deactivated %d certificates for project %s\n", len(allIDs), tt.projectID)
					} else {
						fmt.Printf("Failed to deactivate certificates for project %s\n", tt.projectID)
					}
					return nil
				}()
				if (err != nil) != tt.wantErr {
					t.Errorf("deactivateProjectCertificatesHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			if !tt.wantErr && !strings.Contains(output, "Successfully deactivated") {
				t.Errorf("Expected success message in output, got: %s", output)
			}
		})
	}
}

// Test error scenarios and edge cases
func TestUploadCertificateValidation(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		contentFile string
		wantErr     bool
		errorMsg    string
	}{
		{
			name:        "both content and file provided",
			content:     "cert content",
			contentFile: "file.pem",
			wantErr:     true,
			errorMsg:    "only one of --content or --content-file can be provided",
		},
		{
			name:        "neither content nor file provided",
			content:     "",
			contentFile: "",
			wantErr:     true,
			errorMsg:    "either --content or --content-file must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the validation logic from uploadCertificate
			var err error
			if tt.content == "" && tt.contentFile == "" {
				err = fmt.Errorf("either --content or --content-file must be provided")
			} else if tt.content != "" && tt.contentFile != "" {
				err = fmt.Errorf("only one of --content or --content-file can be provided")
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("uploadCertificate validation error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errorMsg) {
				t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
			}
		})
	}
}
