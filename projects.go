package openaiorgs

import "fmt"

// Package openaiorgs provides functionality for managing OpenAI organization resources.

// Project represents a project within an OpenAI organization.
// Projects are containers for resources, members, billing limits and settings that can be managed independently.
// Each project has its own API keys, rate limits, and user permissions, allowing for isolated management
// of different workloads or teams within the same organization.
type Project struct {
	// Object identifies the type of this resource.
	// This will always be "project" for Project objects.
	Object string `json:"object"`

	// ID is the unique identifier for this project.
	// This is assigned by OpenAI and cannot be changed.
	ID string `json:"id"`

	// Name is the human-readable identifier for the project.
	// This can be updated using ModifyProject.
	Name string `json:"name"`

	// CreatedAt is the timestamp when this project was created.
	// The timestamp is in Unix epoch seconds.
	CreatedAt UnixSeconds `json:"created_at"`

	// ArchivedAt is the timestamp when this project was archived, if applicable.
	// The timestamp is in Unix epoch seconds. Will be nil for active projects.
	ArchivedAt *UnixSeconds `json:"archived_at,omitempty"`

	// Status indicates the current state of the project.
	// Possible values are "active" or "archived".
	Status string `json:"status"`
}

// ProjectsListEndpoint is the base endpoint for project management operations.
// All project-related API requests are made relative to this path.
const ProjectsListEndpoint = "/organization/projects"

// ListProjects retrieves a paginated list of projects in the organization.
// The results can be filtered to include or exclude archived projects.
// Results are ordered by creation date, with newest projects first.
//
// Parameters:
//   - limit: Maximum number of projects to return (0 for default, which is typically 20)
//   - after: Pagination token for fetching next page (empty string for first page)
//   - includeArchived: Whether to include archived projects in the response
//
// Returns a ListResponse containing the projects and pagination metadata.
// The ListResponse includes the next pagination token if more results are available.
func (c *Client) ListProjects(limit int, after string, includeArchived bool) (*ListResponse[Project], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}
	if includeArchived {
		queryParams["include_archived"] = "true"
	}

	return Get[Project](c.client, ProjectsListEndpoint, queryParams)
}

// CreateProject creates a new project in the organization.
// New projects are created in an active state and can be used immediately.
// The project name must be unique within the organization.
//
// Parameters:
//   - name: The human-readable name for the new project (must be unique)
//
// Returns the created Project object or an error if creation fails.
// Common errors include duplicate project names or reaching project limits.
func (c *Client) CreateProject(name string) (*Project, error) {
	body := map[string]string{"name": name}
	return Post[Project](c.client, ProjectsListEndpoint, body)
}

// RetrieveProject fetches details of a specific project.
// This can be used to get the current state of any project, including archived ones.
//
// Parameters:
//   - id: The unique identifier of the project to retrieve
//
// Returns the Project details or an error if retrieval fails.
// Returns an error if the project ID does not exist or if the caller lacks permission.
func (c *Client) RetrieveProject(id string) (*Project, error) {
	return GetSingle[Project](c.client, ProjectsListEndpoint+"/"+id)
}

// ModifyProject updates the properties of an existing project.
// Currently, only the project name can be modified.
// This operation cannot be performed on archived projects.
//
// Parameters:
//   - id: The unique identifier of the project to modify
//   - name: The new name for the project (must be unique within the organization)
//
// Returns the updated Project object or an error if modification fails.
// Common errors include duplicate names or attempting to modify an archived project.
func (c *Client) ModifyProject(id string, name string) (*Project, error) {
	body := map[string]string{"name": name}
	return Post[Project](c.client, ProjectsListEndpoint+"/"+id, body)
}

// ArchiveProject moves a project to an archived state.
// Archived projects retain their data but cannot be modified or used for new operations.
// This operation cannot be undone through the API.
//
// Parameters:
//   - id: The unique identifier of the project to archive
//
// Returns the updated Project object or an error if archiving fails.
// Returns an error if the project is already archived or if the caller lacks permission.
func (c *Client) ArchiveProject(id string) (*Project, error) {
	return Post[Project](c.client, ProjectsListEndpoint+"/"+id+"/archive", nil)
}

// String returns a human-readable string representation of the Project.
// It includes the project's ID, name, and current status (including archived state).
// This is useful for logging and debugging purposes.
func (p *Project) String() string {
	status := p.Status
	if p.ArchivedAt != nil {
		status = "archived"
	}
	return fmt.Sprintf("Project{ID: %s, Name: %s, Status: %s}", p.ID, p.Name, status)
}
