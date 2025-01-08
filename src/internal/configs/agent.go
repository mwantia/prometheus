package configs

import "fmt"

type AgentConfig struct {
	Server *AgentServerConfig `hcl:"server,block"`
	Kafka  *AgentKafkaConfig  `hcl:"kafka,block"`

	PluginDir    string   `hcl:"plugin_dir,optional"`
	EmbedPlugins []string `hcl:"embed_plugins,optional"`
}

type AgentServerConfig struct {
	Address string `hcl:"address,optional"`
}

type AgentKafkaConfig struct {
	Network   string `hcl:"network,optional"`
	Address   string `hcl:"address,optional"`
	Topics    string `hcl:"topics,optional"`
	Partition int    `hcl:"partition,optional"`
}

func (c *AgentConfig) ValidateAgentConfig() error {
	if c.Server == nil {
		return fmt.Errorf("server block is required")
	}

	if c.Server.Address == "" {
		return fmt.Errorf("server address is required")
	}

	return nil
}
