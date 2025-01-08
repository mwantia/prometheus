package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
)

func ParseConfig(path string) (*Config, error) {
	config := CreateDefaultConfig()
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return config, nil
	}

	if err := hclsimple.DecodeFile(path, nil, config); err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	return config, nil
}

func (c *Config) GetPluginConfig(name string) *PluginConfig {
	for _, plugin := range c.Plugins {
		if plugin.Name == name {
			return plugin
		}
	}

	return nil
}

func (cfg *Config) GetPluginConfigMap(name string) (map[string]interface{}, error) {
	for _, plugin := range cfg.Plugins {
		if plugin.Name == name {
			result := make(map[string]interface{})

			if plugin.Config != nil {
				attrs, diags := plugin.Config.Body.JustAttributes()
				if diags.HasErrors() {
					return nil, fmt.Errorf("failed to get config attributes: %s", diags.Error())
				}

				for name, attr := range attrs {
					value, diags := attr.Expr.Value(&hcl.EvalContext{})
					if diags.HasErrors() {
						return nil, fmt.Errorf("failed to evaluate '%s': %s", name, diags.Error())
					}

					switch {
					case value.Type() == cty.String:
						result[name] = value.AsString()
					case value.Type() == cty.Number:
						f, _ := value.AsBigFloat().Float64()
						result[name] = f
					case value.Type() == cty.Bool:
						result[name] = value.True()
					default:
						result[name] = value.GoString()
					}
				}
			}

			return result, nil
		}
	}

	return nil, nil
}
