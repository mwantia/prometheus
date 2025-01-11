package msg

import "context"

type MessageHubProducer interface {
	GetName() string

	Write(context.Context, Message) error

	Cleanup() error
}
