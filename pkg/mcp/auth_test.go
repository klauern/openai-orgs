package mcp

import (
	"context"
	"testing"
)

func TestAuthFromEnvironment(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		wantSet  bool
	}{
		{"with valid key", "sk-test-key", true},
		{"with empty key", "", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				t.Setenv("OPENAI_API_KEY", tc.envValue)
			} else {
				t.Setenv("OPENAI_API_KEY", "")
			}
			ctx := AuthFromEnvironment(context.Background())
			val := ctx.Value(authToken{})
			if tc.wantSet && val == nil {
				t.Error("expected token in context")
			}
			if !tc.wantSet && val != nil {
				t.Error("expected no token in context")
			}
		})
	}
}

func TestWithAuthToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  bool
	}{
		{"empty string", "", false},
		{"non-empty", "sk-123", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := withAuthToken(context.Background(), tc.token)
			val := ctx.Value(authToken{})
			if tc.want && val == nil {
				t.Error("expected token in context")
			}
			if !tc.want && val != nil {
				t.Error("expected no token in context")
			}
			if tc.want {
				got, ok := val.(string)
				if !ok || got != tc.token {
					t.Errorf("got %q, want %q", got, tc.token)
				}
			}
		})
	}
}
