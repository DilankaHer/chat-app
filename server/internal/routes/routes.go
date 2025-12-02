package routes

import (
	auth "duhchat/internal/api/handler/auth"
	"duhchat/internal/api/handler/message"
	room "duhchat/internal/api/handler/room"
	"duhchat/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRoutes(loginHandler *auth.AuthHandler, joinRoomHandler *room.RoomHandler, messageHandler *message.MessageHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://localhost:8080", "http://127.0.0.1:8080", "https://yappr.chat", "http://yappr.chat"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// public routes (NO middleware)
	r.Get("/rooms", joinRoomHandler.GetRooms)
	r.Get("/messages", messageHandler.GetMessagesByRoomId)
	r.Post("/signup", loginHandler.Signup)
	r.Post("/login", loginHandler.Login)

	// private routes (WITH JWT middleware)
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.JWTAuth)
		pr.Get("/joinRoom", joinRoomHandler.JoinRoom)
		pr.Get("/me", loginHandler.GetMe)
		pr.Post("/createRoom", joinRoomHandler.CreateRoom)
		// add more protected routes here
	})

	return r
}
