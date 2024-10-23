package oaiprom

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

const AuditLogsListEndpoint = "/organization/audit_logs"

// AuditLog represents the main audit log object
type AuditLog struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Actor     Actor     `json:"actor"`
	Event     Event     `json:"event"`
}

// Actor represents the entity performing the action
type Actor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Event represents the details of the audit log event
type Event struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Action  string      `json:"action"`
	Auth    Auth        `json:"auth"`
	Payload interface{} `json:"payload"`
}

// Auth represents authentication details
type Auth struct {
	Type      string `json:"type"`
	Transport string `json:"transport"`
}

// AccessPolicyCreated represents the payload for access policy creation
type AccessPolicyCreated struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AccessPolicyDeleted represents the payload for access policy deletion
type AccessPolicyDeleted struct {
	ID string `json:"id"`
}

// AccessPolicyUpdated represents the payload for access policy updates
type AccessPolicyUpdated struct {
	ID      string `json:"id"`
	Changes struct {
		Name struct {
			Old string `json:"old"`
			New string `json:"new"`
		} `json:"name"`
	} `json:"changes"`
}

// APIKeyCreated represents the payload for API key creation
type APIKeyCreated struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// APIKeyDeleted represents the payload for API key deletion
type APIKeyDeleted struct {
	ID string `json:"id"`
}

// AssistantCreated represents the payload for assistant creation
type AssistantCreated struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AssistantDeleted represents the payload for assistant deletion
type AssistantDeleted struct {
	ID string `json:"id"`
}

// AssistantModified represents the payload for assistant modifications
type AssistantModified struct {
	ID      string `json:"id"`
	Changes struct {
		Name struct {
			Old string `json:"old"`
			New string `json:"new"`
		} `json:"name"`
	} `json:"changes"`
}

// FileCreated represents the payload for file creation
type FileCreated struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// FileDeleted represents the payload for file deletion
type FileDeleted struct {
	ID string `json:"id"`
}

// FineTuneCreated represents the payload for fine-tune creation
type FineTuneCreated struct {
	ID string `json:"id"`
}

// FineTuneDeleted represents the payload for fine-tune deletion
type FineTuneDeleted struct {
	ID string `json:"id"`
}

// FineTuneEventCreated represents the payload for fine-tune event creation
type FineTuneEventCreated struct {
	ID           string    `json:"id"`
	FineTuneID   string    `json:"fine_tune_id"`
	Level        string    `json:"level"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at"`
	SerializedAt time.Time `json:"serialized_at"`
}

// ModelCreated represents the payload for model creation
type ModelCreated struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ModelDeleted represents the payload for model deletion
type ModelDeleted struct {
	ID string `json:"id"`
}

// RunCreated represents the payload for run creation
type RunCreated struct {
	ID            string    `json:"id"`
	ThreadID      string    `json:"thread_id"`
	AssistantID   string    `json:"assistant_id"`
	Status        string    `json:"status"`
	StartedAt     time.Time `json:"started_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	CancelledAt   time.Time `json:"cancelled_at"`
	FailedAt      time.Time `json:"failed_at"`
	CompletedAt   time.Time `json:"completed_at"`
	LastErrorCode string    `json:"last_error_code"`
}

// RunModified represents the payload for run modifications
type RunModified struct {
	ID      string `json:"id"`
	Changes struct {
		Status struct {
			Old string `json:"old"`
			New string `json:"new"`
		} `json:"status"`
	} `json:"changes"`
}

// ThreadCreated represents the payload for thread creation
type ThreadCreated struct {
	ID string `json:"id"`
}

// ThreadDeleted represents the payload for thread deletion
type ThreadDeleted struct {
	ID string `json:"id"`
}

// ThreadModified represents the payload for thread modifications
type ThreadModified struct {
	ID      string `json:"id"`
	Changes struct {
		Metadata struct {
			Old map[string]interface{} `json:"old"`
			New map[string]interface{} `json:"new"`
		} `json:"metadata"`
	} `json:"changes"`
}

// AuditLogListParams represents the query parameters for listing audit logs
type AuditLogListParams struct {
	Limit     int       `json:"limit,omitempty"`
	After     string    `json:"after,omitempty"`
	Before    string    `json:"before,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

func (c *Client) ListAuditLogs(params *AuditLogListParams) (*ListResponse[AuditLog], error) {
	queryParams := make(map[string]string)

	if params != nil {
		if params.Limit > 0 {
			queryParams["limit"] = fmt.Sprintf("%d", params.Limit)
		}
		if params.After != "" {
			queryParams["after"] = params.After
		}
		if params.Before != "" {
			queryParams["before"] = params.Before
		}
		if !params.StartDate.IsZero() {
			queryParams["start_date"] = params.StartDate.Format(time.RFC3339)
		}
		if !params.EndDate.IsZero() {
			queryParams["end_date"] = params.EndDate.Format(time.RFC3339)
		}
	}

	return Get[AuditLog](c.client, AuditLogsListEndpoint, queryParams)
}

// ParseAuditLogPayload parses the payload of an AuditLog based on its type
func ParseAuditLogPayload(auditLog *AuditLog) (interface{}, error) {
	var payload interface{}

	switch auditLog.Type {
	case "access_policy.created":
		payload = &AccessPolicyCreated{}
	case "access_policy.deleted":
		payload = &AccessPolicyDeleted{}
	case "access_policy.updated":
		payload = &AccessPolicyUpdated{}
	case "api_key.created":
		payload = &APIKeyCreated{}
	case "api_key.deleted":
		payload = &APIKeyDeleted{}
	case "assistant.created":
		payload = &AssistantCreated{}
	case "assistant.deleted":
		payload = &AssistantDeleted{}
	case "assistant.modified":
		payload = &AssistantModified{}
	case "file.created":
		payload = &FileCreated{}
	case "file.deleted":
		payload = &FileDeleted{}
	case "fine_tune.created":
		payload = &FineTuneCreated{}
	case "fine_tune.deleted":
		payload = &FineTuneDeleted{}
	case "fine_tune.event.created":
		payload = &FineTuneEventCreated{}
	case "model.created":
		payload = &ModelCreated{}
	case "model.deleted":
		payload = &ModelDeleted{}
	case "run.created":
		payload = &RunCreated{}
	case "run.modified":
		payload = &RunModified{}
	case "thread.created":
		payload = &ThreadCreated{}
	case "thread.deleted":
		payload = &ThreadDeleted{}
	case "thread.modified":
		payload = &ThreadModified{}
	case "invite.sent":
		payload = &InviteSent{}
	case "login.succeeded":
		payload = &LoginSucceeded{}
	case "logout.succeeded":
		payload = &LogoutSucceeded{}
	case "organization.updated":
		payload = &OrganizationUpdated{}
	default:
		return nil, fmt.Errorf("unknown audit log type: %s", auditLog.Type)
	}

	err := mapstructure.Decode(auditLog.Event.Payload, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode audit log payload: %w", err)
	}

	return payload, nil
}

// Add these new structs for the missing audit log types
type InviteSent struct {
	Email string `json:"email"`
}

type LoginSucceeded struct {
	// Add relevant fields if available in the API response
}

type LogoutSucceeded struct {
	// Add relevant fields if available in the API response
}

type OrganizationUpdated struct {
	Changes struct {
		// Add relevant fields based on what can be updated in an organization
		Name struct {
			Old string `json:"old"`
			New string `json:"new"`
		} `json:"name,omitempty"`
		// Add other fields as necessary
	} `json:"changes"`
}
