package deepseek

// ChatCompletionRequest is the request payload for /chat/completion API
type ChatCompletionRequest struct {
	Messages         []Message      `json:"messages"`
	Model            string         `json:"model"`
	FrequencyPenalty float64        `json:"frequency_penalty"`
	MaxTokens        int            `json:"max_tokens"`
	PresencePenalty  float64        `json:"presence_penalty"`
	ResponseFormat   ResponseFormat `json:"response_format"`
	Stop             interface{}    `json:"stop"` // Can be null, string, or []string
	Stream           bool           `json:"stream"`
	StreamOptions    interface{}    `json:"stream_options"` // Can be null or a struct
	Temperature      float64        `json:"temperature"`
	TopP             float64        `json:"top_p"`
	Tools            interface{}    `json:"tools"` // Can be null or a slice of tools
	ToolChoice       string         `json:"tool_choice"`
	Logprobs         bool           `json:"logprobs"`
	TopLogprobs      interface{}    `json:"top_logprobs"` // Can be null or an integer
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}
