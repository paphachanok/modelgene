package types

import "encoding/json"

// Message represents a single message in the conversation history across different APIs.
type Message struct {
	// Role specifies the author of the message.
	// Common values: "system", "user", "assistant".
	// Google Vertex AI uses "model" instead of "assistant".
	// OpenAI uses "tool".
	// Anthropic uses only "user", "assistant".
	Role string `json:"role"` // Consider standardizing to "system", "user", "assistant", "tool" and mapping internally

	// Content is the textual content of the message.
	Content string `json:"content"`

	// Images contains base64 encoded image data (primarily for Ollama multimodal).
	// Other providers might use different structures (e.g., OpenAI content array, Anthropic content blocks, Vertex parts).
	// Consider using a more generic `Parts` field or handling multimodal separately.
	Images []string `json:"images,omitempty"` // Used by Ollama

	// Name is used by OpenAI for tool/function messages.
	Name *string `json:"name,omitempty"` // Used by OpenAI

	// ToolCallID is used by OpenAI for tool response messages.
	ToolCallID *string `json:"tool_call_id,omitempty"` // Used by OpenAI

	// ToolCalls are generated by the assistant in OpenAI when invoking tools.
	// This might be part of the *response* structure but needs to be included
	// in subsequent *requests* if the assistant calls a tool.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"` // Used by OpenAI

	// --- Fields potentially needed for other providers' multimodal/tool structures ---
	// Parts could be a more generic way to handle multimodal content blocks (text, image, tool_use, tool_result).
	// Parts []ContentPart `json:"parts,omitempty"` // Adapt for Anthropic/Vertex/OpenAI vision
}

// ToolCall represents a tool invocation requested by the model (mainly OpenAI).
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"` // Usually "function"
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction details the function called by the model.
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string arguments
}

// Tool represents a tool/function definition provided to the model.
type Tool struct {
	// Type is primarily used by OpenAI ("function").
	Type string `json:"type,omitempty"` // Used by OpenAI

	// Function holds the definition for OpenAI tools.
	Function *FunctionDefinition `json:"function,omitempty"` // Used by OpenAI

	// --- Fields for Anthropic/Vertex AI Tools ---
	Name        *string         `json:"name,omitempty"`         // Used by Anthropic/Vertex AI (and OpenAI within Function)
	Description *string         `json:"description,omitempty"`  // Used by Anthropic/Vertex AI (and OpenAI within Function)
	InputSchema json.RawMessage `json:"input_schema,omitempty"` // Used by Anthropic (JSON Schema), Vertex AI uses a specific proto struct
}

// FunctionDefinition describes a function for OpenAI tools.
type FunctionDefinition struct {
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters"` // JSON Schema object
}

// ToolChoice controls how the model uses tools.
type ToolChoice struct {
	// Type controls the mode: "none", "auto", "any" (Anthropic), "tool" (Anthropic), "function" (OpenAI).
	Type string `json:"type"`

	// Name specifies a particular tool/function to call. Required by Anthropic if type="tool",
	// used within an object by OpenAI if specifying a function.
	Name *string `json:"name,omitempty"` // Used by Anthropic directly, OpenAI uses it within a nested structure

	// Function can be used for OpenAI's specific function choice structure.
	Function *ToolChoiceFunction `json:"function,omitempty"` // Used by OpenAI: { "type": "function", "function": { "name": "my_func" } }
}

// ToolChoiceFunction specifies a function name for OpenAI's tool_choice.
type ToolChoiceFunction struct {
	Name string `json:"name"`
}

// ResponseFormat specifies constraints on the output format (e.g., JSON).
type ResponseFormat struct {
	// Type is used by OpenAI (e.g., "text", "json_object").
	Type string `json:"type,omitempty"` // Used by OpenAI

	// Format is used by Ollama (e.g., "json").
	Format string `json:"format,omitempty"` // Used by Ollama
}

