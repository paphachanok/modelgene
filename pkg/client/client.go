package client

import (
	"fmt"

	"github.com/paphachanok/modelgene/pkg/types"
	"github.com/paphachanok/modelgene/providers/ollama"
	"github.com/paphachanok/modelgene/providers/anthropic"
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
