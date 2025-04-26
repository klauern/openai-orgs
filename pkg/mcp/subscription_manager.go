package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// SubscriptionHandler is a function type that handles subscription updates
type SubscriptionHandler func(contents mcp.ResourceContents)

// SubscriptionOptions contains options for creating a subscription
type SubscriptionOptions struct {
	// BufferSize is the size of the channel buffer (default 1)
	BufferSize int
	// Handler is called for each update (optional)
	Handler SubscriptionHandler
}

// DefaultSubscriptionOptions returns the default subscription options
func DefaultSubscriptionOptions() *SubscriptionOptions {
	return &SubscriptionOptions{
		BufferSize: 1,
	}
}

// WithBufferSize sets the channel buffer size
func (o *SubscriptionOptions) WithBufferSize(size int) *SubscriptionOptions {
	o.BufferSize = size
	return o
}

// WithHandler sets the update handler
func (o *SubscriptionOptions) WithHandler(handler SubscriptionHandler) *SubscriptionOptions {
	o.Handler = handler
	return o
}
