package db

import (
	"context"
	"database/sql"
	migration "duhchat/internal/db/migrations"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "postgres", 5432, "admin", "password", "mydatabase")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	fmt.Println("Open")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	if err := migration.RunMigrations(db); err != nil {
		return nil, err
	}
	fmt.Println("Connected!")

	return db, nil
}
