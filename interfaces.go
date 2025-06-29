package openaiorgs

// OpenAIOrgsClient defines the interface for interacting with the OpenAI Organizations API.
type OpenAIOrgsClient interface {
	// Project Management
	ListProjects(limit int, after string, includeArchived bool) (*ListResponse[Project], error)
	CreateProject(name string) (*Project, error)
	RetrieveProject(id string) (*Project, error)
	ModifyProject(id string, name string) (*Project, error)
	ArchiveProject(id string) (*Project, error)

	// Project Users
	ListProjectUsers(projectID string, limit int, after string) (*ListResponse[ProjectUser], error)
	CreateProjectUser(projectID string, userID string, role string) (*ProjectUser, error)
	RetrieveProjectUser(projectID string, userID string) (*ProjectUser, error)
	ModifyProjectUser(projectID string, userID string, role string) (*ProjectUser, error)
	DeleteProjectUser(projectID string, userID string) error

	// Organization Users
	ListUsers(limit int, after string) (*ListResponse[User], error)
	RetrieveUser(id string) (*User, error)
	DeleteUser(id string) error
	ModifyUserRole(id string, role string) error

	// Organization Invites
	ListInvites(limit int, after string) (*ListResponse[Invite], error)
	CreateInvite(email string, role string) (*Invite, error)
	RetrieveInvite(id string) (*Invite, error)
	DeleteInvite(id string) error

	// Project API Keys
	ListProjectApiKeys(projectID string, limit int, after string) (*ListResponse[ProjectApiKey], error)
	RetrieveProjectApiKey(projectID string, apiKeyID string) (*ProjectApiKey, error)
	DeleteProjectApiKey(projectID string, apiKeyID string) error

	// Admin API Keys
	ListAdminAPIKeys(limit int, after string) (*ListResponse[AdminAPIKey], error)
	CreateAdminAPIKey(name string, scopes []string) (*AdminAPIKey, error)
	RetrieveAdminAPIKey(apiKeyID string) (*AdminAPIKey, error)
	DeleteAdminAPIKey(apiKeyID string) error

	// Project Service Accounts
	ListProjectServiceAccounts(projectID string, limit int, after string) (*ListResponse[ProjectServiceAccount], error)
	CreateProjectServiceAccount(projectID string, name string) (*ProjectServiceAccount, error)
	RetrieveProjectServiceAccount(projectID string, serviceAccountID string) (*ProjectServiceAccount, error)
	DeleteProjectServiceAccount(projectID string, serviceAccountID string) error

	// Project Rate Limits
	ListProjectRateLimits(limit int, after string, projectId string) (*ListResponse[ProjectRateLimit], error)

	// Audit Logs
	ListAuditLogs(params *AuditLogListParams) (*ListResponse[AuditLog], error)

	// Organization Certificates
	ListOrganizationCertificates(limit int, after string, order string) (*ListResponse[Certificate], error)
	UploadCertificate(content string, name string) (*Certificate, error)
	GetCertificate(certificateID string, includeContent bool) (*Certificate, error)
	ModifyCertificate(certificateID string, name string) (*Certificate, error)
	DeleteCertificate(certificateID string) (*CertificateDeletedResponse, error)
	ActivateOrganizationCertificates(certificateIDs []string) (*CertificateActivationResponse, error)
	DeactivateOrganizationCertificates(certificateIDs []string) (*CertificateActivationResponse, error)

	// Project Certificates
	ListProjectCertificates(projectID string, limit int, after string, order string) (*ListResponse[Certificate], error)
	ActivateProjectCertificates(projectID string, certificateIDs []string) (*CertificateActivationResponse, error)
	DeactivateProjectCertificates(projectID string, certificateIDs []string) (*CertificateActivationResponse, error)
}
