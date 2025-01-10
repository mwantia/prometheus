package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaMessageHubProducer struct {
	topic string
	conn  *kafka.Conn
}

func (p KafkaMessageHubProducer) GetName() string {
	return p.topic
}

func (p KafkaMessageHubProducer) Write(ctx context.Context, content string) error {
	p.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	c, err := p.conn.WriteMessages(kafka.Message{
		Value: []byte(content),
	})
	log.Printf("Bytes written to Kafka: %v", c)

	return err
}

func (p KafkaMessageHubProducer) Cleanup(ctx context.Context) error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
