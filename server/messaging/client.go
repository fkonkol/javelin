package messaging

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       int `json:"id"`
	Username string
	Socket   *websocket.Conn
}

var clients []*Client

// Listens for inbound messages.
func (c *Client) Listen() {
	var input map[string]string
	for {
		err := c.Socket.ReadJSON(&input)
		if err != nil {
			log.Printf("Releasing websocket connection: %d\n", c.ID)
			c.Socket.Close()
		}

		HandleMessage(c, input)
	}
}

func HandleMessage(c *Client, input map[string]string) {
	switch input["action"] {
	case "INITIAL_CONNECTION":
		c.Username = input["username"]

		newConnMsg := Message{
			SenderID: 0,
			Username: "System",
			Body:     fmt.Sprintf("%s has joined the chat", c.Username),
		}
		newConnMsg.Broadcast()

	case "WRITE_MESSAGE":
		message := Message{
			SenderID: c.ID,
			Username: c.Username,
			Body:     input["body"],
		}

		message.Broadcast()
	}
}
