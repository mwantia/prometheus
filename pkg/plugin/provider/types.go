package provider

type ProviderChatRequest struct {
	Model    string                `json:"model"`
	Messages []ProviderChatMessage `json:"messages"`
	Metadata map[string]any        `json:"metadata,omitempty"`
}

type ProviderChatMessage struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ProviderChatResponse struct {
	Model    string              `json:"model"`
	Message  ProviderChatMessage `json:"message"`
	Metadata map[string]any      `json:"metadata,omitempty"`
}
