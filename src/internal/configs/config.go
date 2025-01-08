package configs

import "fmt"

type Config struct {
	Agent    *AgentConfig    `hcl:"agent,block"`
	LogLevel string          `hcl:"log_level,optional"`
	Plugins  []*PluginConfig `hcl:"plugin,block"`
}

func CreateDefaultConfig() *Config {
	return &Config{
		Agent: &AgentConfig{
			Server: &AgentServerConfig{
				Address: ":8080",
			},
			PluginDir:    "./plugins",
			EmbedPlugins: make([]string, 0),
		},
		LogLevel: "INFO",
		Plugins:  make([]*PluginConfig, 0),
	}
}

func (c *Config) ValidateConfig() error {
	if c.Agent == nil {
		return fmt.Errorf("agent block is required")
	}

	if err := c.Agent.ValidateAgentConfig(); err != nil {
		return fmt.Errorf("invalid agent block: %v", err)
	}

	return nil
}
