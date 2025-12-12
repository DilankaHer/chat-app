package db

import (
	"context"
	"database/sql"
	migration "duhchat/internal/db/migrations"
	"encoding/json"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string `json:"DATABASE_HOST"`
	Port     int    `json:"DATABASE_PORT"`
	User     string `json:"DATABASE_USER"`
	Password string `json:"DATABASE_PASSWORD"`
	DBName   string `json:"DATABASE_NAME"`
}

func ConnectDB() (*sql.DB, error) {
	env := os.Getenv("APP_ENV")
	fileName := ""

	if env == "local" || env == "" {
		fileName = "config/.env.local.json"
	}

	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	config := &DBConfig{}
	if err := json.Unmarshal(file, config); err != nil {
		return nil, err
	}

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
