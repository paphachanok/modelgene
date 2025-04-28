package anthropic

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/paphachanok/modelgene/pkg/types"
	"github.com/paphachanok/modelgene/pkg/utils"
)

// Provider implements types.ProviderClient
type Provider struct {
	client *AnthropicClient
	config *types.AnthropicConfig
}

// NewProvider initializes an Anthropic Provider
func NewProvider(cfg *types.AnthropicConfig) (*Provider, error) {
	if cfg == nil {
		return nil, utils.NewError(types.ProviderAnthropic, "anthropic config is nil", nil)
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
		return nil, utils.NewError(types.ProviderAnthropic, "model name is required", nil)

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
		return nil, utils.NewError(types.ProviderAnthropic, "anthropic chat error", err)
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
