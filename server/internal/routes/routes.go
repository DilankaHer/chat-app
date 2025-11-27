package routes

import (
	"duhchat/internal/api"
	"duhchat/internal/api/handler"
	"duhchat/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(loginHandler *handler.LoginHandler) *chi.Mux {
    r := chi.NewRouter()

    // public routes (NO middleware)
    r.Post("/login", loginHandler.Login)

    // private routes (WITH JWT middleware)
    r.Group(func(pr chi.Router) {
        pr.Use(middleware.JWTAuth)
        pr.Get("/joinRoom", api.JoinRoom)
        // add more protected routes here
    })

    return r	
}
