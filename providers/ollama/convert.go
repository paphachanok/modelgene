// Ensure this file starts with: package ollama
package ollama

import (

	// Corrected import path for your root package types
	modelgene "github.com/paphachanok/coder-gene-test"

	// Import the official ollama client library types
	ollamaapi "github.com/ollama/ollama/api"
)

// convertToOllamaRequest converts a generic APIRequest to an Ollama ChatRequest
// Uses types from the imported modelgene package
func convertToOllamaRequest(req modelgene.APIRequest) (*ollamaapi.ChatRequest, error) {
	ollamaReq := &ollamaapi.ChatRequest{
		Model:    req.Model,
		Messages: make([]ollamaapi.Message, len(req.Messages)),
		Options:  make(map[string]any),
		Stream:   new(bool), // Default to non-streaming
	}

	// Convert messages
	for i, msg := range req.Messages {
		ollamaReq.Messages[i] = ollamaapi.Message{
			Role:    msg.Role,
			Content: msg.Content,
			Images:  ollamaapi.ImageData(msg.Images),
		}
	}

	// --- Map Options ---
	setOptionIfPresent := func(key string, value interface{}) {
		// (Keep the helper function as provided in the previous response)
		switch v := value.(type) {
		case *string:
			if v != nil {
				ollamaReq.Options[key] = *v
			}
		case *int:
			if v != nil {
				ollamaReq.Options[key] = *v
			}
		case *float64:
			if v != nil {
				ollamaReq.Options[key] = *v
			}
		case *bool:
			if v != nil {
				ollamaReq.Options[key] = *v
			}
		case []string:
			if len(v) > 0 {
				ollamaReq.Options[key] = v
			}
		}
	}


	setOptionIfPresent("temperature", req.Temperature)
	setOptionIfPresent("top_p", req.TopP)
	setOptionIfPresent("top_k", req.TopK)
	setOptionIfPresent("num_predict", req.MaxTokens)
	setOptionIfPresent("presence_penalty", req.PresencePenalty)
	setOptionIfPresent("frequency_penalty", req.FrequencyPenalty)
	setOptionIfPresent("seed", req.Seed)
	setOptionIfPresent("stop", req.StopSequences)

	// --- Map Direct Fields ---
	if req.System != nil {
		ollamaReq.System = *req.System
	}
	if req.Template != nil {
		ollamaReq.Template = *req.Template
	}
	if req.ResponseFormat != nil && req.ResponseFormat.Format == "json" {
		ollamaReq.Format = ollamaapi.FormatJSON
	}
	if req.KeepAlive != nil {
		// Adjust conversion based on ollamaapi.Duration type if needed
		ka := ollamaapi.Duration(*req.KeepAlive)
		ollamaReq.KeepAlive = &ka
	}
	if req.OllamaOptions != nil {
		for k, v := range req.OllamaOptions {
			if _, exists := ollamaReq.Options[k]; !exists {
				ollamaReq.Options[k] = v
			}
		}
	}
	// Add other options like Raw if needed: setOptionIfPresent("raw", req.Raw)
	return ollamaReq, nil
}

// convertFromOllamaResponse converts an Ollama ChatResponse to a generic APIResponse
// Uses types from the imported modelgene package
func convertFromOllamaResponse(resp ollamaapi.ChatResponse, provider modelgene.Provider) (*modelgene.APIResponse, error) {
	apiResp := &modelgene.APIResponse{
		Model:    resp.Model,
		Provider: provider, // Use the provider constant from modelgene package
		Choices: []modelgene.Choice{ // Choice struct from modelgene package
			{
				Index: 0,
				Message: modelgene.Message{ // Message struct from modelgene package
					Role:    resp.Message.Role,
					Content: resp.Message.Content,
				},
				FinishReason: "stop",
			},
		},
	}

	promptTokens := resp.PromptEvalCount
	completionTokens := resp.EvalCount
	totalTokens := promptTokens + completionTokens

	if promptTokens > 0 || completionTokens > 0 {
		apiResp.Usage = &modelgene.Usage{ // Usage struct from modelgene package
			PromptTokens:     &promptTokens,
			CompletionTokens: &completionTokens,
			TotalTokens:      &totalTokens,
		}
	}

	return apiResp, nil
}
