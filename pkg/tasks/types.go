package tasks

import "github.com/mwantia/queueverse/pkg/plugin/provider"

const TaskTypeGenerateName = "task:generate"

type GenerateRequest struct {
	Input provider.ChatRequest `json:"input"`
}

type GenerateResponse struct {
	Task   string                `json:"task"`
	State  string                `json:"state"`
	Pool   string                `json:"pool"`
	Output provider.ChatResponse `json:"output,omitempty"`
}
