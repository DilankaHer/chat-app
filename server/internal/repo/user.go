package repo

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *sql.DB
	User *User
}

type User struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password_hash"`
}

func NewUserRepo(db *sql.DB) *UserRepo {	
	return &UserRepo{db: db, User: &User{}}
}

func (ur *UserRepo) AddUser() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(ur.User.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return err
	}
	ur.User.Password = string(hash)
	tx, err := ur.db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2	, $3) RETURNING id`

	err = tx.QueryRow(query, ur.User.Username, ur.User.Email, ur.User.Password).Scan(&ur.User.Id)
	if err != nil {
		fmt.Println("Insert Error", err)
		return err
	}	
	
	tx.Commit()

	return nil
}