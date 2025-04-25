package ollama

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"

	ollamaapi "github.com/ollama/ollama/api"
	modelgene "github.com/paphachanok/coder-gene-test"
)

// Client implements the Ollama provider client
type Client struct {
	client *ollamaapi.Client
}

// NewClient creates a new Ollama client
func NewClient(baseURL *url.URL, httpClient *http.Client) modelgene.ProviderClient {
	return &Client{
		client: ollamaapi.NewClient(baseURL, httpClient),
	}
}

// Chat implements the ProviderClient interface for Ollama
func (c *Client) Chat(ctx context.Context, req modelgene.APIRequest) (*modelgene.APIResponse, error) {
	// Convert the request to Ollama format
	ollamaReq, err := convertToOllamaRequest(req)
	if err != nil {
		return nil, modelgene.NewError(modelgene.ProviderOllama, "failed to convert request", err)
	}

	// Initialize response
	var finalResp ollamaapi.ChatResponse
	var mu sync.Mutex

	// Make the API call
	err = c.client.Chat(ctx, ollamaReq, func(resp ollamaapi.ChatResponse) error {
		// For streaming responses, this callback is called multiple times
		// For non-streaming, it's called once with the complete response
		mu.Lock()
		defer mu.Unlock()

		// If we're streaming, accumulate content
		if ollamaReq.Stream != nil && *ollamaReq.Stream {
			// For streaming, accumulate the content
			if finalResp.Message.Content == "" {
				finalResp = resp
			} else {
				finalResp.Message.Content += resp.Message.Content
			}

			// Update metadata from the latest chunk
			finalResp.Model = resp.Model
			finalResp.CreatedAt = resp.CreatedAt
			finalResp.Done = resp.Done
			finalResp.TotalDuration = resp.TotalDuration
			finalResp.LoadDuration = resp.LoadDuration
			finalResp.PromptEvalCount = resp.PromptEvalCount
			finalResp.PromptEvalDuration = resp.PromptEvalDuration
			finalResp.EvalCount = resp.EvalCount
			finalResp.EvalDuration = resp.EvalDuration
		} else {
			// For non-streaming, just store the complete response
			finalResp = resp
		}

		// Check if this is the final chunk for streaming response
		if resp.Done {
			return errors.New("done") // Signal completion to stop streaming
		}

		return nil
	})

	// Handle errors properly
	if err != nil && err.Error() != "done" {
		return nil, modelgene.NewError(modelgene.ProviderOllama, "ollama chat api call failed", err)
	}

	// Convert to standard response
	return convertFromOllamaResponse(finalResp, modelgene.ProviderOllama)
}
