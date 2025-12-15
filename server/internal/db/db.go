package db

import (
	"context"
	"database/sql"
	migration "duhchat/internal/db/migrations"
	"duhchat/util"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func ConnectDB(config *util.AppConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

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
