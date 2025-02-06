package ollama

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/internal/tools"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/shared"
)

const (
	ConfigPath = "../../tests/config.hcl"
)

func TestOllamaProvider(t *testing.T) {
	cfg, err := config.ParseConfig(ConfigPath)
	if err != nil {
		t.Fatalf("Failed to parse test config: %v", err)
	}

	pcm := cfg.GetPluginConfigMap(PluginName)

	plugin := OllamaProvider{
		Context: context.TODO(),
	}
	if err := plugin.SetConfig(&base.PluginConfig{ConfigMap: pcm}); err != nil {
		t.Fatalf("Failed to set plugin config: %v", err)
	}

	request := shared.ChatRequest{
		Model: "llama3.1:8b",
		Message: shared.Message{
			Content: "Tell me the current time in germany.",
		},
	}

	resp, err := plugin.Chat(request, tools.NewTest())
	if err != nil {
		t.Fatalf("Failed to perform chat request: %v", err)
	}

	debug, _ := json.Marshal(resp)
	log.Println(string(debug))
}
