package openaiorgs

import "fmt"

type ProjectUser struct {
	Object  string     `json:"object"`
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Email   string     `json:"email"`
	Role    string     `json:"role"`
	AddedAt CustomTime `json:"added_at"`
}

const ProjectUsersListEndpoint = "/organization/projects/%s/users"

func (c *Client) ListProjectUsers(projectID string, limit int, after string) (*ListResponse[ProjectUser], error) {
	queryParams := make(map[string]string)
	if limit > 0 {
		queryParams["limit"] = fmt.Sprintf("%d", limit)
	}
	if after != "" {
		queryParams["after"] = after
	}

	return Get[ProjectUser](c.client, fmt.Sprintf(ProjectUsersListEndpoint, projectID), queryParams)
}

func (c *Client) CreateProjectUser(projectID string, userID string) (*ProjectUser, error) {
	body := map[string]string{"user_id": userID}
	return Post[ProjectUser](c.client, fmt.Sprintf(ProjectUsersListEndpoint, projectID), body)
}

func (c *Client) RetrieveProjectUser(projectID string, userID string) (*ProjectUser, error) {
	return GetSingle[ProjectUser](c.client, fmt.Sprintf(ProjectUsersListEndpoint+"/%s", projectID, userID))
}

func (c *Client) ModifyProjectUser(projectID string, userID string, role RoleType) (*ProjectUser, error) {
	body := map[string]string{"role": string(role)}
	return Post[ProjectUser](c.client, fmt.Sprintf(ProjectUsersListEndpoint+"/%s", projectID, userID), body)
}

func (c *Client) DeleteProjectUser(projectID string, userID string) error {
	return Delete[ProjectUser](c.client, fmt.Sprintf(ProjectUsersListEndpoint+"/%s", projectID, userID))
}
