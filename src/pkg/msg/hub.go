package msg

import "context"

type MessageHub interface {
	CreateProducer(string) (MessageHubProducer, error)

	CreateConsumer(string) (MessageHubConsumer, error)

	Cleanup(context.Context) error
}
