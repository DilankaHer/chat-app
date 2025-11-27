package handler

import (
	"duhchat/internal/repo"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type LoginHandler struct {
	userRepo *repo.UserRepo
}	

func NewLoginHandler(userRepo *repo.UserRepo) *LoginHandler {
	return &LoginHandler{userRepo: userRepo}
}

func (uh *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Login Failed at Read All", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &uh.userRepo.User)
	if err != nil {
		http.Error(w, "Login Failed at Unmarshal", http.StatusInternalServerError)
		return
	}

	err = uh.userRepo.AddUser()
	if err != nil {
		http.Error(w, "Login Failed at Add User", http.StatusInternalServerError)
		return
	}

	fmt.Println(uh.userRepo.User)

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": uh.userRepo.User.Username,
		"email": uh.userRepo.User.Email,
		"id": uh.userRepo.User.Id,
	}).SignedString([]byte("k8f9+2aV3b7XcQpL6eR1yT0uN4wZ5vQ2"))
	if err != nil {
		http.Error(w, "Login Failed At JWT New", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,    // JavaScript cannot read it
		Secure:   false,   // set to true in production (HTTPS only)
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login Successful"))
}