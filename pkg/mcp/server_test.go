package mcp

import (
	"context"
	"testing"
)

func TestNewMCPServer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := NewMCPServer(ctx)
	if s == nil {
		t.Fatal("NewMCPServer returned nil")
	}
}
