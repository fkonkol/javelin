package messaging

type Message struct {
	SenderID int    `json:"senderID"`
	Username string `json:"username"`
	Action   string `json:"action"`
	Body     string `json:"body"`
}

var messages []*Message

func (m *Message) Broadcast() {
	for _, client := range clients {
		m.BroadcastTo(client)
	}
	messages = append(messages, m)
}

func (m *Message) BroadcastTo(client *Client) {
	client.Socket.WriteJSON(m)
}
