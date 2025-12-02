package openaiorgs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const AuditLogsListEndpoint = "/organization/audit_logs"

// AuditLogListResponse represents the paginated response from the audit logs endpoint
type AuditLogListResponse struct {
	Object  string     `json:"object"`
	Data    []AuditLog `json:"data"`
	FirstID string     `json:"first_id"`
	LastID  string     `json:"last_id"`
	HasMore bool       `json:"has_more"`
}

// AuditLog represents the main audit log object
type AuditLog struct {
	Object      string        `json:"object,omitempty"`
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	EffectiveAt UnixSeconds   `json:"effective_at"`
	Project     *AuditProject `json:"project,omitempty"`
	Actor       Actor         `json:"actor"`
	Details     any           `json:"-"` // This will be unmarshaled based on Type
}

// AuditProject represents project information in audit logs
type AuditProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Actor represents the entity performing the action
type Actor struct {
	Type    string       `json:"type"` // "session" or "api_key"
	Session *Session     `json:"session,omitempty"`
	APIKey  *APIKeyActor `json:"api_key,omitempty"`
}

// Session represents user session information
type Session struct {
	User             AuditUser         `json:"user"`
	IPAddress        string            `json:"ip_address"`
	UserAgent        string            `json:"user_agent"`
	JA3              string            `json:"ja3,omitempty"`
	JA4              string            `json:"ja4,omitempty"`
	IPAddressDetails *IPAddressDetails `json:"ip_address_details,omitempty"`
}

// IPAddressDetails contains geolocation information for an IP address
type IPAddressDetails struct {
	Country    string `json:"country,omitempty"`
	City       string `json:"city,omitempty"`
	Region     string `json:"region,omitempty"`
	RegionCode string `json:"region_code,omitempty"`
	ASN        string `json:"asn,omitempty"`
	Latitude   string `json:"latitude,omitempty"`
	Longitude  string `json:"longitude,omitempty"`
}

// APIKeyActor represents API key information in the actor field
type APIKeyActor struct {
	Type string    `json:"type"`
	User AuditUser `json:"user"`
}

// AuditUser represents user information in audit logs
type AuditUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Event types and their corresponding payload structures
type APIKeyCreated struct {
	ID   string `json:"id"`
	Data struct {
		Scopes []string `json:"scopes"`
	} `json:"data"`
}

type APIKeyUpdated struct {
	ID               string `json:"id"`
	ChangesRequested struct {
		Scopes []string `json:"scopes"`
	} `json:"changes_requested"`
}

type APIKeyDeleted struct {
	ID string `json:"id"`
}

type InviteSent struct {
	ID   string `json:"id"`
	Data struct {
		Email string `json:"email"`
	} `json:"data"`
}

type InviteAccepted struct {
	ID string `json:"id"`
}

type InviteDeleted struct {
	ID string `json:"id"`
}

type LoginFailed struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type LogoutFailed struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type OrganizationUpdated struct {
	ID               string `json:"id"`
	ChangesRequested struct {
		Name string `json:"name,omitempty"`
	} `json:"changes_requested"`
}

// ProjectCreated represents the details for project.created events
type ProjectCreated struct {
	ID   string `json:"id"`
	Data struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	} `json:"data"`
}

// ProjectUpdated represents the details for project.updated events
type ProjectUpdated struct {
	ID               string `json:"id"`
	ChangesRequested struct {
		Title string `json:"title"`
	} `json:"changes_requested"`
}

// ProjectArchived represents the details for project.archived events
type ProjectArchived struct {
	ID string `json:"id"`
}

// RateLimitUpdated represents the details for rate_limit.updated events
type RateLimitUpdated struct {
	ID               string `json:"id"`
	ChangesRequested struct {
		MaxRequestsPer1Minute       int `json:"max_requests_per_1_minute,omitempty"`
		MaxTokensPer1Minute         int `json:"max_tokens_per_1_minute,omitempty"`
		MaxImagesPer1Minute         int `json:"max_images_per_1_minute,omitempty"`
		MaxAudioMegabytesPer1Minute int `json:"max_audio_megabytes_per_1_minute,omitempty"`
		MaxRequestsPer1Day          int `json:"max_requests_per_1_day,omitempty"`
		Batch1DayMaxInputTokens     int `json:"batch_1_day_max_input_tokens,omitempty"`
	} `json:"changes_requested"`
}

// RateLimitDeleted represents the details for rate_limit.deleted events
type RateLimitDeleted struct {
	ID string `json:"id"`
}

