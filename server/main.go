package main

import (
	auth "duhchat/internal/api/handler/auth"
	"duhchat/internal/api/handler/message"
	room "duhchat/internal/api/handler/room"
	"duhchat/internal/db"
	"duhchat/internal/repo"
	"duhchat/internal/routes"
	"duhchat/internal/ws"
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
	roomRepo := repo.NewRoomRepo(db)
	messageRepo := repo.NewMessageRepo(db)

	hub := ws.NewHub(roomRepo)
	if err := hub.CreateDefaultRooms(); err != nil {
		panic(err)
	}

	loginHandler := auth.NewAuthHandler(userRepo)
	roomUserHandler := room.NewRoomHandler(roomRepo, messageRepo, hub)
	messageHandler := message.NewMessageHandler(messageRepo)

	r := routes.SetupRoutes(loginHandler, roomUserHandler, messageHandler)
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
