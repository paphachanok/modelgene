// This file must have `package ollama`
package ollama

import (
	"context"
	"fmt"

	// Import the new types package
	types "github.com/paphachanok/coder-gene-test/pkg/types"

	// Import the official ollama client library types
	ollamaapi "github.com/ollama/ollama/api"
)

// OllamaProviderClient implements the types.ProviderClient interface for Ollama.
type OllamaProviderClient struct {
	client *ollamaapi.Client
}

// NewOllamaProviderClient creates a new Ollama provider client.
func NewOllamaProviderClient(officialClient *ollamaapi.Client) *OllamaProviderClient {
	return &OllamaProviderClient{
		client: officialClient,
	}
}

// Chat handles non-streaming chat requests to the Ollama API.
// Use types from the imported types package (types.APIRequest, types.APIResponse)
func (c *OllamaProviderClient) Chat(ctx context.Context, req types.APIRequest) (*types.APIResponse, error) {
	ollamaReq, err := convertToOllamaRequest(req) // convertToOllamaRequest is in this package
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to convert request: %w", err)
	}
	ollamaReq.Stream = new(bool) // Ensure non-streaming

	var finalResp ollamaapi.ChatResponse
	err = c.client.Chat(ctx, ollamaReq, func(resp ollamaapi.ChatResponse) error {
		finalResp = resp
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("ollama: chat api call failed: %w", err)
	}

	// Pass the ProviderOllama constant from the types package
	apiResp, err := convertFromOllamaResponse(finalResp, types.ProviderOllama) // convertFromOllamaResponse is in this package
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to convert response: %w", err)
	}

	return apiResp, nil
}
