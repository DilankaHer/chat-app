package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func ReadMessage(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // allow all origins
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Exit upgrader", err)
		return
	}

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Exit for", err)
			return
		}
		dataByte := []byte("Ypoyopypypyp")
		if err := conn.WriteMessage(messageType, dataByte); err != nil {
			log.Println("Exit for",err)
			return
		}
	}

}