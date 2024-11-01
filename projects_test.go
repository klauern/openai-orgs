package openaiorgs

import (
	"testing"
	"time"
)

func TestListProjects(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	// Mock response data
	now := time.Now()
	mockProjects := []Project{
		{
			Object:    "project",
			ID:        "proj_123",
			Name:      "Test Project",
			CreatedAt: UnixSeconds(now),
			Status:    "active",
		},
	}

	// Register mock response
	response := ListResponse[Project]{
		Object:  "list",
		Data:    mockProjects,
		FirstID: "proj_123",
		LastID:  "proj_123",
		HasMore: false,
	}
	h.mockResponse("GET", ProjectsListEndpoint, 200, response)

	// Make the API call
	projects, err := h.client.ListProjects(10, "", false)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
		return
	}
	if len(projects.Data) != 1 {
		t.Errorf("Expected 1 project, got %d", len(projects.Data))
		return
	}
	if mockProjects[0].ID != projects.Data[0].ID {
		t.Errorf("Expected ID %s, got %s", mockProjects[0].ID, projects.Data[0].ID)
	}
	if mockProjects[0].Name != projects.Data[0].Name {
		t.Errorf("Expected Name %s, got %s", mockProjects[0].Name, projects.Data[0].Name)
	}

	// Verify the request was made
	h.assertRequest("GET", ProjectsListEndpoint, 1)
}

func TestCreateProject(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	mockProject := Project{
		Object:    "project",
		ID:        "proj_123",
		Name:      "New Project",
		CreatedAt: UnixSeconds(time.Now()),
		Status:    "active",
	}

	h.mockResponse("POST", ProjectsListEndpoint, 200, mockProject)

	// Make the API call
	project, err := h.client.CreateProject("New Project")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if project == nil {
		t.Error("Expected project, got nil")
		return
	}
	if mockProject.ID != project.ID {
		t.Errorf("Expected ID %s, got %s", mockProject.ID, project.ID)
	}
	if mockProject.Name != project.Name {
		t.Errorf("Expected Name %s, got %s", mockProject.Name, project.Name)
	}

	// Verify the request was made
	h.assertRequest("POST", ProjectsListEndpoint, 1)
}

func TestRetrieveProject(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	projectID := "proj_123"
	mockProject := Project{
		Object:    "project",
		ID:        projectID,
		Name:      "Test Project",
		CreatedAt: UnixSeconds(time.Now()),
		Status:    "active",
	}

	h.mockResponse("GET", ProjectsListEndpoint+"/"+projectID, 200, mockProject)

	// Make the API call
	project, err := h.client.RetrieveProject(projectID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if project == nil {
		t.Error("Expected project, got nil")
		return
	}
	if mockProject.ID != project.ID {
		t.Errorf("Expected ID %s, got %s", mockProject.ID, project.ID)
	}

	// Verify the request was made
	h.assertRequest("GET", ProjectsListEndpoint+"/"+projectID, 1)
}

func TestModifyProject(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	projectID := "proj_123"
	mockProject := Project{
		Object:    "project",
		ID:        projectID,
		Name:      "Updated Project",
		CreatedAt: UnixSeconds(time.Now()),
		Status:    "active",
	}

	h.mockResponse("POST", ProjectsListEndpoint+"/"+projectID, 200, mockProject)

	// Make the API call
	project, err := h.client.ModifyProject(projectID, "Updated Project")
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if project == nil {
		t.Error("Expected project, got nil")
		return
	}
	if mockProject.ID != project.ID {
		t.Errorf("Expected ID %s, got %s", mockProject.ID, project.ID)
	}
	if mockProject.Name != project.Name {
		t.Errorf("Expected Name %s, got %s", mockProject.Name, project.Name)
	}

	// Verify the request was made
	h.assertRequest("POST", ProjectsListEndpoint+"/"+projectID, 1)
}

func TestArchiveProject(t *testing.T) {
	h := newTestHelper(t)
	defer h.cleanup()

	projectID := "proj_123"
	mockProject := Project{
		Object:    "project",
		ID:        projectID,
		Name:      "Test Project",
		CreatedAt: UnixSeconds(time.Now()),
		Status:    "archived",
	}

	h.mockResponse("POST", ProjectsListEndpoint+"/"+projectID+"/archive", 200, mockProject)

	// Make the API call
	project, err := h.client.ArchiveProject(projectID)
	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if project == nil {
		t.Error("Expected project, got nil")
		return
	}
	if mockProject.ID != project.ID {
		t.Errorf("Expected ID %s, got %s", mockProject.ID, project.ID)
	}
	if mockProject.Status != project.Status {
		t.Errorf("Expected Status %s, got %s", mockProject.Status, project.Status)
	}

	// Verify the request was made
	h.assertRequest("POST", ProjectsListEndpoint+"/"+projectID+"/archive", 1)
}
