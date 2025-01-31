package tasks

const TaskTypeGenerateName = "task:generate"

type GenerateRequest struct {
	Content  string         `json:"content"`
	Model    string         `json:"model"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type GenerateResponse struct {
	Task   string         `json:"task"`
	State  string         `json:"state"`
	Pool   string         `json:"pool"`
	Result GenerateResult `json:"result,omitempty"`
}

type GenerateResult struct {
	Content  string         `json:"content"`
	Model    string         `json:"model"`
	Metadata map[string]any `json:"metadata,omitempty"`
}
