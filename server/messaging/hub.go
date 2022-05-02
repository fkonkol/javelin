package messaging

// Primary message handler.
type Hub struct {
	// Register requests from clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// Handle inbound messages.
	ReadMessage chan *Message

	// Handle outbound messages.
	WriteMessage chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		ReadMessage:  make(chan *Message),
		WriteMessage: make(chan *Message),
	}
}
