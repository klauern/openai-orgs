package cmd

import (
	"strings"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
)

// Helper to create mock project data
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

func TestListProjectsCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful list with active project",
			args:       []string{"projects", "list"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.Project]{
				Object:  "list",
				Data:    []openaiorgs.Project{createMockProject("proj_123", "My Project", "active", false)},
				FirstID: "proj_123",
				LastID:  "proj_123",
				HasMore: false,
			},
			wantContains: []string{"ID | Name | Created At | Archived At | Status", "proj_123", "My Project", "N/A", "active"},
		},
		{
			name:       "successful list with archived project",
			args:       []string{"projects", "list"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.Project]{
				Object:  "list",
				Data:    []openaiorgs.Project{createMockProject("proj_456", "Archived Project", "archived", true)},
				FirstID: "proj_456",
				LastID:  "proj_456",
				HasMore: false,
			},
			wantContains: []string{"proj_456", "Archived Project", "archived"},
		},
		{
			name:       "empty list",
			args:       []string{"projects", "list"},
			statusCode: 200,
			response: openaiorgs.ListResponse[openaiorgs.Project]{
				Object: "list",
				Data:   []openaiorgs.Project{},
			},
			wantContains: []string{"ID | Name | Created At | Archived At | Status"},
		},
		{
			name:       "error from API",
			args:       []string{"projects", "list"},
			statusCode: 500,
			response:   map[string]string{"error": "internal server error"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("GET", "/organization/projects", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectsCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}

			if !tt.wantErr {
				h.assertRequest("GET", "/organization/projects", 1)
			}
		})
	}
}

func TestCreateProjectCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful create",
			args:       []string{"projects", "create", "--name", "New Project"},
			statusCode: 200,
			response:   createMockProject("proj_new", "New Project", "active", false),
			wantContains: []string{"Project created:", "proj_new", "New Project", "active"},
		},
		{
			name:       "error from API",
			args:       []string{"projects", "create", "--name", "Bad Project"},
			statusCode: 500,
			response:   map[string]string{"error": "creation failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("POST", "/organization/projects", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectsCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestRetrieveProjectCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful retrieve without archived",
			args:       []string{"projects", "retrieve", "--id", "proj_123"},
			statusCode: 200,
			response:   createMockProject("proj_123", "My Project", "active", false),
			wantContains: []string{"Project details:", "proj_123", "My Project", "active"},
		},
		{
			name:       "successful retrieve with archived",
			args:       []string{"projects", "retrieve", "--id", "proj_456"},
			statusCode: 200,
			response:   createMockProject("proj_456", "Archived Project", "archived", true),
			wantContains: []string{"Project details:", "proj_456", "Archived At:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			// Mock the specific project endpoint
			projectID := tt.args[3] // --id value
			h.mockResponse("GET", "/organization/projects/"+projectID, tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectsCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestModifyProjectCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful modify",
			args:       []string{"projects", "modify", "--id", "proj_123", "--name", "Updated Name"},
			statusCode: 200,
			response:   createMockProject("proj_123", "Updated Name", "active", false),
			wantContains: []string{"Project modified:", "proj_123", "Updated Name"},
		},
		{
			name:       "error from API",
			args:       []string{"projects", "modify", "--id", "proj_123", "--name", "Bad Name"},
			statusCode: 500,
			response:   map[string]string{"error": "modify failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("POST", "/organization/projects/proj_123", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectsCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}

func TestArchiveProjectCommand(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		statusCode   int
		response     any
		wantErr      bool
		wantContains []string
	}{
		{
			name:       "successful archive",
			args:       []string{"projects", "archive", "--id", "proj_123"},
			statusCode: 200,
			response:   createMockProject("proj_123", "My Project", "archived", true),
			wantContains: []string{"Project archived:", "proj_123", "My Project", "archived"},
		},
		{
			name:       "error from API",
			args:       []string{"projects", "archive", "--id", "proj_123"},
			statusCode: 500,
			response:   map[string]string{"error": "archive failed"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newCmdTestHelper(t)
			defer h.cleanup()

			h.mockResponse("POST", "/organization/projects/proj_123/archive", tt.statusCode, tt.response)

			var output string
			var runErr error
			output = captureOutput(func() {
				runErr = h.runCmd(ProjectsCommand(), tt.args)
			})

			if (runErr != nil) != tt.wantErr {
				t.Errorf("runCmd() error = %v, wantErr %v", runErr, tt.wantErr)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Expected output to contain %q, got: %s", want, output)
				}
			}
		})
	}
}
