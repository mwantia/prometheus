package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mwantia/queueverse/pkg/msg"
	"github.com/segmentio/kafka-go"
)

type KafkaMessageHubConsumer struct {
	uuid   [16]byte
	topic  string
	reader *kafka.Reader
}

func (p KafkaMessageHubConsumer) GetName() string {
	return string(p.uuid[:])
}

func (p KafkaMessageHubConsumer) Read(ctx context.Context, handler interface{}) error {
	ev := handlerForInterface(handler)
	if ev == nil {
		return fmt.Errorf("unable to create handler: %v", handler)
	}

	var end bool

	for !end {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		km, err := p.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Failed to read kafka message: %v", err)
			break
		}

		var m msg.Message
		if err := json.Unmarshal(km.Value, &m); err != nil {
			log.Printf("Failed to unmarschal kafka message: %v", err)
			continue
		}

		switch m.Type {
		case "end":
			end = true
		case "msg":
			switch ev.Type() {
			case "message":
				ev.Handle(m)
			case "content":
				ev.Handle(m.Content)
			}
		}
	}

	return p.reader.Close()
}

func handlerForInterface(handler interface{}) msg.EventHandler {
	switch v := handler.(type) {
	case func(string):
		return msg.ContentEventHandler(v)
	case func(msg.Message):
		return msg.MessageEventHandler(v)
	}
	return nil
}

func (p KafkaMessageHubConsumer) Cleanup() error {
	return nil
}
