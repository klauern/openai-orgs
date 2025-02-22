package openaiorgs

import (
	"net/http"
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
	h.mockResponse("GET", "/organizations", http.StatusInternalServerError, map[string]string{
		"error": "internal server error",
	})

	_, err := h.client.ListOrganizations(10, "")
	if err == nil {
		t.Error("Expected error from ListOrganizations, got none")
	}

	// Test get error
	h.mockResponse("GET", "/organizations/org_123", http.StatusNotFound, map[string]string{
		"error": "not found",
	})

	_, err = h.client.GetOrganization("org_123")
	if err == nil {
		t.Error("Expected error from GetOrganization, got none")
	}
}
