package tools

import "github.com/mwantia/queueverse/pkg/plugin/base"

type ToolPlugin interface {
	base.BasePlugin

	GetParameters() (*ToolParameters, error)

	Handle(ctx *ToolContext) error
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
