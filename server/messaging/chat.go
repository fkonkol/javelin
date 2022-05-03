package messaging

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleNewConnection(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	newClient := Client{
		ID:     getNextID(),
		Socket: conn,
	}

	clients = append(clients, &newClient)

	newClient.Listen()
}

var currentID = 0

func getNextID() int {
	currentID += 1
	return currentID
}
