package messaging

import "github.com/gorilla/websocket"

type Client struct {
	ID     int `json:"id"`
	Socket *websocket.Conn
}

var clients []*Client

// Listens for inbound messages.
func (c *Client) Listen() {
	var message *Message
	for {
		err := c.Socket.ReadJSON(&message)
		if err != nil {
			// Release connection
		}

		HandleMessage(c, message)
	}
}

func HandleMessage(c *Client, m *Message) {
	switch m.Action {
	case "INITIAL_CONNECTION":
		// Client joined a thread specified in message struct

	case "WRITE_MESSAGE":
	}
}
