package msg

import "context"

type MessageHubConsumer interface {
	GetName() string

	Read(context.Context, interface{}) error

	Cleanup() error
}

type EventHandler interface {
	Type() string

	Handle(interface{})
}

type ContentEventHandler func(string)

func (ev ContentEventHandler) Type() string {
	return "content"
}

func (ev ContentEventHandler) Handle(i interface{}) {
	if t, ok := i.(string); ok {
		ev(t)
	}
}

type MessageEventHandler func(Message)

func (ev MessageEventHandler) Type() string {
	return "message"
}

func (ev MessageEventHandler) Handle(i interface{}) {
	if t, ok := i.(Message); ok {
		ev(t)
	}
}