// ServiceAccountCreated represents the details for service_account.created events
type ServiceAccountCreated struct {
	ID   string `json:"id"`
	Data struct {
		Role string `json:"role"` // Either "owner" or "member"
	} `json:"data"`
}

// ServiceAccountUpdated represents the details for service_account.updated events
type ServiceAccountUpdated struct {
	ID               string `json:"id"`
	ChangesRequested struct {
		Role string `json:"role"` // Either "owner" or "member"
	} `json:"changes_requested"`
}

// ServiceAccountDeleted represents the details for service_account.deleted events
type ServiceAccountDeleted struct {
	ID string `json:"id"`
}

// UserAdded represents the details for user.added events
type UserAdded struct {
	ID   string `json:"id"`
	Data struct {
		Role string `json:"role"` // Either "owner" or "member"
	} `json:"data"`
}

// UserUpdated represents the details for user.updated events
type UserUpdated struct {
	ID               string `json:"id"`
	ChangesRequested struct {
		Role string `json:"role"` // Either "owner" or "member"
	} `json:"changes_requested"`
}

// UserDeleted represents the details for user.deleted events
type UserDeleted struct {
	ID string `json:"id"`
}

type LoginSucceeded struct {
	Object      string `json:"object"`
	ID          string `json:"id"`
	Type        string `json:"type"`
	EffectiveAt int64  `json:"effective_at"`
	Actor       Actor  `json:"actor"`
}

type EffectiveAt struct {
	Gte int64 `json:"gte,omitempty"`
	Gt  int64 `json:"gt,omitempty"`
	Lte int64 `json:"lte,omitempty"`
	Lt  int64 `json:"lt,omitempty"`
}

// AuditLogListParams represents the query parameters for listing audit logs
type AuditLogListParams struct {
	EffectiveAt *EffectiveAt `json:"effective_at,omitempty"`
	ProjectIDs  []string     `json:"project_ids,omitempty"`
	EventTypes  []string     `json:"event_types,omitempty"`
	ActorIDs    []string     `json:"actor_ids,omitempty"`
	ActorEmails []string     `json:"actor_emails,omitempty"`
	ResourceIDs []string     `json:"resource_ids,omitempty"`
	Limit       int          `json:"limit,omitempty"`
	After       string       `json:"after,omitempty"`
	Before      string       `json:"before,omitempty"`
}

// rawAuditLog is used internally to parse audit log JSON with dynamic event keys
type rawAuditLog struct {
	Object      string        `json:"object,omitempty"`
	ID          string        `json:"id"`
	Type        string        `json:"type"`
	EffectiveAt UnixSeconds   `json:"effective_at"`
	Project     *AuditProject `json:"project,omitempty"`
	Actor       Actor         `json:"actor"`
}

// UnmarshalJSON handles the event-specific details using dynamic keys.
// The OpenAI API returns event details under a key matching the event type
// (e.g., "invite.deleted": {...}) rather than a static "details" key.
func (a *AuditLog) UnmarshalJSON(data []byte) error {
	// First, unmarshal the common fields
	var raw rawAuditLog
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Copy the common fields
	a.Object = raw.Object
	a.ID = raw.ID
	a.Type = raw.Type
	a.EffectiveAt = raw.EffectiveAt
	a.Project = raw.Project
	a.Actor = raw.Actor

	// Now unmarshal into a map to extract the dynamic event key
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// Look for event details under the dynamic key (e.g., "invite.deleted")
	eventData, hasEventData := rawMap[raw.Type]
	if !hasEventData || len(eventData) == 0 || string(eventData) == "null" {
		a.Details = nil
		return nil
	}

	// Parse the details based on the event type
	var details any
	switch raw.Type {
	case "api_key.created":
		details = &APIKeyCreated{}
	case "api_key.updated":
		details = &APIKeyUpdated{}
	case "api_key.deleted":
		details = &APIKeyDeleted{}
	case "invite.sent":
		details = &InviteSent{}
	case "invite.accepted":
		details = &InviteAccepted{}
	case "invite.deleted":
		details = &InviteDeleted{}
	case "login.failed":
		details = &LoginFailed{}
	case "login.succeeded":
		details = &LoginSucceeded{}
	case "logout.failed":
		details = &LogoutFailed{}
	case "logout.succeeded":
		// logout.succeeded has no additional details
		a.Details = nil
		return nil
	case "organization.updated":
		details = &OrganizationUpdated{}
	case "project.created":
		details = &ProjectCreated{}
	case "project.updated":
		details = &ProjectUpdated{}
	case "project.archived":
		details = &ProjectArchived{}
	case "rate_limit.updated":
		details = &RateLimitUpdated{}
	case "rate_limit.deleted":
		details = &RateLimitDeleted{}
	case "service_account.created":
		details = &ServiceAccountCreated{}
	case "service_account.updated":
		details = &ServiceAccountUpdated{}
	case "service_account.deleted":
		details = &ServiceAccountDeleted{}
	case "user.added":
		details = &UserAdded{}
	case "user.updated":
		details = &UserUpdated{}
	case "user.deleted":
		details = &UserDeleted{}
	default:
		// For unknown event types, store the raw JSON as a map
		var rawDetails map[string]any
		if err := json.Unmarshal(eventData, &rawDetails); err != nil {
			// If we can't parse as a map, just store nil
			a.Details = nil
			return nil
		}
		a.Details = rawDetails
		return nil
	}

	if err := json.Unmarshal(eventData, details); err != nil {
		return fmt.Errorf("failed to unmarshal details for type %s: %w", raw.Type, err)
	}

	a.Details = details
	return nil
}

