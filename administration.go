package openaiorgs

import (
	"encoding/json"
	"time"
)

// UnixSeconds represents a Unix timestamp as a time.Time
type UnixSeconds time.Time

func (us UnixSeconds) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(us).Unix())
}

// UnmarshalJSON implements json.Unmarshaler interface
func (us *UnixSeconds) UnmarshalJSON(data []byte) error {
	// Try parsing as string first (for RFC3339 format)
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err == nil {
		t, err := time.Parse(time.RFC3339, timeStr)
		if err == nil {
			*us = UnixSeconds(t)
			return nil
		}
	}

	// Fall back to parsing as integer (Unix timestamp)
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
