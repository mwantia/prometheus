package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mwantia/prometheus/pkg/msg"
	"github.com/segmentio/kafka-go"
)

type KafkaMessageHubProducer struct {
	uuid  [16]byte
	topic string
	conn  *kafka.Conn
}

func (p KafkaMessageHubProducer) GetName() string {
	return string(p.uuid[:])
}

func (p KafkaMessageHubProducer) Write(ctx context.Context, msg msg.Message) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	deadline, ok := ctx.Deadline()
	if ok {
		p.conn.SetWriteDeadline(deadline)
	} else {
		p.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	}

	if _, err := p.conn.WriteMessages(kafka.Message{
		Value: value,
		Headers: []kafka.Header{
			{Key: "stream_uuid", Value: p.uuid[:]},
		},
	}); err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return nil
}

func (p KafkaMessageHubProducer) Cleanup() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