func (c *Client) ListAuditLogs(params *AuditLogListParams) (*ListResponse[AuditLog], error) {
	queryParams := make(map[string]string)

	if params != nil {
		if params.Limit > 0 {
			queryParams["limit"] = strconv.Itoa(params.Limit)
		}
		if params.After != "" {
			queryParams["after"] = params.After
		}
		if params.Before != "" {
			queryParams["before"] = params.Before
		}
		if params.EffectiveAt != nil {
			if params.EffectiveAt.Gte != 0 {
				queryParams["effective_at[gte]"] = strconv.FormatInt(params.EffectiveAt.Gte, 10)
			}
			if params.EffectiveAt.Gt != 0 {
				queryParams["effective_at[gt]"] = strconv.FormatInt(params.EffectiveAt.Gt, 10)
			}
			if params.EffectiveAt.Lte != 0 {
				queryParams["effective_at[lte]"] = strconv.FormatInt(params.EffectiveAt.Lte, 10)
			}
			if params.EffectiveAt.Lt != 0 {
				queryParams["effective_at[lt]"] = strconv.FormatInt(params.EffectiveAt.Lt, 10)
			}
		}
		if len(params.ProjectIDs) > 0 {
			queryParams["project_ids"] = strings.Join(params.ProjectIDs, ",")
		}
		if len(params.EventTypes) > 0 {
			queryParams["event_types"] = strings.Join(params.EventTypes, ",")
		}
		if len(params.ActorIDs) > 0 {
			queryParams["actor_ids"] = strings.Join(params.ActorIDs, ",")
		}
		if len(params.ActorEmails) > 0 {
			queryParams["actor_emails"] = strings.Join(params.ActorEmails, ",")
		}
		if len(params.ResourceIDs) > 0 {
			queryParams["resource_ids"] = strings.Join(params.ResourceIDs, ",")
		}
	}

	return Get[AuditLog](c.client, AuditLogsListEndpoint, queryParams)
}

// String returns a human-readable string representation of the AuditLog
func (al *AuditLog) String() string {
	projectInfo := "no project"
	if al.Project != nil {
		projectInfo = fmt.Sprintf("%s(%s)", al.Project.Name, al.Project.ID)
	}

	actorInfo := "unknown"
	switch {
	case al.Actor.Session != nil:
		actorInfo = fmt.Sprintf("user:%s", al.Actor.Session.User.Email)
	case al.Actor.APIKey != nil:
		actorInfo = fmt.Sprintf("apikey:%s", al.Actor.APIKey.User.Email)
	}

	return fmt.Sprintf("AuditLog{ID: %s, Type: %s, Project: %s, Actor: %s, Time: %s}",
		al.ID, al.Type, projectInfo, actorInfo, al.EffectiveAt.String())
}

// MarshalJSON implements json.Marshaler to properly serialize the AuditLog
// including the event-specific details under the dynamic key (e.g., "invite.deleted")
func (al AuditLog) MarshalJSON() ([]byte, error) {
	// Create a map to hold all fields
	result := make(map[string]any)

	// Add standard fields
	if al.Object != "" {
		result["object"] = al.Object
	}
	result["id"] = al.ID
	result["type"] = al.Type
	result["effective_at"] = al.EffectiveAt
	result["actor"] = al.Actor

	if al.Project != nil {
		result["project"] = al.Project
	}

	// Add event-specific details under the dynamic key
	if al.Details != nil && al.Type != "" {
		result[al.Type] = al.Details
	}

	return json.Marshal(result)
}
