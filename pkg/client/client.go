package client

import (
	"context"

	"github.com/paphachanok/modelgene/pkg/types"
	"github.com/paphachanok/modelgene/providers/anthropic"
	"github.com/paphachanok/modelgene/providers/ollama"
	"github.com/paphachanok/modelgene/pkg/utils"
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
			return nil, utils.NewError(types.ProviderOllama, "failed to initialize Ollama provider", err)
		}
		providers[types.ProviderOllama] = ollamaProvider
	}

	if cfg.AnthropicConfig != nil {
		anthropicProvider, err := anthropic.NewProvider(cfg.AnthropicConfig)
		if err != nil {
			return nil, utils.NewError(types.ProviderAnthropic, "failed to initialize Anthropic provider", err)
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
		return nil, utils.NewError(provider, "provider is not configured", nil)
	}
	return prov.Chat(ctx, req)
}

