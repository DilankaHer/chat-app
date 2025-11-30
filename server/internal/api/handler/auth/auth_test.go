package handler

import (
	"bytes"
	"duhchat/internal/repo"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserRepo struct {
	AddUserFn func(email string, username string, password string) (*repo.User, error)
	GetUserFn func(userId string) (*repo.User, error)
	LoginFn   func(emailUsername string, password string) (*repo.User, error)
}

func (m *mockUserRepo) AddUser(email string, username string, password string) (*repo.User, error) {
	return m.AddUserFn(email, username, password)
}

func (m *mockUserRepo) GetUser(userId string) (*repo.User, error) {
	return m.GetUserFn(userId)
}

func (m *mockUserRepo) Login(emailUsername string, password string) (*repo.User, error) {
	return m.LoginFn(emailUsername, password)
}
func TestLogin(t *testing.T) {
	mockRepo := &mockUserRepo{
		AddUserFn: func(email string, username string, password string) (*repo.User, error) {
			user := &repo.User{
				Id:       "1",
				Username: username,
				Email:    email,
				Password: password,
			}
			return user, nil
		},
		GetUserFn: func(userId string) (*repo.User, error) {
			return &repo.User{
				Id:       "1",
				Username: "john",
				Email:    "john@example.com",
				Password: "",
			}, nil
		},
		LoginFn: func(emailUsername string, password string) (*repo.User, error) {
			return &repo.User{
				Id:       "1",
				Username: "john",
				Email:    "john@example.com",
				Password: "",
			}, nil
		},
	}

	lh := NewLoginHandler(mockRepo)

	jsonBody := `{
		"username": "john",
		"email": "john@example.com",
		"password": "password"
	}`

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(jsonBody))
	w := httptest.NewRecorder()

	lh.Login(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// Check response body
	body := w.Body.String()
	if body != "Login Successful" {
		t.Fatalf("expected Login Successful, got %s", body)
	}

	// Check if cookie is set
	cookie := resp.Cookies()
	if len(cookie) == 0 {
		t.Fatal("expected JWT cookie to be set")
	}

	if cookie[0].Name != "jwt" {
		t.Fatalf("expected cookie name 'jwt', got %s", cookie[0].Name)
	}
}
