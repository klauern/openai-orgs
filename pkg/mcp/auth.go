package mcp

import (
	"context"
	"errors"
	"os"
)

type authToken struct{}

var (
	ErrNoAuthToken = errors.New("no auth token found in context")
	ErrEmptyToken  = errors.New("empty auth token provided")
)

func withAuthToken(c context.Context, auth string) context.Context {
	if auth == "" {
		return c
	}
	return context.WithValue(c, authToken{}, auth)
}

func AuthFromEnvironment(c context.Context) context.Context {
	token := os.Getenv("OPENAI_API_KEY")
	return withAuthToken(c, token)
}
