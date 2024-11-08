package openaiorgs

import (
	"time"
)

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
