package mcp

import (
	"fmt"
	"strings"
)

const (
	uriPrefix = "openai-orgs://"
)

// ResourceURI represents a parsed OpenAI organizations URI
type ResourceURI struct {
	Type           string
	ProjectID      string
	ServiceAccount string
	MemberID       string
}

// ParseURI parses a raw URI string into a structured ResourceURI
func ParseURI(uri string) (*ResourceURI, error) {
	if !strings.HasPrefix(uri, uriPrefix) {
		return nil, fmt.Errorf("invalid URI prefix: must start with %s", uriPrefix)
	}

	path := strings.TrimPrefix(uri, uriPrefix)
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid URI: empty path")
	}

	r := &ResourceURI{
		Type: parts[0],
	}

	switch r.Type {
	case "project":
		if len(parts) >= 2 {
			r.ProjectID = parts[1]
		}
		if len(parts) >= 4 && parts[2] == "service-account" {
			r.ServiceAccount = parts[3]
		}
	case "members":
		if len(parts) >= 2 {
			r.MemberID = parts[1]
		}
	case "active-projects", "current-members", "usage-dashboard":
		// These are valid static resources with no additional parsing needed
	default:
		return nil, fmt.Errorf("invalid resource type: %s", r.Type)
	}

	return r, nil
}

// String converts the ResourceURI back to its string representation
func (r *ResourceURI) String() string {
	var builder strings.Builder
	builder.WriteString(uriPrefix)
	builder.WriteString(r.Type)

	switch r.Type {
	case "project":
		if r.ProjectID != "" {
			builder.WriteString("/")
			builder.WriteString(r.ProjectID)
			if r.ServiceAccount != "" {
				builder.WriteString("/service-account/")
				builder.WriteString(r.ServiceAccount)
			}
		}
	case "members":
		if r.MemberID != "" {
			builder.WriteString("/")
			builder.WriteString(r.MemberID)
		}
	}

	return builder.String()
}
