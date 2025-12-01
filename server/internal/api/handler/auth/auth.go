package handler

import (
	"duhchat/internal/repo"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type Response struct {
	Message string      `json:"message"`
	Rooms   []repo.Room `json:"rooms"`
	UserId  string      `json:"userId"`
}

type Login struct {
	EmailUsername string `json:"emailUsername" validate:"required"`
	Password      string `json:"password" validate:"required"`
}

type Signup struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginHandler struct {
	userRepo repo.UserRepository
}

func NewLoginHandler(userRepo repo.UserRepository) *LoginHandler {
	return &LoginHandler{userRepo: userRepo}
}

func (uh *LoginHandler) Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signup")
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Signup Failed at Read All", http.StatusInternalServerError)
		return
	}

	signup := &Signup{}
	err = json.Unmarshal(body, signup)
	if err != nil {
		http.Error(w, "SignUp Failed at Unmarshal: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = validator.New().Struct(signup)
	if err != nil {
		http.Error(w, "SignUp Failed at Struct Level: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uh.userRepo.AddUser(signup.Email, signup.Username, signup.Password)
	if err != nil {
		http.Error(w, "SignUp Failed at Add User", http.StatusInternalServerError)
		return
	}

	err = SetJWTCookie(&w, user)
	if err != nil {
		http.Error(w, "SignUp Failed at Set JWT Cookie", http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: "SignUp Successful",
		UserId:  user.Id,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (uh *LoginHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(jwt.MapClaims)["userId"].(string)
	user, err := uh.userRepo.GetUser(userId)
	if err != nil {
		http.Error(w, "Failed to Get User", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uh *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Login Failed at Read All", http.StatusInternalServerError)
		return
	}

	login := &Login{}
	err = json.Unmarshal(body, login)
	if err != nil {
		http.Error(w, "Login Failed at Unmarshal: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = validator.New().Struct(login)
	if err != nil {
		http.Error(w, "Login Failed at Struct Level: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := uh.userRepo.Login(login.EmailUsername, login.Password)
	if err != nil {
		http.Error(w, "Credentials are not correct", http.StatusUnauthorized)
		return
	}

	err = SetJWTCookie(&w, user)
	if err != nil {
		http.Error(w, "Login Failed at Set JWT Cookie", http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: "Login Successful",
		UserId:  user.Id,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func SetJWTCookie(w *http.ResponseWriter, user *repo.User) error {
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"userId":   user.Id,
	}).SignedString([]byte("k8f9+2aV3b7XcQpL6eR1yT0uN4wZ5vQ2"))
	if err != nil {
		return err
	}

	http.SetCookie(*w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,  // JavaScript cannot read it
		Secure:   false, // set to true in production (HTTPS only)
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}
