package openaiorgs

import (
	"fmt"
	"strings"
)

// ListResponse is a generic container for paginated API responses.
// It provides metadata about the list (such as pagination tokens) and the actual data items.
type ListResponse[T any] struct {
	// Object identifies the type of the response, typically "list".
	Object string `json:"object"`
	// Data contains the actual list of items of type T.
	Data []T `json:"data"`
	// FirstID is the ID of the first item in this page.
	FirstID string `json:"first_id"`
	// LastID is the ID of the last item in this page.
	LastID string `json:"last_id"`
	// HasMore indicates whether there are more items available in subsequent pages.
	HasMore bool `json:"has_more"`
}

// String returns a pretty-printed string representation of the ListResponse.
// It includes all metadata and a formatted list of all items in the Data field.
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

// Owner represents an entity that can own or be assigned to resources.
// An owner can be either a user or a service account, as indicated by the Type field.
type Owner struct {
	// Object identifies the type of this resource.
	Object string `json:"object"`
	// ID is the unique identifier for this owner.
	ID string `json:"id"`
	// Name is the display name of the owner.
	Name string `json:"name"`
	// Type indicates whether this is a user or service account owner.
	Type OwnerType `json:"type"`
	// User contains user-specific details if Type is OwnerTypeUser.
	User *User `json:"user,omitempty"`
	// SA contains service account details if Type is OwnerTypeServiceAccount.
	SA *ProjectServiceAccount `json:"service_account,omitempty"`
}

// OwnerType represents the type of an owner entity.
// It distinguishes between user accounts and service accounts.
type OwnerType string

const (
	// OwnerTypeUser indicates the owner is a human user account.
	OwnerTypeUser OwnerType = "user"
	// OwnerTypeServiceAccount indicates the owner is a service account.
	OwnerTypeServiceAccount OwnerType = "service_account"
)

// RoleType represents the level of access and permissions an entity has.
type RoleType string

const (
	// RoleTypeOwner grants full administrative access.
	RoleTypeOwner RoleType = "owner"
	// RoleTypeMember grants standard member access.
	RoleTypeMember RoleType = "member"
)

// ParseRoleType converts a string to a RoleType.
// Returns an empty RoleType if the string doesn't match a known role.
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

// String returns a human-readable string representation of the Owner.
// It includes basic metadata and owner-specific information based on the owner type.
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
