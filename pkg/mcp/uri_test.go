package mcp

import "testing"

func TestParseURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		want    *ResourceURI
		wantErr bool
	}{
		{"valid project", "openai-orgs://project/proj_123", &ResourceURI{Type: "project", ProjectID: "proj_123"}, false},
		{"valid project with service account", "openai-orgs://project/proj_123/service-account/sa_456", &ResourceURI{Type: "project", ProjectID: "proj_123", ServiceAccount: "sa_456"}, false},
		{"valid members", "openai-orgs://members/user_789", &ResourceURI{Type: "members", MemberID: "user_789"}, false},
		{"valid active-projects", "openai-orgs://active-projects", &ResourceURI{Type: "active-projects"}, false},
		{"valid current-members", "openai-orgs://current-members", &ResourceURI{Type: "current-members"}, false},
		{"valid usage-dashboard", "openai-orgs://usage-dashboard", &ResourceURI{Type: "usage-dashboard"}, false},
		{"project without id", "openai-orgs://project", &ResourceURI{Type: "project"}, false},
		{"members without id", "openai-orgs://members", &ResourceURI{Type: "members"}, false},
		{"wrong prefix", "http://wrong/path", nil, true},
		{"unknown type", "openai-orgs://unknown", nil, true},
		{"empty after prefix", "openai-orgs://", nil, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseURI(tc.uri)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseURI(%q) error = %v, wantErr %v", tc.uri, err, tc.wantErr)
				return
			}
			if tc.wantErr {
				return
			}
			if got.Type != tc.want.Type {
				t.Errorf("Type = %q, want %q", got.Type, tc.want.Type)
			}
			if got.ProjectID != tc.want.ProjectID {
				t.Errorf("ProjectID = %q, want %q", got.ProjectID, tc.want.ProjectID)
			}
			if got.ServiceAccount != tc.want.ServiceAccount {
				t.Errorf("ServiceAccount = %q, want %q", got.ServiceAccount, tc.want.ServiceAccount)
			}
			if got.MemberID != tc.want.MemberID {
				t.Errorf("MemberID = %q, want %q", got.MemberID, tc.want.MemberID)
			}
		})
	}
}

func TestResourceURI_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		uri  string
	}{
		{"project", "openai-orgs://project/proj_123"},
		{"project with sa", "openai-orgs://project/proj_123/service-account/sa_456"},
		{"members", "openai-orgs://members/user_789"},
		{"active-projects", "openai-orgs://active-projects"},
		{"current-members", "openai-orgs://current-members"},
		{"usage-dashboard", "openai-orgs://usage-dashboard"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := ParseURI(tc.uri)
			if err != nil {
				t.Fatalf("ParseURI failed: %v", err)
			}
			got := parsed.String()
			if got != tc.uri {
				t.Errorf("round-trip: got %q, want %q", got, tc.uri)
			}
		})
	}
}
