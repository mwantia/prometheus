package provider

type DataType string

const (
	TypeObject  DataType = "object"
	TypeString  DataType = "string"
	TypeBoolean DataType = "boolean"
	TypeInteger DataType = "integer"
)

type Model struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type Message struct {
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string           `json:"model"`
	Message  Message          `json:"message"`
	Tools    []ToolDefinition `json:"tools,omitempty"`
	Metadata map[string]any   `json:"metadata,omitempty"`
}

type ChatResponse struct {
	Model    string         `json:"model"`
	Message  Message        `json:"message"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type EmbedRequest struct {
	Model    string         `json:"model"`
	Message  Message        `json:"message"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type EmbedResponse struct {
	Model      string         `json:"model"`
	Embeddings [][]float32    `json:"embeddings"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type ToolDefinition struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Type        DataType                `json:"type"`
	Required    []string                `json:"required,omitempty"`
	Properties  map[string]ToolProperty `json:"properties,omitempty"`
}

type ToolProperty struct {
	Type        DataType `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type ToolFunction struct {
	Index     int            `json:"index,omitempty"`
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}
