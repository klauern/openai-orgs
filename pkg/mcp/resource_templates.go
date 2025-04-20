package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Resource template types and MIME types
const (
	ResourceTemplateTypeProject               = "project"
	ResourceTemplateTypeMember                = "member"
	ResourceTemplateTypeUsage                 = "usage"
	ResourceTemplateTypeProjectServiceAccount = "project-service-account"

	MIMETypeProjectTemplate               = "application/vnd.openai-orgs.project+json"
	MIMETypeMemberTemplate                = "application/vnd.openai-orgs.member+json"
	MIMETypeUsageTemplate                 = "application/vnd.openai-orgs.usage+json"
	MIMETypeProjectServiceAccountTemplate = "application/vnd.openai-orgs.project-service-account+json"
)

// templateHandler is a generic handler for resource templates
type templateHandler func(ctx context.Context, client *openaiorgs.Client, uri *ResourceURI) (interface{}, error)

// AddResourceTemplates adds resource template capabilities to the MCP server
func AddResourceTemplates(s *server.MCPServer) {
	templates := []struct {
		path     string
		name     string
		desc     string
		mimeType string
		handler  templateHandler
	}{
		{
			path:     "openai-orgs://project/{id}",
			name:     "OpenAI Project",
			desc:     "OpenAI project information",
			mimeType: MIMETypeProjectTemplate,
			handler:  handleProject,
		},
		{
			path:     "openai-orgs://members/{id}",
			name:     "Organization Member",
			desc:     "OpenAI organization member information",
			mimeType: MIMETypeMemberTemplate,
			handler:  handleMember,
		},
		{
			path:     "openai-orgs://usage/{id}",
			name:     "Usage Statistics",
			desc:     "OpenAI organization usage statistics",
			mimeType: MIMETypeUsageTemplate,
			handler:  handleUsage,
		},
		{
			path:     "openai-orgs://project/{id}/service-account/{serviceAccountID}",
			name:     "Project Service Account",
			desc:     "OpenAI project service account information",
			mimeType: MIMETypeProjectServiceAccountTemplate,
			handler:  handleProjectServiceAccount,
		},
	}

	for _, t := range templates {
		template := mcp.NewResourceTemplate(
			t.path,
			t.name,
			mcp.WithTemplateDescription(t.desc),
			mcp.WithTemplateMIMEType(t.mimeType),
		)

		// Create a closure to capture the handler and MIME type
		handler := createTemplateHandler(t.handler, t.mimeType)
		s.AddResourceTemplate(template, handler)
	}
}

// createTemplateHandler creates a standard MCP handler from our template handler
func createTemplateHandler(h templateHandler, mimeType string) func(context.Context, mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		authToken := ctx.Value(authToken{}).(string)
		if authToken == "" {
			return nil, ErrNoAuthToken
		}

		uri, err := ParseURI(request.Params.URI)
		if err != nil {
			return nil, fmt.Errorf("invalid URI: %w", err)
		}

		client := openaiorgs.NewClient(openaiorgs.DefaultBaseURL, authToken)

		data, err := h(ctx, client, uri)
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

		return []mcp.ResourceContents{contents}, nil
	}
}

// Individual handlers for each template type
func handleProject(_ context.Context, client *openaiorgs.Client, uri *ResourceURI) (interface{}, error) {
	return client.RetrieveProject(uri.ProjectID)
}

func handleMember(_ context.Context, client *openaiorgs.Client, uri *ResourceURI) (interface{}, error) {
	return client.RetrieveUser(uri.MemberID)
}

func handleProjectServiceAccount(_ context.Context, client *openaiorgs.Client, uri *ResourceURI) (interface{}, error) {
	return client.RetrieveProjectServiceAccount(uri.ProjectID, uri.ServiceAccount)
}

func handleUsage(_ context.Context, client *openaiorgs.Client, uri *ResourceURI) (interface{}, error) {
	usageData := make(map[string]interface{})
	params := map[string]string{
		"start_time": "0",
		"project_id": uri.ProjectID,
	}

	// Get all types of usage for comprehensive data
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
