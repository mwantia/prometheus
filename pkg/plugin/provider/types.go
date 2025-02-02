package provider

const (
	ChatRoleSystem    = "system"
	ChatRoleUser      = "user"
	ChatRoleAssistant = "assistant"
	ChatRoleTool      = "tool"
	ChatRoleDocument  = "document"
	ChatRoleImage     = "image"
)

type ToolDefinitionType string

const (
	ToolDefinitionFunction ToolDefinitionType = "function"
)

type ToolType string

const (
	ToolTypeObject  ToolType = "object"
	ToolTypeString  ToolType = "string"
	ToolTypeBoolean ToolType = "boolean"
	ToolTypeInteger ToolType = "integer"
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
	Type       ToolType                `json:"type"`
	Required   []string                `json:"required,omitempty"`
	Properties map[string]ToolProperty `json:"properties,omitempty"`
}

type ToolProperty struct {
	Type        ToolType `json:"type"`
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
