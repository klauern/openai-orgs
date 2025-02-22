package openaiorgs

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestListOrganizations(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	expectedOrg := Organization{
		ID:          "org_123",
		Object:      "organization",
		Name:        "Test Org",
		Description: "Test Description",
		CreatedAt:   UnixSeconds(time.Now()),
		IsBlocked:   false,
	}

	response := ListResponse[Organization]{
		Object:  "list",
		Data:    []Organization{expectedOrg},
		HasMore: false,
	}

	h.mockResponse("GET", "/organizations", http.StatusOK, response)

	orgs, err := h.client.ListOrganizations(10, "org_456")
	if err != nil {
		t.Errorf("ListOrganizations returned error: %v", err)
	}

	if len(orgs.Data) != 1 {
		t.Errorf("Expected 1 organization, got %d", len(orgs.Data))
	}

	if orgs.Data[0].ID != expectedOrg.ID {
		t.Errorf("Expected org_123, got %s", orgs.Data[0].ID)
	}

	h.assertRequest("GET", "/organizations", 1)
}

func TestGetOrganization(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	expectedOrg := Organization{
		ID:          "org_123",
		Object:      "organization",
		Name:        "Test Org",
		Description: "Test Description",
		CreatedAt:   UnixSeconds(time.Now()),
		IsBlocked:   false,
	}

	h.mockResponse("GET", "/organizations/org_123", http.StatusOK, expectedOrg)

	org, err := h.client.GetOrganization("org_123")
	if err != nil {
		t.Errorf("GetOrganization returned error: %v", err)
	}

	if org.ID != expectedOrg.ID {
		t.Errorf("Expected org_123, got %s", org.ID)
	}

	h.assertRequest("GET", "/organizations/org_123", 1)
}

func TestUpdateOrganization(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	name := "Updated Org"
	desc := "Updated Description"
	req := &UpdateOrganizationRequest{
		Name:        &name,
		Description: &desc,
	}

	expectedOrg := Organization{
		ID:          "org_123",
		Object:      "organization",
		Name:        name,
		Description: desc,
		CreatedAt:   UnixSeconds(time.Now()),
		IsBlocked:   false,
	}

	h.mockResponse("POST", "/organizations/org_123", http.StatusOK, expectedOrg)

	org, err := h.client.UpdateOrganization("org_123", req)
	if err != nil {
		t.Errorf("UpdateOrganization returned error: %v", err)
	}

	if org.Name != name {
		t.Errorf("Expected name %s, got %s", name, org.Name)
	}

	if org.Description != desc {
		t.Errorf("Expected description %s, got %s", desc, org.Description)
	}

	h.assertRequest("POST", "/organizations/org_123", 1)
}

func TestDeleteOrganization(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	h.mockResponse("DELETE", "/organizations/org_123", http.StatusNoContent, nil)

	err := h.client.DeleteOrganization("org_123")
	if err != nil {
		t.Errorf("DeleteOrganization returned error: %v", err)
	}

	h.assertRequest("DELETE", "/organizations/org_123", 1)
}

func TestOrganizationErrors(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Test list error
	h.mockResponse("GET", "/organizations", http.StatusInternalServerError, map[string]interface{}{
		"error": map[string]interface{}{
			"message": "internal server error",
			"type":    "internal_error",
			"code":    "internal_error",
		},
	})

	_, err := h.client.ListOrganizations(10, "")
	if err == nil {
		t.Error("Expected error from ListOrganizations, got none")
	}
	h.assertRequest("GET", "/organizations", 1)

	// Test get error
	h.mockResponse("GET", "/organizations/org_123", http.StatusNotFound, map[string]interface{}{
		"error": map[string]interface{}{
			"message": "organization not found",
			"type":    "not_found",
			"code":    "not_found",
		},
	})

	_, err = h.client.GetOrganization("org_123")
	if err == nil {
		t.Error("Expected error from GetOrganization, got none")
	}
	if err != nil && !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
	h.assertRequest("GET", "/organizations/org_123", 1)

	// Test update error
	h.mockResponse("POST", "/organizations/org_123", http.StatusBadRequest, map[string]interface{}{
		"error": map[string]interface{}{
			"message": "invalid request",
			"type":    "invalid_request",
			"code":    "invalid_request",
		},
	})

	name := "Test"
	_, err = h.client.UpdateOrganization("org_123", &UpdateOrganizationRequest{
		Name: &name,
	})
	if err == nil {
		t.Error("Expected error from UpdateOrganization, got none")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid request") {
		t.Errorf("Expected 'invalid request' error, got: %v", err)
	}
	h.assertRequest("POST", "/organizations/org_123", 1)

	// Test delete error
	h.mockResponse("DELETE", "/organizations/org_123", http.StatusForbidden, map[string]interface{}{
		"error": map[string]interface{}{
			"message": "permission denied",
			"type":    "permission_denied",
			"code":    "permission_denied",
		},
	})

	err = h.client.DeleteOrganization("org_123")
	if err == nil {
		t.Error("Expected error from DeleteOrganization, got none")
	}
	if err != nil && !strings.Contains(err.Error(), "permission denied") {
		t.Errorf("Expected 'permission denied' error, got: %v", err)
	}
	h.assertRequest("DELETE", "/organizations/org_123", 1)
}
