package openaiorgs

import (
	"encoding/json"
	"time"
)

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

type UnixSeconds time.Time

func (ct *UnixSeconds) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		// Handle string format
		t, err := time.Parse(time.RFC3339, string(b[1:len(b)-1]))
		if err != nil {
			return err
		}
		*ct = UnixSeconds(t)
	} else {
		// Handle numeric format (assume Unix timestamp in seconds)
		var timestamp int64
		err := json.Unmarshal(b, &timestamp)
		if err != nil {
			return err
		}
		*ct = UnixSeconds(time.Unix(timestamp, 0))
	}
	return nil
}

// Add this method to the CustomTime type
func (ct UnixSeconds) String() string {
	return time.Time(ct).Format(time.RFC3339)
}

func (ct UnixSeconds) Time() time.Time {
	return time.Time(ct)
}
