package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Mock client interface for testing projects
type mockProjectClient interface {
	ListProjects(limit int, after string, includeArchived bool) (*openaiorgs.ListResponse[openaiorgs.Project], error)
	CreateProject(name string) (*openaiorgs.Project, error)
	RetrieveProject(id string) (*openaiorgs.Project, error)
	ModifyProject(id, name string) (*openaiorgs.Project, error)
	ArchiveProject(id string) (*openaiorgs.Project, error)
}

// Mock implementation
type mockProjectClientImpl struct {
	ListProjectsFunc    func(limit int, after string, includeArchived bool) (*openaiorgs.ListResponse[openaiorgs.Project], error)
	CreateProjectFunc   func(name string) (*openaiorgs.Project, error)
	RetrieveProjectFunc func(id string) (*openaiorgs.Project, error)
	ModifyProjectFunc   func(id, name string) (*openaiorgs.Project, error)
	ArchiveProjectFunc  func(id string) (*openaiorgs.Project, error)
}

func (m *mockProjectClientImpl) ListProjects(limit int, after string, includeArchived bool) (*openaiorgs.ListResponse[openaiorgs.Project], error) {
	if m.ListProjectsFunc != nil {
		return m.ListProjectsFunc(limit, after, includeArchived)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectClientImpl) CreateProject(name string) (*openaiorgs.Project, error) {
	if m.CreateProjectFunc != nil {
		return m.CreateProjectFunc(name)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectClientImpl) RetrieveProject(id string) (*openaiorgs.Project, error) {
	if m.RetrieveProjectFunc != nil {
		return m.RetrieveProjectFunc(id)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectClientImpl) ModifyProject(id, name string) (*openaiorgs.Project, error) {
	if m.ModifyProjectFunc != nil {
		return m.ModifyProjectFunc(id, name)
	}
	return nil, fmt.Errorf("not implemented")
}

func (m *mockProjectClientImpl) ArchiveProject(id string) (*openaiorgs.Project, error) {
	if m.ArchiveProjectFunc != nil {
		return m.ArchiveProjectFunc(id)
	}
	return nil, fmt.Errorf("not implemented")
}

// Testable handlers

func listProjectsHandler(client mockProjectClient, limit int, after string, includeArchived bool) error {
	projects, err := client.ListProjects(limit, after, includeArchived)
	if err != nil {
		return wrapError("list projects", err)
	}

	data := TableData{
		Headers: []string{"ID", "Name", "Created At", "Archived At", "Status"},
		Rows:    make([][]string, len(projects.Data)),
	}

	for i, project := range projects.Data {
		archivedAt := "N/A"
		if project.ArchivedAt != nil {
			archivedAt = project.ArchivedAt.String()
		}
		data.Rows[i] = []string{
			project.ID,
			project.Name,
			project.CreatedAt.String(),
			archivedAt,
			project.Status,
		}
	}

	printTableData(data)
	return nil
}

func createProjectHandler(client mockProjectClient, name string) error {
	project, err := client.CreateProject(name)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	fmt.Printf("Project created:\n")
	fmt.Printf("ID: %s\nName: %s\nCreated At: %s\nStatus: %s\n",
		project.ID,
		project.Name,
		project.CreatedAt.String(),
		project.Status,
	)
	return nil
}

func retrieveProjectHandler(client mockProjectClient, id string) error {
	project, err := client.RetrieveProject(id)
	if err != nil {
		return fmt.Errorf("failed to retrieve project: %w", err)
	}

	fmt.Printf("Project details:\n")
	fmt.Printf("ID: %s\nName: %s\nCreated At: %s\nStatus: %s\n",
		project.ID,
		project.Name,
		project.CreatedAt.String(),
		project.Status,
	)
	if project.ArchivedAt != nil {
		fmt.Printf("Archived At: %s\n", project.ArchivedAt.String())
	}
	return nil
}

func modifyProjectHandler(client mockProjectClient, id, name string) error {
	project, err := client.ModifyProject(id, name)
	if err != nil {
		return fmt.Errorf("failed to modify project: %w", err)
	}

	fmt.Printf("Project modified:\n")
	fmt.Printf("ID: %s\nNew Name: %s\nCreated At: %s\nStatus: %s\n",
		project.ID,
		project.Name,
		project.CreatedAt.String(),
		project.Status,
	)
	return nil
}

func archiveProjectHandler(client mockProjectClient, id string) error {
	project, err := client.ArchiveProject(id)
	if err != nil {
		return fmt.Errorf("failed to archive project: %w", err)
	}

	fmt.Printf("Project archived:\n")
	fmt.Printf("ID: %s\nName: %s\nCreated At: %s\nArchived At: %s\nStatus: %s\n",
		project.ID,
		project.Name,
		project.CreatedAt.String(),
		project.ArchivedAt.String(),
		project.Status,
	)
	return nil
}

// Helper to create mock project
func createMockProject(id, name, status string, archived bool) openaiorgs.Project {
	now := openaiorgs.UnixSeconds(time.Now())
	p := openaiorgs.Project{
		Object:    "organization.project",
		ID:        id,
		Name:      name,
		CreatedAt: now,
		Status:    status,
	}
	if archived {
		archivedAt := openaiorgs.UnixSeconds(time.Now())
		p.ArchivedAt = &archivedAt
	}
	return p
}

// Tests

func TestListProjectsHandler(t *testing.T) {
	tests := []struct {
		name            string
		limit           int
		after           string
		includeArchived bool
		mockFn          func(*mockProjectClientImpl)
		wantErr         bool
		wantContains    []string
	}{
		{
			name:            "successful list with active project",
			limit:           10,
			after:           "",
			includeArchived: false,
			mockFn: func(m *mockProjectClientImpl) {
				proj := createMockProject("proj_123", "My Project", "active", false)
				m.ListProjectsFunc = func(limit int, after string, includeArchived bool) (*openaiorgs.ListResponse[openaiorgs.Project], error) {
					return &openaiorgs.ListResponse[openaiorgs.Project]{
						Object:  "list",
						Data:    []openaiorgs.Project{proj},
						FirstID: "proj_123",
						LastID:  "proj_123",
						HasMore: false,
					}, nil
				}
			},
			wantContains: []string{"ID | Name | Created At | Archived At | Status", "proj_123", "My Project", "N/A", "active"},
		},
		{
			name:            "successful list with archived project",
			limit:           10,
			after:           "",
			includeArchived: true,
			mockFn: func(m *mockProjectClientImpl) {
				proj := createMockProject("proj_456", "Archived Project", "archived", true)
				m.ListProjectsFunc = func(limit int, after string, includeArchived bool) (*openaiorgs.ListResponse[openaiorgs.Project], error) {
					return &openaiorgs.ListResponse[openaiorgs.Project]{
						Object:  "list",
						Data:    []openaiorgs.Project{proj},
						FirstID: "proj_456",
						LastID:  "proj_456",
						HasMore: false,
					}, nil
				}
			},
			wantContains: []string{"proj_456", "Archived Project", "archived"},
		},
		{
			name:            "empty list",
			limit:           10,
			after:           "",
			includeArchived: false,
			mockFn: func(m *mockProjectClientImpl) {
				m.ListProjectsFunc = func(limit int, after string, includeArchived bool) (*openaiorgs.ListResponse[openaiorgs.Project], error) {
					return &openaiorgs.ListResponse[openaiorgs.Project]{
						Object: "list",
						Data:   []openaiorgs.Project{},
					}, nil
				}
			},
			wantContains: []string{"ID | Name | Created At | Archived At | Status"},
		},
		{
			name:            "error from client",
			limit:           10,
			after:           "",
			includeArchived: false,
			mockFn: func(m *mockProjectClientImpl) {
				m.ListProjectsFunc = func(limit int, after string, includeArchived bool) (*openaiorgs.ListResponse[openaiorgs.Project], error) {
					return nil, fmt.Errorf("API error")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := listProjectsHandler(mock, tt.limit, tt.after, tt.includeArchived)
				if (err != nil) != tt.wantErr {
					t.Errorf("listProjectsHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestCreateProjectHandler(t *testing.T) {
	tests := []struct {
		name         string
		projectName  string
		mockFn       func(*mockProjectClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:        "successful create",
			projectName: "New Project",
			mockFn: func(m *mockProjectClientImpl) {
				m.CreateProjectFunc = func(name string) (*openaiorgs.Project, error) {
					if name != "New Project" {
						t.Errorf("unexpected name: %s", name)
					}
					proj := createMockProject("proj_new", "New Project", "active", false)
					return &proj, nil
				}
			},
			wantContains: []string{"Project created:", "proj_new", "New Project", "active"},
		},
		{
			name:        "error from client",
			projectName: "Bad Project",
			mockFn: func(m *mockProjectClientImpl) {
				m.CreateProjectFunc = func(name string) (*openaiorgs.Project, error) {
					return nil, fmt.Errorf("creation failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := createProjectHandler(mock, tt.projectName)
				if (err != nil) != tt.wantErr {
					t.Errorf("createProjectHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestRetrieveProjectHandler(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		mockFn       func(*mockProjectClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name: "successful retrieve without archived",
			id:   "proj_123",
			mockFn: func(m *mockProjectClientImpl) {
				m.RetrieveProjectFunc = func(id string) (*openaiorgs.Project, error) {
					if id != "proj_123" {
						t.Errorf("unexpected id: %s", id)
					}
					proj := createMockProject("proj_123", "My Project", "active", false)
					return &proj, nil
				}
			},
			wantContains: []string{"Project details:", "proj_123", "My Project", "active"},
		},
		{
			name: "successful retrieve with archived",
			id:   "proj_456",
			mockFn: func(m *mockProjectClientImpl) {
				m.RetrieveProjectFunc = func(id string) (*openaiorgs.Project, error) {
					proj := createMockProject("proj_456", "Archived Project", "archived", true)
					return &proj, nil
				}
			},
			wantContains: []string{"Project details:", "proj_456", "Archived At:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := retrieveProjectHandler(mock, tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("retrieveProjectHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestModifyProjectHandler(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		newName      string
		mockFn       func(*mockProjectClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name:    "successful modify",
			id:      "proj_123",
			newName: "Updated Name",
			mockFn: func(m *mockProjectClientImpl) {
				m.ModifyProjectFunc = func(id, name string) (*openaiorgs.Project, error) {
					if id != "proj_123" || name != "Updated Name" {
						t.Errorf("unexpected params: id=%s, name=%s", id, name)
					}
					proj := createMockProject("proj_123", "Updated Name", "active", false)
					return &proj, nil
				}
			},
			wantContains: []string{"Project modified:", "proj_123", "Updated Name"},
		},
		{
			name:    "error from client",
			id:      "proj_123",
			newName: "Bad Name",
			mockFn: func(m *mockProjectClientImpl) {
				m.ModifyProjectFunc = func(id, name string) (*openaiorgs.Project, error) {
					return nil, fmt.Errorf("modify failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := modifyProjectHandler(mock, tt.id, tt.newName)
				if (err != nil) != tt.wantErr {
					t.Errorf("modifyProjectHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestArchiveProjectHandler(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		mockFn       func(*mockProjectClientImpl)
		wantErr      bool
		wantContains []string
	}{
		{
			name: "successful archive",
			id:   "proj_123",
			mockFn: func(m *mockProjectClientImpl) {
				m.ArchiveProjectFunc = func(id string) (*openaiorgs.Project, error) {
					if id != "proj_123" {
						t.Errorf("unexpected id: %s", id)
					}
					proj := createMockProject("proj_123", "My Project", "archived", true)
					return &proj, nil
				}
			},
			wantContains: []string{"Project archived:", "proj_123", "My Project", "archived"},
		},
		{
			name: "error from client",
			id:   "proj_123",
			mockFn: func(m *mockProjectClientImpl) {
				m.ArchiveProjectFunc = func(id string) (*openaiorgs.Project, error) {
					return nil, fmt.Errorf("archive failed")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectClientImpl{}
			tt.mockFn(mock)

			output := captureOutput(func() {
				err := archiveProjectHandler(mock, tt.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("archiveProjectHandler() error = %v, wantErr %v", err, tt.wantErr)
				}
			})

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}
