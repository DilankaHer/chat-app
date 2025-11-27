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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "admin", "password", "mydatabase")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	fmt.Println("Open")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	chn := make(chan error, 1)
	
	go func(ctx context.Context) {
		select {
		case chn <- db.Ping():
			return 
		case <- ctx.Done():
			chn <- ctx.Err() 
			return 
		}
	}(ctx)

	result := <-chn
	if result != nil {
		return nil, result
	}

	err = migration.RunMigrations(db)
	if err != nil {
		return nil, err
	}
    fmt.Println("Connected!")

	return db, nil
}