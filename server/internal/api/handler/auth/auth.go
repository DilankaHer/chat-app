package handler

import (
	"duhchat/internal/repo"
	"duhchat/util"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type Response struct {
	Message string `json:"message"`
	UserId  string `json:"userId"`
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

type AuthHandler struct {
	userRepository repo.UserRepository
}

func NewAuthHandler(userRepository repo.UserRepository) *AuthHandler {
	return &AuthHandler{userRepository: userRepository}
}

func (ah *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
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

	user := repo.User{Email: signup.Email, Username: signup.Username, Password: signup.Password}
	err = ah.userRepository.AddUser(&user)
	if err != nil {
		http.Error(w, "SignUp Failed at Add User", http.StatusInternalServerError)
		return
	}

	err = SetJWTCookie(&w, &user, false)
	if err != nil {
		http.Error(w, "SignUp Failed at Set JWT Cookie", http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: "SignUp Successful",
		UserId:  user.UserId,
	}

	util.JSONMarshaller(w, http.StatusOK, response, "SignUp Successful")
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		util.JSONMarshaller(w, http.StatusInternalServerError, "Login Failed", http.StatusText(http.StatusInternalServerError))
		return
	}

	login := &Login{}
	err = json.Unmarshal(body, login)
	if err != nil {
		util.JSONMarshaller(w, http.StatusBadRequest, "Login Failed", http.StatusText(http.StatusBadRequest))
		return
	}

	err = validator.New().Struct(login)
	if err != nil {
		util.JSONMarshaller(w, http.StatusBadRequest, "Login failed at struct level: "+err.Error(), http.StatusText(http.StatusBadRequest))
		return
	}

	user := repo.User{EmailUsername: login.EmailUsername, Password: login.Password}
	err = ah.userRepository.Login(&user)
	if err != nil {
		util.JSONMarshaller(w, http.StatusUnauthorized, err.Error(), http.StatusText(http.StatusUnauthorized))
		return
	}

	err = SetJWTCookie(&w, &user, false)
	if err != nil {
		util.JSONMarshaller(w, http.StatusInternalServerError, "Login Failed", http.StatusText(http.StatusInternalServerError))
		return
	}

	response := Response{
		Message: "Login Successful",
		UserId:  user.UserId,
	}

	util.JSONMarshaller(w, http.StatusOK, response, "Login Successful")
}

func (ah *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(jwt.MapClaims)["userId"].(string)
	user := repo.User{UserId: userId}
	err := ah.userRepository.GetUser(&user)
	if err != nil {
		_ = SetJWTCookie(&w, nil, true)
		util.JSONMarshaller(w, http.StatusNotFound, "please login again", http.StatusText(http.StatusNotFound))
		return
	}
	util.JSONMarshaller(w, http.StatusOK, user, http.StatusText(http.StatusOK))
}

func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := SetJWTCookie(&w, nil, true)
	if err != nil {
		util.JSONMarshaller(w, http.StatusInternalServerError, "Logout Failed", http.StatusText(http.StatusInternalServerError))
		return
	}
	util.JSONMarshaller(w, http.StatusOK, nil, "Logout Successful")
}

func SetJWTCookie(w *http.ResponseWriter, user *repo.User, isLogout bool) error {
	var tokenString = ""
	var maxAge = -1
	var err error

	if !isLogout {
		maxAge = 60 * 60 * 24 * 2
		tokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"email":    user.Email,
			"userId":   user.UserId,
		}).SignedString([]byte("k8f9+2aV3b7XcQpL6eR1yT0uN4wZ5vQ2"))
		if err != nil {
			return err
		}
	}

	http.SetCookie(*w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set to true in production (HTTPS only)
		MaxAge:   maxAge,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}
