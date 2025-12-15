package util

import (
	"encoding/json"
	"errors"
	"os"
)

type AppConfig struct {
	Host      string `json:"DATABASE_HOST"`
	Port      int    `json:"DATABASE_PORT"`
	User      string `json:"DATABASE_USER"`
	Password  string `json:"DATABASE_PASSWORD"`
	DBName    string `json:"DATABASE_NAME"`
	JWTSecret string `json:"JWT_SECRET"`
}

func GetAppConfig() (*AppConfig, error) {
	env := os.Getenv("APP_ENV")
	fileName := ""

	switch env {
	case "local", "":
		fileName = "config/.env.local.json"
	case "development":
		fileName = "config/.env.dev.json"
	case "production":
		fileName = "config/.env.prod.json"
	default:
		return nil, errors.New("invalid app environment")
	}

	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	config := &AppConfig{}
	if err := json.Unmarshal(file, config); err != nil {
		return nil, err
	}

	return config, nil
}
