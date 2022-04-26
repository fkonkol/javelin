package messaging

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type MessagingHandler struct{}

func NewHandler() *MessagingHandler {
	return &MessagingHandler{}
}

func (m *MessagingHandler) ConnectionHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrade error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := ws.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
