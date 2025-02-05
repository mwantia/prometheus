package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
)

func ParseConfig(path string) (*Config, error) {
	config := CreateDefault()
	if path == "" {
		return config, nil
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Printf("Unable to access config file '%s': %v", path, err)
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

func (cfg *Config) GetPluginConfigMap(name string) map[string]interface{} {
	pcm := make(map[string]interface{}, 0)

	for _, plugin := range cfg.Plugins {
		if plugin.Name == name {

			if plugin.Config != nil {
				attrs, diags := plugin.Config.Body.JustAttributes()
				if diags.HasErrors() {
					return pcm
				}

				for name, attr := range attrs {
					value, diags := attr.Expr.Value(&hcl.EvalContext{})
					if diags.HasErrors() {
						return pcm
					}

					switch {
					case value.Type() == cty.String:
						pcm[name] = value.AsString()
					case value.Type() == cty.Number:
						f, _ := value.AsBigFloat().Float64()
						pcm[name] = f
					case value.Type() == cty.Bool:
						pcm[name] = value.True()
					default:
						pcm[name] = value.GoString()
					}
				}
			}

			return pcm
		}
	}

	return pcm
}
