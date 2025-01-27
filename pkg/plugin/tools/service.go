package tools

type ToolService interface {
	GetName() (string, error)

	GetParameters() (*ToolParameters, error)

	Handle(ctx *ToolContext) error

	Probe() error
}

type ToolParameters struct {
	ReturnType  string         `json:"return_type"`
	Description string         `json:"description"`
	Properties  []ToolProperty `json:"properties"`
}

type ToolProperty struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
	Required    bool     `json:"required,omitempty"`
}
