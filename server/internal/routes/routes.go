package routes

import (
	auth "duhchat/internal/api/handler/auth"
	"duhchat/internal/api/handler/message"
	room "duhchat/internal/api/handler/room"
	"duhchat/middleware"
	"duhchat/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRoutes(loginHandler *auth.AuthHandler, joinRoomHandler *room.RoomHandler, messageHandler *message.MessageHandler, config *util.AppConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// public routes (NO middleware)
	r.Get("/rooms", middleware.StandardResponse(joinRoomHandler.GetRooms))
	r.Get("/messages", middleware.StandardResponse(messageHandler.GetMessagesByRoomId))
	r.Post("/signup", middleware.StandardResponse(loginHandler.Signup))
	r.Post("/login", middleware.StandardResponse(loginHandler.Login))

	// private routes (WITH JWT middleware)
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.JWTAuth(config.JWTSecret))
		pr.Get("/joinRoom", joinRoomHandler.JoinRoom)
		pr.Get("/me", middleware.StandardResponse(loginHandler.GetMe))
		pr.Post("/createRoom", middleware.StandardResponse(joinRoomHandler.CreateRoom))
		pr.Delete("/deleteRoom", middleware.StandardResponse(joinRoomHandler.DeleteRoom))
		pr.Post("/logout", middleware.StandardResponse(loginHandler.Logout))
		// add more protected routes here
	})

	// r.Handle("/assets/*", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	// Catch-all to serve index.html for React/Vite SPA
	// r.NotFound(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./static/index.html")
	// })

	return r
}
