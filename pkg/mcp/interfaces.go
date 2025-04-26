package mcp

import (
	"context"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/mark3labs/mcp-go/mcp"
)

// ResourceHandler is a function type that handles resource requests
type ResourceHandler func(context.Context, mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)

// ResourceProvider defines the interface for providing MCP resources
type ResourceProvider interface {
	// GetResource retrieves a resource by its URI and parameters
	GetResource(ctx context.Context, uri string, params map[string]any) (any, error)

	// Subscribe subscribes to updates for a resource
	Subscribe(uri string) (<-chan mcp.ResourceContents, func())
}

// ClientProvider defines the interface for creating OpenAI clients
type ClientProvider interface {
	// NewClient creates a new OpenAI client with the given auth token
	NewClient(authToken string) *openaiorgs.Client
}

// ResourceManager defines the interface for managing MCP resources
type ResourceManager interface {
	// AddResource adds a new resource to the MCP server
	AddResource(resource *mcp.Resource, handler ResourceHandler)

	// ListResources returns all registered resources
	ListResources() []*mcp.Resource
}

// SubscriptionManager defines the interface for managing resource subscriptions
type SubscriptionManager interface {
	// Subscribe creates a new subscription for the given URI
	Subscribe(uri string) chan mcp.ResourceContents

	// Unsubscribe removes a subscription
	Unsubscribe(uri string, ch chan mcp.ResourceContents)

	// Notify notifies all subscribers of a resource update
	Notify(uri string, contents mcp.ResourceContents)
}

// DefaultClientProvider implements the ClientProvider interface
type DefaultClientProvider struct{}

func (p *DefaultClientProvider) NewClient(authToken string) *openaiorgs.Client {
	return openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
}
