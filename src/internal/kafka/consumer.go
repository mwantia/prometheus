package kafka

import (
	"context"
	"fmt"

	"github.com/mwantia/prometheus/pkg/msg"
	"github.com/segmentio/kafka-go"
)

type KafkaMessageHubConsumer struct {
	Topic  string
	Reader *kafka.Reader
}

func (p KafkaMessageHubConsumer) Read(ctx context.Context, handler interface{}) error {
	ev := handlerForInterface(handler)
	if ev == nil {
		return fmt.Errorf("unable to create handler: %v", handler)
	}

	for {
		m, err := p.Reader.ReadMessage(ctx)
		if err != nil {
			break
		}

		ev.Handle(string(m.Value))
	}

	return p.Reader.Close()
}

func handlerForInterface(handler interface{}) msg.EventHandler {
	switch v := handler.(type) {
	case func(string):
		return msg.MessageReadEventHandler(v)
	}
	return nil
}

func (p KafkaMessageHubConsumer) Cleanup(ctx context.Context) error {
	return nil
}
