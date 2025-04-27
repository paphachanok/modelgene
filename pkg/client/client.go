package client

import (
	"context"
	"fmt"

	"github.com/paphachanok/modelgene/pkg/types"
	"github.com/paphachanok/modelgene/providers/anthropic"
	"github.com/paphachanok/modelgene/providers/ollama"
)

type Client struct {
	providers map[types.Provider]types.ProviderClient
}

// NewClient initializes the modelgene client based on available configs
func NewClient(cfg *types.Config) (*Client, error) {
	providers := make(map[types.Provider]types.ProviderClient)

	if cfg.OllamaConfig != nil {
		ollamaProvider, err := ollama.NewProvider(cfg.OllamaConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to init Ollama provider: %w", err)
		}
		providers[types.ProviderOllama] = ollamaProvider
	}

	if cfg.AnthropicConfig != nil {
		anthropicProvider, err := anthropic.NewProvider(cfg.AnthropicConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to init Anthropic provider: %w", err)
		}
		providers[types.ProviderAnthropic] = anthropicProvider
	}

	return &Client{
		providers: providers,
	}, nil
}

func (c *Client) Chat(ctx context.Context, provider types.Provider, req types.APIRequest) (*types.APIResponse, error) {
	prov, ok := c.providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s is not configured", provider)
	}
	return prov.Chat(ctx, req)
}
