package messaging

type Message struct {
	ThreadID string `json:"threadID"`
	Action   string `json:"action"`
	Body     string `json:"body"`
}

var messages []*Message

func (m *Message) Broadcast() {}
