package essentials

import (
	"fmt"
	"log"
	"time"

	"github.com/mwantia/prometheus/pkg/plugin/tools"
)

type GetCurrentTimeTool struct {
	tools.DefaultUnimplementedService
}

func (t *GetCurrentTimeTool) GetName() (string, error) {
	return "get_current_time", nil
}

func (t *GetCurrentTimeTool) GetParameters() (*tools.ToolParameters, error) {
	return &tools.ToolParameters{
		ReturnType:  "string",
		Description: "Get the current time in the specified timezone",
		Properties: []tools.ToolProperty{
			{
				Name:        "timezone",
				Type:        "string",
				Description: "The timezone to use. Must be a IANA Time Zone",
				Required:    true,
			},
		},
	}, nil
}

func (t *GetCurrentTimeTool) Handle(ctx *tools.ToolContext) error {
	tz, err := ctx.GetString("timezone")
	if err != nil {
		return fmt.Errorf("failed to receive property '%s': %w", "timezone", err)
	}

	location, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("failed to load timezone location: %w", err)
	}

	now := time.Now().In(location)
	log.Printf("Time: %s", now.Format("Mon Jan 2 15:04:05"))

	return nil
}

func (t *GetCurrentTimeTool) Probe() error {
	return nil
}
