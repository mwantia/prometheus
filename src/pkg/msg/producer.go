package msg

import "context"

type MessageHubProducer interface {
	GetName() string

	Write(context.Context, string) error

	Cleanup(context.Context) error
}
