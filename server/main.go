package main

import (
	"duhchat/internal/db"
	migration "duhchat/internal/db/migrations"
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

	// app.Logger.Println("app has started")

	r := routes.SetupRoutes()
	// r.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
    // r.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	db, err := db.ConnectDB()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = migration.RunMigrations(db)
	if err != nil {
		fmt.Println(err.Error())
	}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Fatal Error")
	}
}
