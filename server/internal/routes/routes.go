package routes

import (
	"duhchat/internal/api"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	// r.Get("/test-message", func(w http.ResponseWriter, r *http.Request) {
	// 	ws.ReadMessage(hub, w, r)
	// })
	r.Get("/joinRoom", api.JoinRoom)
	return r
}