// SafetySetting configures content safety filters (primarily Google Vertex AI).
type SafetySetting struct {
	Category  string `json:"category"`  // e.g., "HARM_CATEGORY_SEXUALLY_EXPLICIT"
	Threshold string `json:"threshold"` // e.g., "BLOCK_LOW_AND_ABOVE"
}

// APIRequest represents the consolidated parameters for calling various LLM APIs.
type APIRequest struct {
	// --- Core Identification & Input ---
	Model    string    `json:"model"`              // Required by all
	Messages []Message `json:"messages,omitempty"` // Required by OpenAI, Anthropic, Ollama (/chat), Vertex AI (within instances)
	Prompt   *string   `json:"prompt,omitempty"`   // Used by Ollama (/generate)

	// --- System/Context Prompts ---
	System  *string `json:"system,omitempty"`  // Used by Ollama, Anthropic, OpenAI (as system message)
	Context *string `json:"context,omitempty"` // Used by Vertex AI (within instances)

	// --- Generation Control Parameters ---
	MaxTokens        *int     `json:"max_tokens,omitempty"`        // OpenAI, Anthropic (required), Vertex AI (`maxOutputTokens`), Ollama (`num_predict`)
	Temperature      *float64 `json:"temperature,omitempty"`       // All providers
	TopP             *float64 `json:"top_p,omitempty"`             // All providers
	TopK             *int     `json:"top_k,omitempty"`             // All providers (except perhaps base Ollama options)
	StopSequences    []string `json:"stop_sequences,omitempty"`    // All providers (named `stop` in OpenAI/Ollama, `stopSequences` in Anthropic/Vertex)
	PresencePenalty  *float64 `json:"presence_penalty,omitempty"`  // OpenAI, Vertex AI, Ollama
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"` // OpenAI, Vertex AI, Ollama
	Seed             *int     `json:"seed,omitempty"`              // OpenAI, Ollama
	Stream           *bool    `json:"stream,omitempty"`            // All providers

	// --- Response & Output Formatting ---
	N              *int            `json:"n,omitempty"`               // OpenAI (`n`), Vertex AI (`candidateCount`)
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"` // OpenAI (`response_format`), Ollama (`format`)
	LogProbs       *bool           `json:"logprobs,omitempty"`        // OpenAI (bool), Vertex AI (int -> number of tokens) - *Type difference!*
	TopLogProbs    *int            `json:"top_logprobs,omitempty"`    // OpenAI

	// --- Tool/Function Calling ---
	Tools      []Tool      `json:"tools,omitempty"`       // OpenAI, Anthropic, Vertex AI
	ToolChoice *ToolChoice `json:"tool_choice,omitempty"` // OpenAI, Anthropic

	// --- Provider-Specific Parameters ---

	// OpenAI Specific
	LogitBias map[string]int `json:"logit_bias,omitempty"` // OpenAI
	User      *string        `json:"user,omitempty"`       // OpenAI

	// Anthropic Specific
	Metadata map[string]interface{} `json:"metadata,omitempty"` // Anthropic

	// Vertex AI Specific
	SafetySettings []SafetySetting `json:"safety_settings,omitempty"` // Vertex AI (`safetySettings` within parameters)
	Echo           *bool           `json:"echo,omitempty"`            // Vertex AI (`echo` within parameters)
	// Vertex AI 'instances' structure needs specific handling outside this unified struct.
	// Vertex AI 'logprobs' is an integer, not bool. Need to handle this conversion.

	// Ollama Specific
	Template  *string                `json:"template,omitempty"`   // Ollama (overrides messages/prompt/system)
	Raw       *bool                  `json:"raw,omitempty"`        // Ollama
	KeepAlive *string                `json:"keep_alive,omitempty"` // Ollama (duration string or seconds)
	OllamaOptions    map[string]interface{} `json:"ollama_options,omitempty"`    // Ollama (catch-all for model-specific options not covered above)
	// Ollama `context` (stateful int array) is intentionally omitted here as it requires state management
	// between calls, which is usually handled outside the request parameters struct itself.
}
