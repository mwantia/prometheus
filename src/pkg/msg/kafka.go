package msg

import (
	"context"
	"encoding/gob"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaMessageHub struct {
	Network   string
	Address   string
	Topic     string
	Partition int

	conn *kafka.Conn
}

func (k *KafkaMessageHub) Setup() error {
	conn, err := kafka.DialLeader(context.Background(), k.Network, k.Address, k.Topic, k.Partition)
	if err != nil {
		return err
	}

	k.conn = conn
	return conn.CreateTopics(kafka.TopicConfig{
		Topic:             k.Topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
}

func (k *KafkaMessageHub) WriteMessage(msg string) error {
	k.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := k.conn.WriteMessages(kafka.Message{
		Value: []byte(msg),
	})

	return err
}

func (k *KafkaMessageHub) Cleanup() error {
	if k.conn != nil {
		return k.conn.Close()
	}
	return nil
}

func init() {
	gob.Register(&KafkaMessageHub{})
}
