package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "admin", "password", "mydatabase")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
    if pingErr != nil {
        return nil, pingErr
    }
    fmt.Println("Connected!")

	return db, nil
}