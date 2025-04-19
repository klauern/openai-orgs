package openaiorgs

import (
	"fmt"
	"strings"
)

// Common response type for paginated lists
type ListResponse[T any] struct {
	Object  string `json:"object"`
	Data    []T    `json:"data"`
	FirstID string `json:"first_id"`
	LastID  string `json:"last_id"`
	HasMore bool   `json:"has_more"`
}

// String returns a pretty-printed string representation of the ListResponse
func (lr *ListResponse[T]) String() string {
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Object: %s\n", lr.Object))
	result.WriteString(fmt.Sprintf("First ID: %s\n", lr.FirstID))
	result.WriteString(fmt.Sprintf("Last ID: %s\n", lr.LastID))
	result.WriteString(fmt.Sprintf("Has More: %v\n", lr.HasMore))
	result.WriteString("Data:\n")

	for i, item := range lr.Data {
		result.WriteString(fmt.Sprintf("  [%d] %+v\n", i, item))
	}

	return result.String()
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

// Common role types
type RoleType string

const (
	RoleTypeOwner  RoleType = "owner"
	RoleTypeMember RoleType = "member"
)

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

// String returns a human-readable string representation of the Owner
func (o *Owner) String() string {
	ownerInfo := "unknown"
	switch {
	case o.User != nil:
		ownerInfo = fmt.Sprintf("user:%s", o.User.Email)
	case o.SA != nil:
		ownerInfo = fmt.Sprintf("sa:%s", o.SA.Name)
	}
	return fmt.Sprintf("Owner{ID: %s, Name: %s, Type: %s, Info: %s}",
		o.ID, o.Name, o.Type, ownerInfo)
}
