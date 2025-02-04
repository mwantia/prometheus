package old

type ChatRoleType string

type ChatMessageType string

type ToolDefinitionType string

type ToolDataType string

const (
	ChatRoleSystem    ChatRoleType = "system"
	ChatRoleUser      ChatRoleType = "user"
	ChatRoleAssistant ChatRoleType = "assistant"

	ChatMessageText       ChatMessageType = "text"
	ChatMessageToolResult ChatMessageType = "tool_result"
	ChatMessageToolUse    ChatMessageType = "tool_use"
	ChatMessageDocument   ChatMessageType = "document"
	ChatMessageImage      ChatMessageType = "image"

	ToolDefinitionFunction ToolDefinitionType = "function"

	ToolTypeObject  ToolDataType = "object"
	ToolTypeString  ToolDataType = "string"
	ToolTypeBoolean ToolDataType = "boolean"
	ToolTypeInteger ToolDataType = "integer"
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
	ID      string               `json:"id,omitempty"`
	Role    ChatRoleType         `json:"role"`
	Content []ChatMessageContent `json:"content"`
}

type ChatMessageContent struct {
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Type      ChatMessageType `json:"type"`
	Text      string          `json:"text,omitempty"`
	ToolCalls []ToolCall      `json:"toolcalls,omitempty"`
}

type ChatResponse struct {
	Model    string         `json:"model"`
	Messages []ChatMessage  `json:"messages"`
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
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  ToolParameters `json:"parameters"`
}

type ToolParameters struct {
	Type       ToolDataType            `json:"type"`
	Required   []string                `json:"required,omitempty"`
	Properties map[string]ToolProperty `json:"properties,omitempty"`
}

type ToolProperty struct {
	Type        ToolDataType `json:"type"`
	Description string       `json:"description"`
	Enum        []string     `json:"enum,omitempty"`
}

type ToolCall struct {
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Index     int            `json:"index,omitempty"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}
