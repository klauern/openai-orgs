package openaiorgs

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnixSeconds_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid unix timestamp",
			input:   "1640995200",
			want:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "valid RFC3339 string",
			input:   `"2022-01-01T00:00:00Z"`,
			want:    time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid RFC3339 string",
			input:   `"2022-13-01T00:00:00Z"`,
			wantErr: true,
		},
		{
			name:    "invalid json",
			input:   "not_json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got UnixSeconds
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnixSeconds.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !got.Time().Equal(tt.want) {
				t.Errorf("UnixSeconds.UnmarshalJSON() = %v, want %v", got.Time(), tt.want)
			}
		})
	}
}

func TestUnixSeconds_String(t *testing.T) {
	timestamp := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	us := UnixSeconds(timestamp)

	want := "2022-01-01T00:00:00Z"
	if got := us.String(); got != want {
		t.Errorf("UnixSeconds.String() = %v, want %v", got, want)
	}
}

func TestUnixSeconds_Time(t *testing.T) {
	timestamp := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	us := UnixSeconds(timestamp)

	if got := us.Time(); !got.Equal(timestamp) {
		t.Errorf("UnixSeconds.Time() = %v, want %v", got, timestamp)
	}
}

func TestParseRoleType(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want RoleType
	}{
		{
			name: "owner",
			args: args{s: "owner"},
			want: RoleTypeOwner,
		},
		{
			name: "member",
			args: args{s: "member"},
			want: RoleTypeMember,
		},
		{
			name: "unknown",
			args: args{s: "unknown"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseRoleType(tt.args.s); got != tt.want {
				t.Errorf("ParseRoleType() = %v, want %v", got, tt.want)
			}
		})
	}
}
