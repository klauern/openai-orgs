package mcp

import (
	"context"
	"fmt"
	"math"
)

// authTokenFromContext safely extracts the auth token from context
func authTokenFromContext(ctx context.Context) (string, error) {
	v := ctx.Value(authToken{})
	if v == nil {
		return "", ErrNoAuthToken
	}
	token, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("auth token is not a string: %T", v)
	}
	if token == "" {
		return "", ErrEmptyToken
	}
	return token, nil
}

// requireString safely extracts a required string parameter
func requireString(params map[string]any, key string) (string, error) {
	v, ok := params[key]
	if !ok {
		return "", fmt.Errorf("missing required parameter: %s", key)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("parameter '%s' must be a string, got %T", key, v)
	}
	return s, nil
}

// optionalString safely extracts an optional string parameter
func optionalString(params map[string]any, key string) (string, bool, error) {
	v, ok := params[key]
	if !ok {
		return "", false, nil
	}
	s, ok := v.(string)
	if !ok {
		return "", false, fmt.Errorf("parameter '%s' must be a string, got %T", key, v)
	}
	return s, true, nil
}

// optionalBool safely extracts an optional bool parameter
func optionalBool(params map[string]any, key string) (bool, bool, error) {
	v, ok := params[key]
	if !ok {
		return false, false, nil
	}
	b, ok := v.(bool)
	if !ok {
		return false, false, fmt.Errorf("parameter '%s' must be a bool, got %T", key, v)
	}
	return b, true, nil
}

// optionalIntFromFloat safely extracts an optional int from a float64 parameter (JSON numbers come as float64)
func optionalIntFromFloat(params map[string]any, key string) (int, bool, error) {
	v, ok := params[key]
	if !ok {
		return 0, false, nil
	}
	f, ok := v.(float64)
	if !ok {
		return 0, false, fmt.Errorf("parameter '%s' must be a number, got %T", key, v)
	}
	if f < 0 {
		return 0, false, fmt.Errorf("parameter '%s' must be non-negative, got %v", key, f)
	}
	if f != math.Trunc(f) {
		return 0, false, fmt.Errorf("parameter '%s' must be a whole number, got %v", key, f)
	}
	if f > math.MaxInt32 {
		return 0, false, fmt.Errorf("parameter '%s' value %v exceeds maximum", key, f)
	}
	return int(f), true, nil
}
