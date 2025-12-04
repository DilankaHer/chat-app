package model

type User struct {
	UserId        string `json:"userId"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	EmailUsername string `json:"emailUsername"`
}
