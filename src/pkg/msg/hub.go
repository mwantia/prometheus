package msg

type MessageHub interface {
	CreateProducer(string) (MessageHubProducer, error)

	CreateConsumer(string) (MessageHubConsumer, error)
}
