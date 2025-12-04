package repo

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Login(user *User) error
	AddUser(user *User) error
	GetUser(user *User) error
}
type UserRepo struct {
	db *sql.DB
}

type User struct {
	UserId        string `json:"userId"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	EmailUsername string `json:"-"`
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &UserRepo{db: db}
}

func (ur *UserRepo) Login(user *User) error {
	fmt.Println("Login", user)
	query := `SELECT id, username, email, password_hash FROM users WHERE email = $1 OR username = $1`

	var passwordHash []byte
	err := ur.db.QueryRow(query, user.EmailUsername).Scan(&user.UserId, &user.Username, &user.Email, &passwordHash)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(user.Password))
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepo) AddUser(user *User) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return err
	}
	tx, err := ur.db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`

	err = tx.QueryRow(query, user.Username, user.Email, passwordHash).Scan(&user.UserId)
	if err != nil {
		fmt.Println("Insert Error", err)
		return err
	}

	tx.Commit()

	return nil
}

func (ur *UserRepo) GetUser(user *User) error {
	query := `SELECT username, email FROM users WHERE id = $1`
	err := ur.db.QueryRow(query, user.UserId).Scan(&user.Username, &user.Email)
	if err != nil {
		return err
	}
	return nil
}
