package mcp

import (
	"context"
	"math"
	"testing"
)

func TestAuthTokenFromContext(t *testing.T) {
	t.Run("nil context value", func(t *testing.T) {
		ctx := context.Background()
		_, err := authTokenFromContext(ctx)
		if err == nil {
			t.Fatal("expected error for nil context value")
		}
		if err != ErrNoAuthToken {
			t.Errorf("expected ErrNoAuthToken, got: %v", err)
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), authToken{}, 12345)
		_, err := authTokenFromContext(ctx)
		if err == nil {
			t.Fatal("expected error for wrong type")
		}
		if err == ErrNoAuthToken || err == ErrEmptyToken {
			t.Errorf("expected type error, got: %v", err)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), authToken{}, "")
		_, err := authTokenFromContext(ctx)
		if err == nil {
			t.Fatal("expected error for empty string")
		}
		if err != ErrEmptyToken {
			t.Errorf("expected ErrEmptyToken, got: %v", err)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), authToken{}, "sk-test-123")
		token, err := authTokenFromContext(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "sk-test-123" {
			t.Errorf("expected 'sk-test-123', got %q", token)
		}
	})
}

func TestRequireString(t *testing.T) {
	t.Run("missing key", func(t *testing.T) {
		params := map[string]any{}
		_, err := requireString(params, "name")
		if err == nil {
			t.Fatal("expected error for missing key")
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		params := map[string]any{"name": 123}
		_, err := requireString(params, "name")
		if err == nil {
			t.Fatal("expected error for wrong type")
		}
	})

	t.Run("valid string", func(t *testing.T) {
		params := map[string]any{"name": "test-project"}
		val, err := requireString(params, "name")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "test-project" {
			t.Errorf("expected 'test-project', got %q", val)
		}
	})
}

func TestOptionalString(t *testing.T) {
	t.Run("missing key", func(t *testing.T) {
		params := map[string]any{}
		val, ok, err := optionalString(params, "after")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Error("expected ok=false for missing key")
		}
		if val != "" {
			t.Errorf("expected empty string, got %q", val)
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		params := map[string]any{"after": 42}
		_, _, err := optionalString(params, "after")
		if err == nil {
			t.Fatal("expected error for wrong type")
		}
	})

	t.Run("valid string", func(t *testing.T) {
		params := map[string]any{"after": "cursor-123"}
		val, ok, err := optionalString(params, "after")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Error("expected ok=true")
		}
		if val != "cursor-123" {
			t.Errorf("expected 'cursor-123', got %q", val)
		}
	})
}

func TestOptionalBool(t *testing.T) {
	t.Run("missing key", func(t *testing.T) {
		params := map[string]any{}
		val, ok, err := optionalBool(params, "activeOnly")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Error("expected ok=false for missing key")
		}
		if val != false {
			t.Errorf("expected false, got %v", val)
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		params := map[string]any{"activeOnly": "yes"}
		_, _, err := optionalBool(params, "activeOnly")
		if err == nil {
			t.Fatal("expected error for wrong type")
		}
	})

	t.Run("valid bool", func(t *testing.T) {
		params := map[string]any{"activeOnly": true}
		val, ok, err := optionalBool(params, "activeOnly")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Error("expected ok=true")
		}
		if val != true {
			t.Errorf("expected true, got %v", val)
		}
	})
}

func TestOptionalIntFromFloat(t *testing.T) {
	t.Run("missing key", func(t *testing.T) {
		params := map[string]any{}
		val, ok, err := optionalIntFromFloat(params, "limit")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ok {
			t.Error("expected ok=false for missing key")
		}
		if val != 0 {
			t.Errorf("expected 0, got %d", val)
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		params := map[string]any{"limit": "fifty"}
		_, _, err := optionalIntFromFloat(params, "limit")
		if err == nil {
			t.Fatal("expected error for wrong type")
		}
	})

	t.Run("negative", func(t *testing.T) {
		params := map[string]any{"limit": float64(-5)}
		_, _, err := optionalIntFromFloat(params, "limit")
		if err == nil {
			t.Fatal("expected error for negative value")
		}
	})

	t.Run("fractional", func(t *testing.T) {
		params := map[string]any{"limit": float64(3.5)}
		_, _, err := optionalIntFromFloat(params, "limit")
		if err == nil {
			t.Fatal("expected error for fractional value")
		}
	})

	t.Run("overflow", func(t *testing.T) {
		params := map[string]any{"limit": float64(math.MaxInt32 + 1)}
		_, _, err := optionalIntFromFloat(params, "limit")
		if err == nil {
			t.Fatal("expected error for overflow value")
		}
	})

	t.Run("zero", func(t *testing.T) {
		params := map[string]any{"limit": float64(0)}
		val, ok, err := optionalIntFromFloat(params, "limit")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Error("expected ok=true")
		}
		if val != 0 {
			t.Errorf("expected 0, got %d", val)
		}
	})

	t.Run("valid", func(t *testing.T) {
		params := map[string]any{"limit": float64(50)}
		val, ok, err := optionalIntFromFloat(params, "limit")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Error("expected ok=true")
		}
		if val != 50 {
			t.Errorf("expected 50, got %d", val)
		}
	})
}
