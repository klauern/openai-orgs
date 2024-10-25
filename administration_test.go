package openaiorgs

import (
	"encoding/json"
	"testing"
	"time"
)

func TestRoleTypeString(t *testing.T) {
	tests := []struct {
		name     string
		roleType RoleType
		expected string
	}{
		{"Owner", RoleTypeOwner, "owner"},
		{"Member", RoleTypeMember, "member"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.roleType.String(); got != tt.expected {
				t.Errorf("RoleType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCustomTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		{"RFC3339 format", `"2023-08-01T12:34:56Z"`, time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC), false},
		{"Unix timestamp", `1627843200`, time.Unix(1627843200, 0), false},
		{"Invalid format", `"invalid"`, time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct CustomTime
			err := json.Unmarshal([]byte(tt.input), &ct)
			if (err != nil) != tt.wantErr {
				t.Errorf("CustomTime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && time.Time(ct) != tt.expected {
				t.Errorf("CustomTime.UnmarshalJSON() = %v, want %v", ct, tt.expected)
			}
		})
	}
}

func TestCustomTimeString(t *testing.T) {
	tests := []struct {
		name     string
		input    CustomTime
		expected string
	}{
		{"RFC3339 format", CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)), "2023-08-01T12:34:56Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.expected {
				t.Errorf("CustomTime.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCustomTimeTime(t *testing.T) {
	tests := []struct {
		name     string
		input    CustomTime
		expected time.Time
	}{
		{"Time conversion", CustomTime(time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)), time.Date(2023, 8, 1, 12, 34, 56, 0, time.UTC)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.Time(); got != tt.expected {
				t.Errorf("CustomTime.Time() = %v, want %v", got, tt.expected)
			}
		})
	}
}
