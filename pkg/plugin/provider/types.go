package provider

const (
	ChatRoleSystem    = "system"
	ChatRoleUser      = "user"
	ChatRoleAssistant = "assistant"
	ChatRoleTool      = "tool"
)

type Model struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type ChatRequest struct {
	Model    string           `json:"model"`
	Messages []ChatMessage    `json:"messages"`
	Tools    []ToolDefinition `json:"tools,omitempty"`
	Metadata map[string]any   `json:"metadata,omitempty"`
}

type ChatMessage struct {
	ID        string     `json:"id,omitempty"`
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"toolcalls,omitempty"`
}

type ChatResponse struct {
	Model    string         `json:"model"`
	Message  ChatMessage    `json:"message"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type EmbedRequest struct {
	Model    string         `json:"model"`
	Input    string         `json:"input"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type EmbedResponse struct {
	Model      string         `json:"model"`
	Embeddings [][]float32    `json:"embeddings"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type ToolDefinition struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Parameter   ToolParameter `json:"parameter"`
}

type ToolParameter struct {
	Type       string         `json:"type"`
	Required   []string       `json:"required,omitempty"`
	Properties []ToolProperty `json:"properties,omitempty"`
}

type ToolProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type ToolCall struct {
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Index     int            `json:"index,omitempty"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}
