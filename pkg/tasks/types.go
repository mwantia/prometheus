package tasks

import "github.com/mwantia/queueverse/pkg/plugin/shared"

const TaskTypeGenerateName = "task:generate"

type GenerateRequest struct {
	Input shared.ChatRequest `json:"input"`
}

type GenerateResponse struct {
	Task   string              `json:"task"`
	State  string              `json:"state"`
	Pool   string              `json:"pool"`
	Output shared.ChatResponse `json:"output,omitempty"`
}
