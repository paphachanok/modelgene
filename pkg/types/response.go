package types

// APIResponse represents a unified response structure from LLM APIs.
type APIResponse struct {
	ID       string   `json:"id,omitempty"`    // Completion ID, if available
	Model    string   `json:"model"`           // Model used for the response
	Provider Provider `json:"provider"`        // Provider that generated the response
	Choices  []Choice `json:"choices"`         // List of response choices
	Usage    *Usage   `json:"usage,omitempty"` // Token usage information, if available
}

// Choice represents a single completion choice within an APIResponse.
// It uses the Message type defined in this package.
type Choice struct {
	Index        int          `json:"index"`
	Message      Message      `json:"message"`            // The generated message (Role="assistant")
	FinishReason string       `json:"finish_reason"`      // e.g., "stop", "length", "tool_calls", "content_filter"
	LogProbs     *LogProbInfo `json:"logprobs,omitempty"` // Log probability info, if requested
}
