package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mwantia/queueverse/internal/config"
	"github.com/mwantia/queueverse/pkg/plugin/base"
	"github.com/mwantia/queueverse/pkg/plugin/provider"
	"github.com/mwantia/queueverse/plugins/ollama/tools"
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

	request := provider.ChatRequest{
		Model: "llama3.2:latest",
		Message: provider.Message{
			Content: "Send a message to Roman Blake over Discord and tell him that I might arrive late to the meeting.",
		},
		Tools: []provider.ToolDefinition{
			tools.TimeGetCurrent,
			tools.DiscordListContact,
			tools.DiscordSendPM,
		},
	}

	resp, err := plugin.Chat(request)
	if err != nil {
		t.Fatalf("Failed to perform chat request: %v", err)
	}

	debug, _ := json.Marshal(resp)
	log.Println(string(debug))
}

func executeToolCall(function provider.ToolFunction) (string, error) {
	switch function.Name {
	case tools.TimeGetCurrent.Name:
		timezone, exist := function.Arguments["timezone"]
		if !exist {
			return "", fmt.Errorf("failed too call '%s': argument 'timezone' not provided", function.Name)
		}

		location, err := time.LoadLocation(timezone.(string))
		if err != nil {
			return "", fmt.Errorf("failed to load location: ")
		}

		return time.Now().In(location).Format("Mon Jan 2 15:04:05"), nil

	case tools.DiscordListContact.Name:

	case tools.DiscordSendPM.Name:
	}
	return "", nil
}
