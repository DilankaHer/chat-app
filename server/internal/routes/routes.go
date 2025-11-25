package routes

import (
	"duhchat/internal/ws"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/test-message", ws.ReadMessage)
	// r.Get("/projects", app.ApiHandler.ProjectsPage)
	// r.Get("/experience", app.ApiHandler.ExperiencePage)
	// r.Get("/download", app.ApiHandler.MakeDownloadLink)
	return r
}
