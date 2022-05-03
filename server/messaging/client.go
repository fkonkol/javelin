package messaging

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

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
			log.Printf("Releasing websocket connection: %d\n", c.ID)
			c.Socket.Close()
		}

		HandleMessage(c, message)
	}
}

func HandleMessage(c *Client, m *Message) {
	switch m.Action {
	case "INITIAL_CONNECTION":
		newConnMsg := Message{
			SenderID: 0,
			Username: "System",
			Body:     fmt.Sprintf("%d has joined the chat", c.ID),
		}
		newConnMsg.Broadcast()

	case "WRITE_MESSAGE":
	}
}
