package kafka

import (
	"context"
	"encoding/gob"
	"fmt"

	"github.com/mwantia/prometheus/pkg/msg"
	"github.com/segmentio/kafka-go"
)

type KafkaMessageHub struct {
	Network   string `json:"network,omitempty"`
	Address   string `json:"address,omitempty"`
	Topic     string `json:"topic,omitempty"`
	Partition int    `json:"partition,omitempty"`
}

func (k *KafkaMessageHub) CreateProducer(topic string) (msg.MessageHubProducer, error) {
	topic = k.combineTopic(topic)
	conn, err := kafka.DialLeader(context.Background(), k.Network, k.Address, topic, k.Partition)
	if err != nil {
		return nil, err
	}

	conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})

	return KafkaMessageHubProducer{
		Topic: topic,
		conn:  conn,
	}, nil
}

func (k *KafkaMessageHub) CreateConsumer(topic string) (msg.MessageHubConsumer, error) {
	topic = k.combineTopic(topic)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.Address},
		Topic:    topic,
		MaxBytes: 10e6, // 10MB
	})
	return KafkaMessageHubConsumer{
		Topic:  topic,
		Reader: reader,
	}, nil
}

func (k *KafkaMessageHub) combineTopic(topic string) string {
	return fmt.Sprintf("%s_%s", k.Topic, topic)
}

func init() {
	gob.Register(&KafkaMessageHub{})
	gob.Register(&KafkaMessageHubProducer{})
	gob.Register(&KafkaMessageHubConsumer{})
}
