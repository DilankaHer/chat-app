package api

import (
	"duhchat/internal/ws"
	"fmt"
	"net/http"
)

// func Message(w http.ResponseWriter, r *http.Request) {
// 	hub := ws.NewHub()
// 	go hub.Run()
// 	fmt.Println("Received request at /test-message")
// 	fmt.Println(r.Body)
// 	ws.ReadMessage(hub, w, r)
// }

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	hubs := ws.GetHubs()
	fmt.Println("Available Hub IDs:", hubs)

	roomId := r.URL.Query().Get("roomId")
	clientId := r.URL.Query().Get("clientId")
	hub, exists := hubs[roomId]
	if !exists {
		hub = ws.NewHub(roomId)
		hubs[roomId] = hub
		go hub.Run()
	}
	ws.ReadMessage(hub, clientId, w, r)

	fmt.Println("Join room endpoint hit")
}
