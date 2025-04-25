package modelgene

import (
	"context"

	"gitpkg/types/types.go"
)

// Provider enumerates the supported LLM providers.
type Provider string

const (
	ProviderOllama    Provider = "ollama"
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
	ProviderVertexAI  Provider = "vertexai"
)

// ProviderClient defines the interface that all provider clients must implement
type ProviderClient interface {
	Chat(ctx context.Context, req APIRequest) (*APIResponse, error)
}
