package msg

import "time"

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Sequence  int       `json:"sequence,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}
