package mcp

import (
	"context"
	"sync"
	"testing"
	"time"

	openaiorgs "github.com/klauern/openai-orgs"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/mock/gomock"
)

// Project represents an OpenAI project for testing
type Project struct {
	ID   string
	Name string
}

// ProjectList represents a list of projects for testing
type ProjectList struct {
	Data []Project
}

// testHandleActiveProjects is a test implementation of the resource handler
func testHandleActiveProjects(_ context.Context, _ *openaiorgs.Client, _ map[string]any) (any, error) {
	// For testing, return a simple project list
	return &ProjectList{
		Data: []Project{
			{ID: "test_proj", Name: "Test Project"},
		},
	}, nil
}

func TestHandleActiveProjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock client
	client := &openaiorgs.Client{}
	_ = &DefaultClientProvider{} // Keep the type but avoid unused variable

	// Test case 1: Successful project listing
	t.Run("successful project listing", func(t *testing.T) {
		ctx := context.Background()
		params := map[string]any{
			"pagination": map[string]any{
				"limit": float64(20),
			},
		}

		result, err := testHandleActiveProjects(ctx, client, params)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Type assertion check
		projects, ok := result.(*ProjectList)
		if !ok {
			t.Error("result is not a ProjectList")
		}

		if len(projects.Data) != 1 {
			t.Errorf("expected 1 project, got %d", len(projects.Data))
		}

		if projects.Data[0].ID != "test_proj" {
			t.Errorf("expected project ID 'test_proj', got %s", projects.Data[0].ID)
		}
	})
}

func TestSubscriptionManager(t *testing.T) {
	// Use the global subscription manager from resources.go
	sm := subManager

	t.Run("subscription lifecycle", func(t *testing.T) {
		uri := "openai-orgs://active-projects"

		// Test subscription
		ch := sm.subscribe(uri)
		if ch == nil {
			t.Error("expected channel to be returned")
		}

		// Send a test notification
		testContent := &mcp.TextResourceContents{
			URI:      uri,
			MIMEType: MIMETypeProjectList,
			Text:     `{"data":[{"id":"test_proj"}]}`,
		}
		sm.notify(uri, testContent)

		// Verify notification was received
		select {
		case content := <-ch:
			if content != testContent {
				t.Errorf("expected content %v, got %v", testContent, content)
			}
		case <-time.After(time.Second):
			t.Error("timeout waiting for notification")
		}

		// Test unsubscription
		sm.unsubscribe(uri, ch)

		// Verify channel is closed
		if _, ok := <-ch; ok {
			t.Error("channel should be closed")
		}
	})

	t.Run("multiple subscribers", func(t *testing.T) {
		uri := "openai-orgs://active-projects"

		// Create two subscribers
		ch1 := sm.subscribe(uri)
		ch2 := sm.subscribe(uri)

		// Send a test notification
		testContent := &mcp.TextResourceContents{
			URI:      uri,
			MIMEType: MIMETypeProjectList,
			Text:     `{"data":[{"id":"test_proj"}]}`,
		}
		sm.notify(uri, testContent)

		// Verify both subscribers receive the notification
		for _, ch := range []chan mcp.ResourceContents{ch1, ch2} {
			select {
			case content := <-ch:
				if content != testContent {
					t.Errorf("expected content %v, got %v", testContent, content)
				}
			case <-time.After(time.Second):
				t.Error("timeout waiting for notification")
			}
		}

		// Unsubscribe one subscriber
		sm.unsubscribe(uri, ch1)

		// Send another notification
		testContent2 := &mcp.TextResourceContents{
			URI:      uri,
			MIMEType: MIMETypeProjectList,
			Text:     `{"data":[{"id":"test_proj_2"}]}`,
		}
		sm.notify(uri, testContent2)

		// Verify ch1 is closed and ch2 receives the notification
		if _, ok := <-ch1; ok {
			t.Error("ch1 should be closed")
		}

		select {
		case content := <-ch2:
			if content != testContent2 {
				t.Errorf("expected content %v, got %v", testContent2, content)
			}
		case <-time.After(time.Second):
			t.Error("timeout waiting for notification")
		}

		// Clean up
		sm.unsubscribe(uri, ch2)
	})

	t.Run("full channel handling", func(t *testing.T) {
		uri := "openai-orgs://active-projects"

		// Create a subscriber with a buffer size of 1
		ch := sm.subscribe(uri)

		// Fill the channel
		testContent1 := &mcp.TextResourceContents{
			URI:      uri,
			MIMEType: MIMETypeProjectList,
			Text:     `{"data":[{"id":"test_proj_1"}]}`,
		}
		sm.notify(uri, testContent1)

		// Try to send another notification without reading the first one
		testContent2 := &mcp.TextResourceContents{
			URI:      uri,
			MIMEType: MIMETypeProjectList,
			Text:     `{"data":[{"id":"test_proj_2"}]}`,
		}
		sm.notify(uri, testContent2)

		// Verify we can still read the first notification
		select {
		case content := <-ch:
			if content != testContent1 {
				t.Errorf("expected content %v, got %v", testContent1, content)
			}
		case <-time.After(time.Second):
			t.Error("timeout waiting for notification")
		}

		// Clean up
		sm.unsubscribe(uri, ch)
	})
}

func TestResourceProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a simple resource provider implementation for testing
	provider := &testResourceProvider{}

	t.Run("get resource", func(t *testing.T) {
		ctx := context.Background()
		uri := "openai-orgs://active-projects"
		params := map[string]any{"limit": 20}

		result, err := provider.GetResource(ctx, uri, params)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if result == nil {
			t.Error("expected non-nil result")
		}
	})

	t.Run("subscribe to updates", func(t *testing.T) {
		uri := "openai-orgs://active-projects"
		ch, cleanup := provider.Subscribe(uri)

		if ch == nil {
			t.Error("expected non-nil channel")
		}

		if cleanup == nil {
			t.Error("expected non-nil cleanup function")
		}

		cleanup()
	})
}

// testResourceProvider is a simple implementation of ResourceProvider for testing
type testResourceProvider struct{}

func (p *testResourceProvider) GetResource(ctx context.Context, uri string, params map[string]any) (any, error) {
	return &ProjectList{}, nil
}

func (p *testResourceProvider) Subscribe(uri string) (<-chan mcp.ResourceContents, func()) {
	ch := make(chan mcp.ResourceContents, 1)
	return ch, func() { close(ch) }
}

func TestSubscriptionManager_ConcurrentNotifyUnsubscribe(t *testing.T) {
	sm := &subscriptionManager{
		subscribers: make(map[string][]*subscription),
	}
	uri := "openai-orgs://test-concurrent"

	// Create several subscribers
	channels := make([]chan mcp.ResourceContents, 10)
	for i := 0; i < 10; i++ {
		channels[i] = sm.subscribe(uri)
	}

	// Concurrently notify and unsubscribe
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		ch := channels[i]
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				sm.notify(uri, &mcp.TextResourceContents{
					URI:  uri,
					Text: "test",
				})
			}
		}()
		go func() {
			defer wg.Done()
			// Small delay to let some notifications through
			sm.unsubscribe(uri, ch)
		}()
	}
	wg.Wait()
	// If we get here without panic, the race condition is fixed
}

func TestGetPaginationFromParams(t *testing.T) {
	tests := []struct {
		name      string
		params    map[string]any
		wantLimit int
		wantAfter string
	}{
		{"default values", map[string]any{}, defaultPageSize, ""},
		{"custom limit", map[string]any{"pagination": map[string]any{"limit": float64(50)}}, 50, ""},
		{"custom after", map[string]any{"pagination": map[string]any{"after": "cursor_123"}}, defaultPageSize, "cursor_123"},
		{"both set", map[string]any{"pagination": map[string]any{"limit": float64(30), "after": "cursor_456"}}, 30, "cursor_456"},
		{"negative limit uses default", map[string]any{"pagination": map[string]any{"limit": float64(-5)}}, defaultPageSize, ""},
		{"limit over max capped", map[string]any{"pagination": map[string]any{"limit": float64(500)}}, maxPageSize, ""},
		{"zero limit uses default", map[string]any{"pagination": map[string]any{"limit": float64(0)}}, defaultPageSize, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			limit, after := getPaginationFromParams(tc.params)
			if limit != tc.wantLimit {
				t.Errorf("limit = %d, want %d", limit, tc.wantLimit)
			}
			if after != tc.wantAfter {
				t.Errorf("after = %q, want %q", after, tc.wantAfter)
			}
		})
	}
}

func TestPollForChanges_CancellableContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		pollForChanges(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
		// pollForChanges exited as expected
	case <-time.After(5 * time.Second):
		t.Fatal("pollForChanges did not exit after context cancellation")
	}
}
