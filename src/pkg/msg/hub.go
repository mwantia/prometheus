package msg

type MessageHub interface {
	Setup() error

	WriteMessage(string) error

	Cleanup() error
}
