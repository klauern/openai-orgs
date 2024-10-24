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

type CustomTime time.Time

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		// Handle string format
		t, err := time.Parse(time.RFC3339, string(b[1:len(b)-1]))
		if err != nil {
			return err
		}
		*ct = CustomTime(t)
	} else {
		// Handle numeric format (assume Unix timestamp in seconds)
		var timestamp int64
		err := json.Unmarshal(b, &timestamp)
		if err != nil {
			return err
		}
		*ct = CustomTime(time.Unix(timestamp, 0))
	}
	return nil
}

// Add this method to the CustomTime type
func (ct CustomTime) String() string {
	return time.Time(ct).Format(time.RFC3339)
}

func (ct CustomTime) Time() time.Time {
	return time.Time(ct)
}
