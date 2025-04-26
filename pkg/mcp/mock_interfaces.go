package mcp

//go:generate mockgen -destination=mock_interfaces_test.go -package=mcp github.com/klauer/openai-orgs/pkg/mcp ResourceProvider,ClientProvider,ResourceManager,SubscriptionManager
