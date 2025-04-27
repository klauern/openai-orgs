package mcp

//go:generate go tool mockgen -destination=mock_interfaces_test.go -package=mcp github.com/klauern/openai-orgs/pkg/mcp ResourceProvider,ClientProvider,ResourceManager,SubscriptionManager
