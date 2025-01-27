package api

const (
	useragent = ""
)

type ClientConfig struct {
	Endpoint string `json:"endpoint"`
	Token    string `json:"token,omitempty"`
	Model    string `json:"model,omitempty"`
}

type TaskConfig struct {
	Endpoint string `json:"endpoint"`
	Token    string `json:"token,omitempty"`
	Model    string `json:"model,omitempty"`
	Style    string `json:"style,omitempty"`
}

type QueueRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model,omitempty"`
	Style  string `json:"style,omitempty"`
}

type QueueResponse struct {
	Task   string      `json:"task"`
	State  string      `json:"state"`
	Pool   string      `json:"pool"`
	Result QueueResult `json:"result,omitempty"`
}

type QueueResult struct {
	Text     string  `json:"text"`
	Model    string  `json:"model,omitempty"`
	Style    string  `json:"style,omitempty"`
	Duration float64 `json:"duration,omitempty"`
}
