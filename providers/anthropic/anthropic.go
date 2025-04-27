package anthropic

import (
	"context"
	"errors"
	"fmt"

	"github.com/paphachanok/modelgene/pkg/types"
	"github.com/anthropics/anthropic-sdk-go"
)

// Provider implements types.ProviderClient
type Provider struct {
	client *AnthropicClient
	config *types.AnthropicConfig
}

// NewProvider initializes an Anthropic Provider
func NewProvider(cfg *types.AnthropicConfig) (*Provider, error) {
	if cfg == nil {
		return nil, errors.New("anthropic config is nil")
	}
	cli := NewAnthropicClient(cfg.APIKey)
	return &Provider{
		client: cli,
		config: cfg,
	}, nil
}

// Chat implements types.ProviderClient
func (p *Provider) Chat(ctx context.Context, req types.APIRequest) (*types.APIResponse, error) {
	if req.Model == "" {
		return nil, errors.New("model name is required")
	}

	// Convert types.Message to anthropic.MessageParam
	messages := make([]anthropic.MessageParam, 0, len(req.Messages))
	for _, m := range req.Messages {
		messages = append(messages, anthropic.MessageParam{
			Role: anthropic.MessageParamRole(m.Role),
			Content: []anthropic.ContentBlockParamUnion{
				{
					OfRequestTextBlock: &anthropic.TextBlockParam{
						Text: m.Content,
					},
				},
			},
		})
	}

	params := anthropic.MessageNewParams{
		Model:     req.Model,
		Messages:  messages,
		MaxTokens: getMaxTokens(req.MaxTokens),
	}

	resp, err := p.client.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("anthropic chat error: %w", err)
	}

	response := &types.APIResponse{
		Model:    req.Model,
		Provider: types.ProviderAnthropic,
		Choices: []types.Choice{
			{
				Index: 0,
				Message: types.Message{
					Role:    "assistant",
					Content: mergeContent(resp.Content),
				},
				FinishReason: "stop",
			},
		},
	}

	return response, nil
}

// --- Helper functions ---

func getMaxTokens(max *int) int64 {
	if max != nil {
		return int64(*max)
	}
	return 1024
}

func mergeContent(blocks []anthropic.ContentBlockUnion) string {
	var combined string
	for _, block := range blocks {
		if block.Text != "" {
			combined += block.Text
		}
	}
	return combined
}
