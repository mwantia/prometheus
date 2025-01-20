package config

type KafkaConfig struct {
	Network   string `hcl:"network,optional"`
	Address   string `hcl:"address,optional"`
	Topics    string `hcl:"topics,optional"`
	Partition int    `hcl:"partition,optional"`
}

func (c *KafkaConfig) ValidateConfig() error {
	return nil
}
