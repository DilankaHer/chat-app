package repo

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	AddUser(email string, username string, password string) (*User, error)
	GetUser(userId string) (*User, error)
	Login(emailUsername string, password string) (*User, error)
}

type UserRepo struct {
	db *sql.DB
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &UserRepo{db: db}
}

func (ur *UserRepo) AddUser(email string, username string, password string) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	tx, err := ur.db.Begin()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`

	user := &User{}
	err = tx.QueryRow(query, username, email, passwordHash).Scan(&user.Id)
	if err != nil {
		fmt.Println("Insert Error", err)
		return nil, err
	}

	tx.Commit()

	return user, nil
}

func (ur *UserRepo) GetUser(userId string) (*User, error) {
	query := `SELECT id, username, email FROM users WHERE id = $1`
	row := ur.db.QueryRow(query, userId)
	user := &User{}
	err := row.Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepo) Login(emailUsername string, password string) (*User, error) {
	query := `SELECT id, username, email, password_hash FROM users WHERE email = $1 OR username = $1`
	row := ur.db.QueryRow(query, emailUsername)
	user := &User{}
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}
