package types

// Config wraps all provider configs together for NewClient
type Config struct {
	OllamaConfig    *OllamaConfig
	AnthropicConfig *AnthropicConfig
	// OpenAIConfig    *OpenAIConfig
	// VertexAIConfig  *VertexAIConfig
}

// AnthropicConfig
type AnthropicConfig struct {
	APIKey string
}
