package configs

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

type PluginConfig struct {
	Name    string            `hcl:"name,label"`
	Enabled bool              `hcl:"enabled,optional"`
	Config  *PluginConfigBody `hcl:"config,block"`
}

type PluginConfigBody struct {
	Body hcl.Body `hcl:",remain"`
}

func (c *PluginConfig) ValidatePluginConfig() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}
