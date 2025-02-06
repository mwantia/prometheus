package tools

import (
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/mwantia/queueverse/pkg/plugin/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProviderToolHandlerTest struct {
	shared.UnimplementedProviderToolHandler
	Tools []shared.ToolDefinition
}

func init() {
	gob.Register(&ProviderToolHandlerTest{})
}

func NewTest() *ProviderToolHandlerTest {
	return &ProviderToolHandlerTest{
		Tools: []shared.ToolDefinition{
			TimeGetCurrent,
			DiscordListContact,
			DiscordSendPM,
		},
	}
}

func (h *ProviderToolHandlerTest) GetTools() []shared.ToolDefinition {
	return h.Tools
}

func (*ProviderToolHandlerTest) Execute(_ context.Context, function shared.ToolFunction) (string, error) {
	switch function.Name {
	case TimeGetCurrent.Name:
		timezone, exist := function.Arguments["timezone"]
		if !exist {
			return "", fmt.Errorf("failed too call '%s': argument 'timezone' not provided", function.Name)
		}

		location, err := time.LoadLocation(timezone.(string))
		if err != nil {
			return "", fmt.Errorf("failed to load location: ")
		}

		return time.Now().In(location).Format("Mon Jan 2 15:04:05"), nil

	case DiscordListContact.Name:
		return "", status.Error(codes.Unimplemented, "Not implemented")

	case DiscordSendPM.Name:
		return "", status.Error(codes.Unimplemented, "Not implemented")
	}

	return "", nil
}
