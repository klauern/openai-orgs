package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Resource types
const (
	ResourceTypeProject = "project"
	ResourceTypeMember  = "member"
	ResourceTypeUsage   = "usage"
)

// MIME types for our resources
const (
	MIMETypeProject = "application/vnd.openai-orgs.project+json"
	MIMETypeMember  = "application/vnd.openai-orgs.member+json"
	MIMETypeUsage   = "application/vnd.openai-orgs.usage+json"
	MIMETypeList    = "application/vnd.openai-orgs.list+json"
)

// Default pagination values
const (
	defaultPageSize = 20
	maxPageSize     = 100
)

// ResourceParams extends the base ReadResourceRequest parameters
type ResourceParams struct {
	URI       string                 `json:"uri"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	Subscribe bool                   `json:"subscribe,omitempty"`
	Context   context.Context        `json:"-"`
}

// ResourceRequest wraps the resource request with extended parameters
type ResourceRequest struct {
	Params ResourceParams
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

// AddResources adds resource capabilities to the MCP server
func AddResources(s *server.MCPServer) {
	// Create resource templates for each type
	projectTemplate := mcp.NewResourceTemplate(
		"openai-orgs://projects/{id}",
		"OpenAI Project",
		mcp.WithTemplateDescription("OpenAI project information"),
		mcp.WithTemplateMIMEType(MIMETypeProject),
	)

	memberTemplate := mcp.NewResourceTemplate(
		"openai-orgs://members/{id}",
		"Organization Member",
		mcp.WithTemplateDescription("OpenAI organization member information"),
		mcp.WithTemplateMIMEType(MIMETypeMember),
	)

	usageTemplate := mcp.NewResourceTemplate(
		"openai-orgs://usage/{id}",
		"Usage Statistics",
		mcp.WithTemplateDescription("OpenAI organization usage statistics"),
		mcp.WithTemplateMIMEType(MIMETypeUsage),
	)

	// List templates
	projectListTemplate := mcp.NewResourceTemplate(
		"openai-orgs://projects",
		"OpenAI Projects List",
		mcp.WithTemplateDescription("List of OpenAI projects"),
		mcp.WithTemplateMIMEType(MIMETypeList),
	)

	memberListTemplate := mcp.NewResourceTemplate(
		"openai-orgs://members",
		"Organization Members List",
		mcp.WithTemplateDescription("List of organization members"),
		mcp.WithTemplateMIMEType(MIMETypeList),
	)

	usageListTemplate := mcp.NewResourceTemplate(
		"openai-orgs://usage",
		"Usage Statistics List",
		mcp.WithTemplateDescription("List of usage statistics"),
		mcp.WithTemplateMIMEType(MIMETypeList),
	)

	// Add resource templates with their handlers
	s.AddResourceTemplate(projectTemplate, handleProjectResource)
	s.AddResourceTemplate(memberTemplate, handleMemberResource)
	s.AddResourceTemplate(usageTemplate, handleUsageResource)
	s.AddResourceTemplate(projectListTemplate, handleListProjectsResource)
	s.AddResourceTemplate(memberListTemplate, handleListMembersResource)
	s.AddResourceTemplate(usageListTemplate, handleListUsageResource)

	// Start background polling for changes
	go pollForChanges(context.Background())
}

// getPaginationParams extracts pagination parameters from the request
func getPaginationParams(request mcp.ReadResourceRequest) (page int, limit string) {
	if args, ok := request.Params.Arguments["pagination"].(map[string]interface{}); ok {
		if pageNum, ok := args["page"].(float64); ok {
			page = int(pageNum)
		}
		if limitStr, ok := args["limit"].(string); ok {
			limit = limitStr
		}
	}

	if page < 1 {
		page = 1
	}

	if limit == "" {
		limit = strconv.Itoa(defaultPageSize)
	}

	if limitInt, err := strconv.Atoi(limit); err == nil && limitInt > maxPageSize {
		limit = strconv.Itoa(maxPageSize)
	}

	return page, limit
}

// handleListProjectsResource handles project list requests
func handleListProjectsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return nil, ErrNoAuthToken
	}

	page, limit := getPaginationParams(request)
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)

	projects, err := client.ListProjects(page, limit, false)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}

	data, err := json.Marshal(projects)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal projects data: %w", err)
	}

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: MIMETypeList,
			Text:     string(data),
		},
	}, nil
}

// handleListMembersResource handles member list requests
func handleListMembersResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return nil, ErrNoAuthToken
	}

	page, limit := getPaginationParams(request)
	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)

	// For now, we'll use ListUsers as a proxy for members since the API doesn't have a direct members endpoint
	users, err := client.ListUsers(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list members: %w", err)
	}

	data, err := json.Marshal(users)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal members data: %w", err)
	}

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: MIMETypeList,
			Text:     string(data),
		},
	}, nil
}

// handleListUsageResource handles usage list requests
func handleListUsageResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return nil, ErrNoAuthToken
	}

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)

	// For usage, we'll use GetCompletionsUsage since there's no list endpoint
	// We'll paginate the results client-side if needed
	usage, err := client.GetCompletionsUsage(map[string]string{
		"start_time": "0", // Get all available data
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get usage data: %w", err)
	}

	data, err := json.Marshal(usage)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal usage data: %w", err)
	}

	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: MIMETypeList,
			Text:     string(data),
		},
	}, nil
}

// handleProjectResource handles project resource requests
func handleProjectResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return nil, ErrNoAuthToken
	}

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
	projectID := request.Params.URI[len("openai-orgs://projects/"):]

	project, err := client.RetrieveProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve project: %w", err)
	}

	data, err := json.Marshal(project)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal project data: %w", err)
	}

	contents := &mcp.TextResourceContents{
		URI:      request.Params.URI,
		MIMEType: MIMETypeProject,
		Text:     string(data),
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

// handleMemberResource handles member resource requests
func handleMemberResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return nil, ErrNoAuthToken
	}

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
	userID := request.Params.URI[len("openai-orgs://members/"):]

	user, err := client.RetrieveUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	data, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user data: %w", err)
	}

	contents := &mcp.TextResourceContents{
		URI:      request.Params.URI,
		MIMEType: MIMETypeMember,
		Text:     string(data),
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

// handleUsageResource handles usage resource requests
func handleUsageResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return nil, ErrNoAuthToken
	}

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
	projectID := request.Params.URI[len("openai-orgs://usage/"):]

	// Get all types of usage for comprehensive data
	usageData := make(map[string]interface{})

	// Get completions usage
	completions, err := client.GetCompletionsUsage(map[string]string{
		"start_time": "0",
		"project_id": projectID,
	})
	if err == nil {
		usageData["completions"] = completions
	}

	// Get embeddings usage
	embeddings, err := client.GetEmbeddingsUsage(map[string]string{
		"start_time": "0",
		"project_id": projectID,
	})
	if err == nil {
		usageData["embeddings"] = embeddings
	}

	// Get images usage
	images, err := client.GetImagesUsage(map[string]string{
		"start_time": "0",
		"project_id": projectID,
	})
	if err == nil {
		usageData["images"] = images
	}

	data, err := json.Marshal(usageData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal usage data: %w", err)
	}

	contents := &mcp.TextResourceContents{
		URI:      request.Params.URI,
		MIMEType: MIMETypeUsage,
		Text:     string(data),
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

// pollForChanges periodically checks for changes in resources
func pollForChanges(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Poll for changes and notify subscribers
			checkForResourceChanges(ctx)
		}
	}
}

func checkForResourceChanges(ctx context.Context) {
	// This is a simplified implementation. In production, you'd want to:
	// 1. Keep track of last known state
	// 2. Compare with current state
	// 3. Only notify on actual changes
	// 4. Handle rate limiting
	// 5. Implement proper error handling and backoff

	subManager.RLock()
	defer subManager.RUnlock()

	for uri := range subManager.subscribers {
		// Based on URI pattern, fetch appropriate resource
		switch {
		case strings.HasPrefix(uri, "openai-orgs://projects/"):
			go updateProjectResource(ctx, uri)
		case strings.HasPrefix(uri, "openai-orgs://members/"):
			go updateMemberResource(ctx, uri)
		case strings.HasPrefix(uri, "openai-orgs://usage/"):
			go updateUsageResource(ctx, uri)
		}
	}
}

// updateProjectResource updates a project's resource data
func updateProjectResource(ctx context.Context, uri string) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return
	}

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
	projectID := uri[len("openai-orgs://projects/"):]

	project, err := client.RetrieveProject(projectID)
	if err != nil {
		return
	}

	data, err := json.Marshal(project)
	if err != nil {
		return
	}

	contents := &mcp.TextResourceContents{
		URI:      uri,
		MIMEType: MIMETypeProject,
		Text:     string(data),
	}

	subManager.notify(uri, contents)
}

// updateMemberResource updates a member's resource data
func updateMemberResource(ctx context.Context, uri string) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return
	}

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
	userID := uri[len("openai-orgs://members/"):]

	user, err := client.RetrieveUser(userID)
	if err != nil {
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	contents := &mcp.TextResourceContents{
		URI:      uri,
		MIMEType: MIMETypeMember,
		Text:     string(data),
	}

	subManager.notify(uri, contents)
}

// updateUsageResource updates usage resource data
func updateUsageResource(ctx context.Context, uri string) {
	authToken := ctx.Value(authToken{}).(string)
	if authToken == "" {
		return
	}

	client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)
	projectID := uri[len("openai-orgs://usage/"):]

	// Get all types of usage for comprehensive data
	usageData := make(map[string]interface{})

	// Get completions usage
	completions, err := client.GetCompletionsUsage(map[string]string{
		"start_time": "0",
		"project_id": projectID,
	})
	if err == nil {
		usageData["completions"] = completions
	}

	// Get embeddings usage
	embeddings, err := client.GetEmbeddingsUsage(map[string]string{
		"start_time": "0",
		"project_id": projectID,
	})
	if err == nil {
		usageData["embeddings"] = embeddings
	}

	// Get images usage
	images, err := client.GetImagesUsage(map[string]string{
		"start_time": "0",
		"project_id": projectID,
	})
	if err == nil {
		usageData["images"] = images
	}

	data, err := json.Marshal(usageData)
	if err != nil {
		return
	}

	contents := &mcp.TextResourceContents{
		URI:      uri,
		MIMEType: MIMETypeUsage,
		Text:     string(data),
	}

	subManager.notify(uri, contents)
}
