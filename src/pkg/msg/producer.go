package msg

import "context"

type MessageHubProducer interface {
	Write(context.Context, string) error

	Cleanup(context.Context) error
}
