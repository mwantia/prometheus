package msg

import "context"

type MessageHubConsumer interface {
	GetName() string

	Read(context.Context, interface{}) error

	Cleanup(context.Context) error
}

type EventHandler interface {
	Type() string

	Handle(interface{})
}

type MessageReadEventHandler func(string)

func (ev MessageReadEventHandler) Type() string {
	return "string"
}

func (ev MessageReadEventHandler) Handle(i interface{}) {
	if t, ok := i.(string); ok {
		ev(t)
	}
}
