package openaiorgs

import "time"

// Common response type for paginated lists
type ListResponse[T any] struct {
	Object  string `json:"object"`
	Data    []T    `json:"data"`
	FirstID string `json:"first_id"`
	LastID  string `json:"last_id"`
	HasMore bool   `json:"has_more"`
}

// Common owner types
type Owner struct {
	Object string                 `json:"object"`
	ID     string                 `json:"id"`
	Name   string                 `json:"name"`
	Type   OwnerType              `json:"type"`
	User   *User                  `json:"user,omitempty"`
	SA     *ProjectServiceAccount `json:"service_account,omitempty"`
}

type OwnerType string

const (
	OwnerTypeUser           OwnerType = "user"
	OwnerTypeServiceAccount OwnerType = "service_account"
)

// Common time handling
type UnixSeconds time.Time

// Common role types
type RoleType string

const (
	RoleTypeOwner  RoleType = "owner"
	RoleTypeMember RoleType = "member"
)

func (rt RoleType) String() string {
	return string(rt)
}

func ParseRoleType(s string) RoleType {
	switch s {
	case "owner":
		return RoleTypeOwner
	case "member":
		return RoleTypeMember
	default:
		return ""
	}
}

func (us UnixSeconds) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(us).Unix())
}

func (us *UnixSeconds) UnmarshalJSON(data []byte) error {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err == nil {
		t, err := time.Parse(time.RFC3339, timeStr)
		if err == nil {
			*us = UnixSeconds(t)
			return nil
		}
	}

	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	*us = UnixSeconds(time.Unix(timestamp, 0))
	return nil
}

func (ct UnixSeconds) String() string {
	return time.Time(ct).Format(time.RFC3339)
}

func (ct UnixSeconds) Time() time.Time {
	return time.Time(ct)
}
