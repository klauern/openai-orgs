package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Resource types and MIME types
const (
	ResourceTypeActiveProjects = "active-projects"
	ResourceTypeCurrentMembers = "current-members"
	ResourceTypeUsageDashboard = "usage-dashboard"

	MIMETypeProjectList    = "application/vnd.openai-orgs.project-list+json"
	MIMETypeMemberList     = "application/vnd.openai-orgs.member-list+json"
	MIMETypeUsageDashboard = "application/vnd.openai-orgs.usage+json"

	defaultPageSize = 20
	maxPageSize     = 100
)

// resourceHandler is a generic handler for resources
type resourceHandler func(ctx context.Context, client *openaiorgs.Client, params map[string]any) (any, error)

// AddResources adds static resource capabilities to the MCP server
func AddResources(s *server.MCPServer) {
	resources := []struct {
		uri      string
		name     string
		desc     string
		mimeType string
		handler  resourceHandler
	}{
		{
			uri:      "openai-orgs://active-projects",
			name:     "Active Projects",
			desc:     "Currently active projects in the organization",
			mimeType: MIMETypeProjectList,
			handler:  handleActiveProjects,
		},
		{
			uri:      "openai-orgs://current-members",
			name:     "Current Members",
			desc:     "Current organization members",
			mimeType: MIMETypeMemberList,
			handler:  handleCurrentMembers,
		},
		{
			uri:      "openai-orgs://usage-dashboard",
			name:     "Usage Dashboard",
			desc:     "Current usage statistics dashboard",
			mimeType: MIMETypeUsageDashboard,
			handler:  handleUsageDashboard,
		},
	}

	for _, r := range resources {
		resource := mcp.NewResource(
			r.uri,
			r.name,
			mcp.WithResourceDescription(r.desc),
			mcp.WithMIMEType(r.mimeType),
		)

		handler := createResourceHandler(r.handler, r.mimeType)
		s.AddResource(resource, handler)
	}

	// Start background polling for changes
	go pollForChanges(context.Background())
}

// createResourceHandler creates a standard MCP handler from our resource handler
func createResourceHandler(h resourceHandler, mimeType string) func(context.Context, mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		authToken := ctx.Value(authToken{}).(string)
		if authToken == "" {
			return nil, ErrNoAuthToken
		}

		client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)

		// Extract pagination and other parameters
		params := make(map[string]any)
		if request.Params.Arguments != nil {
			params = request.Params.Arguments
		}

		data, err := h(ctx, client, params)
		if err != nil {
			return nil, err
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		contents := &mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: mimeType,
			Text:     string(jsonData),
		}

		// Handle subscription if requested
		if sub, ok := request.Params.Arguments["subscribe"].(bool); ok && sub {
			ch := subManager.subscribe(request.Params.URI)
			go func() {
				<-ctx.Done()
				subManager.unsubscribe(request.Params.URI, ch)
			}()
		}

		return []mcp.ResourceContents{contents}, nil
	}
}

// Individual handlers for each resource type
func handleActiveProjects(_ context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
	limit, after := getPaginationFromParams(params)
	return client.ListProjects(limit, after, true)
}

func handleCurrentMembers(_ context.Context, client *openaiorgs.Client, params map[string]any) (any, error) {
	limit, after := getPaginationFromParams(params)
	return client.ListUsers(limit, after)
}

func handleUsageDashboard(_ context.Context, client *openaiorgs.Client, _ map[string]any) (any, error) {
	startTime := time.Now().AddDate(0, -1, 0).Format(time.RFC3339) // Last month
	params := map[string]string{"start_time": startTime}

	usageData := make(map[string]any)

	if completions, err := client.GetCompletionsUsage(params); err == nil {
		usageData["completions"] = completions
	}
	if embeddings, err := client.GetEmbeddingsUsage(params); err == nil {
		usageData["embeddings"] = embeddings
	}
	if images, err := client.GetImagesUsage(params); err == nil {
		usageData["images"] = images
	}

	return usageData, nil
}

// Helper functions

