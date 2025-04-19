package openaiorgs

import "fmt"

type Project struct {
	Object     string       `json:"object"`
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	CreatedAt  UnixSeconds  `json:"created_at"`
	ArchivedAt *UnixSeconds `json:"archived_at,omitempty"`
	Status     string       `json:"status"`
}

const ProjectsListEndpoint = "/organization/projects"

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

func (c *Client) CreateProject(name string) (*Project, error) {
	body := map[string]string{"name": name}
	return Post[Project](c.client, ProjectsListEndpoint, body)
}

func (c *Client) RetrieveProject(id string) (*Project, error) {
	return GetSingle[Project](c.client, ProjectsListEndpoint+"/"+id)
}

func (c *Client) ModifyProject(id string, name string) (*Project, error) {
	body := map[string]string{"name": name}
	return Post[Project](c.client, ProjectsListEndpoint+"/"+id, body)
}

func (c *Client) ArchiveProject(id string) (*Project, error) {
	return Post[Project](c.client, ProjectsListEndpoint+"/"+id+"/archive", nil)
}

// String returns a human-readable string representation of the Project
func (p *Project) String() string {
	status := p.Status
	if p.ArchivedAt != nil {
		status = "archived"
	}
	return fmt.Sprintf("Project{ID: %s, Name: %s, Status: %s}", p.ID, p.Name, status)
}
