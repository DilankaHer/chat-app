package api

import (
	"duhchat/internal/ws"
	"fmt"
	"net/http"
)

func Message(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Innnnnnnnnn")
	ws.ReadMessage(w, r)
}