func getPaginationFromParams(params map[string]any) (limit int, after string) {
	limit = defaultPageSize

	if pagination, ok := params["pagination"].(map[string]any); ok {
		if afterVal, ok := pagination["after"].(string); ok {
			after = afterVal
		}
		if limitVal, ok := pagination["limit"].(float64); ok {
			limit = int(limitVal)
		}
	}

	if limit <= 0 {
		limit = defaultPageSize
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}

	return limit, after
}

// Subscription management
type subscriptionManager struct {
	sync.RWMutex
	subscribers map[string][]chan mcp.ResourceContents
}

var subManager = &subscriptionManager{
	subscribers: make(map[string][]chan mcp.ResourceContents),
}

func (sm *subscriptionManager) subscribe(uri string) chan mcp.ResourceContents {
	sm.Lock()
	defer sm.Unlock()
	ch := make(chan mcp.ResourceContents, 1)
	sm.subscribers[uri] = append(sm.subscribers[uri], ch)
	return ch
}

func (sm *subscriptionManager) unsubscribe(uri string, ch chan mcp.ResourceContents) {
	sm.Lock()
	defer sm.Unlock()
	subs := sm.subscribers[uri]
	for i, sub := range subs {
		if sub == ch {
			sm.subscribers[uri] = append(subs[:i], subs[i+1:]...)
			close(ch)
			break
		}
	}
}

func (sm *subscriptionManager) notify(uri string, contents mcp.ResourceContents) {
	sm.RLock()
	defer sm.RUnlock()
	for _, ch := range sm.subscribers[uri] {
		select {
		case ch <- contents:
		default:
			// Channel is full, skip notification
		}
	}
}

// Background polling
func pollForChanges(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			checkForResourceChanges(ctx)
		}
	}
}

func checkForResourceChanges(ctx context.Context) {
	subManager.RLock()
	defer subManager.RUnlock()

	for uri := range subManager.subscribers {
		parsedURI, err := ParseURI(uri)
		if err != nil {
			continue
		}

		switch parsedURI.Type {
		case ResourceTypeActiveProjects:
			go updateActiveProjects(ctx)
		case ResourceTypeCurrentMembers:
			go updateCurrentMembers(ctx)
		case ResourceTypeUsageDashboard:
			go updateUsageDashboard(ctx)
		}
	}
}

// Update functions for each resource type
func updateActiveProjects(ctx context.Context) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return
	}

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
	projects, err := client.ListProjects(defaultPageSize, "", true)
	if err != nil {
		return
	}

	data, err := json.Marshal(projects)
	if err != nil {
		return
	}

	contents := &mcp.TextResourceContents{
		URI:      "openai-orgs://active-projects",
		MIMEType: MIMETypeProjectList,
		Text:     string(data),
	}

	subManager.notify("openai-orgs://active-projects", contents)
}

func updateCurrentMembers(ctx context.Context) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return
	}

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
	members, err := client.ListUsers(defaultPageSize, "")
	if err != nil {
		return
	}

	data, err := json.Marshal(members)
	if err != nil {
		return
	}

	contents := &mcp.TextResourceContents{
		URI:      "openai-orgs://current-members",
		MIMEType: MIMETypeMemberList,
		Text:     string(data),
	}

	subManager.notify("openai-orgs://current-members", contents)
}

func updateUsageDashboard(ctx context.Context) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return
	}

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
	startTime := time.Now().AddDate(0, -1, 0).Format(time.RFC3339)
	params := map[string]string{"start_time": startTime}

	usageData := make(map[string]any)
	if completions, err := client.GetCompletionsUsage(params); err == nil {
		usageData["completions"] = completions
	}
	if embeddings, err := client.GetEmbeddingsUsage(params); err == nil {
		usageData["embeddings"] = embeddings
	}
	if images, err := client.GetImagesUsage(params); err == nil {
		usageData["images"] = images
	}

	data, err := json.Marshal(usageData)
	if err != nil {
		return
	}

	contents := &mcp.TextResourceContents{
		URI:      "openai-orgs://usage-dashboard",
		MIMEType: MIMETypeUsageDashboard,
		Text:     string(data),
	}

	subManager.notify("openai-orgs://usage-dashboard", contents)
}
