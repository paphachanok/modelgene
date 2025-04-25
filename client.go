// This file should now have `package modelgene` (or `main` if it's the entry point)
package modelgene

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	// Import the new types package
	types "github.com/paphachanok/coder-gene-test/pkg/types"

	// Import the ollama provider implementation
	ollama "github.com/paphachanok/coder-gene-test/providers/ollama"

	// Import the official ollama client library
	ollamaapi "github.com/ollama/ollama/api"
)

// Client holds configuration and initialized clients for different providers.
type Client struct {
	// Use the ProviderClient interface from the types package
	providerClients map[types.Provider]types.ProviderClient
}

// NewClient creates a new modelgene client.
// Use OllamaConfig from the types package
func NewClient(ollamaCfg *types.OllamaConfig) (*Client, error) {
	c := &Client{
		// Use ProviderClient interface from types package
		providerClients: make(map[types.Provider]types.ProviderClient),
	}

	// --- Initialize Ollama Provider Client ---
	if ollamaCfg != nil {
		baseURL := ollamaCfg.BaseURL
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		parsedURL, err := url.Parse(strings.TrimSuffix(baseURL, "/"))
		if err != nil {
			return nil, fmt.Errorf("invalid ollama base url: %w", err)
		}

		var httpClient *http.Client
		if ollamaCfg.HTTPClient != nil {
			httpClient = ollamaCfg.HTTPClient
		} else {
			httpClient = http.DefaultClient
		}

		officialOllamaClient := ollamaapi.NewClient(parsedURL, httpClient)

		// Create and store our Ollama provider implementation
		// Use ProviderOllama constant from types package
		c.providerClients[types.ProviderOllama] = ollama.NewOllamaProviderClient(officialOllamaClient)
	}

	return c, nil
}

// Chat sends a request to the specified provider.
// Use Provider enum, APIRequest, APIResponse from the types package
func (c *Client) Chat(ctx context.Context, provider types.Provider, req types.APIRequest) (*types.APIResponse, error) {
	if len(req.Messages) == 0 && (req.Prompt == nil || *req.Prompt == "") {
		// Use types.Provider in the error message
		return nil, fmt.Errorf("modelgene [%s]: request must contain messages or a prompt", provider)
	}
	if req.Model == "" {
		return nil, fmt.Errorf("modelgene [%s]: request must specify a model", provider)
	}

	clientImpl, found := c.providerClients[provider]
	if !found {
		return nil, fmt.Errorf("modelgene [%s]: provider not configured or supported in client", provider)
	}

	return clientImpl.Chat(ctx, req) // Call the interface method
}
