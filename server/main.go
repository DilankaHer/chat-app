package main

import (
	"duhchat/internal/api/handler"
	"duhchat/internal/db"
	"duhchat/internal/repo"
	"duhchat/internal/routes"
	"fmt"
	"net/http"
	"time"
)

type Count struct {
	Count int
}

func main() {
	// app, err := app.NewApplication()
	// if err != nil {
	// 	panic(err)
	// }

	db, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	userRepo := repo.NewUserRepo(db)	
	loginHandler := handler.NewLoginHandler(userRepo)	

	r := routes.SetupRoutes(loginHandler)
	// r.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	// r.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Fatal Error")
		panic(err)
	}
}
