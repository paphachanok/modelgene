package types

// Config wraps all provider configs together for NewClient
type Config struct {
	OllamaConfig    *OllamaConfig
	// OpenAIConfig    *OpenAIConfig
	// AnthropicConfig *AnthropicConfig
	// VertexAIConfig  *VertexAIConfig
}
