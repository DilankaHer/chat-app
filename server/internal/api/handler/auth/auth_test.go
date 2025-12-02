package handler

import (
	"bytes"
	"duhchat/internal/repo"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserRepo struct {
}

func (m *mockUserRepo) AddUser(user *repo.User) error {
	user.UserId = "283120389013489210"
	user.Email = "john@gmail.com"
	user.Username = "john"
	user.EmailUsername = "john"
	return nil
}

func (m *mockUserRepo) GetUser(user *repo.User) error {
	user.UserId = "283120389013489210"
	user.Email = "john@gmail.com"
	user.Username = "john"
	user.EmailUsername = "john"
	return nil
}

func (m *mockUserRepo) Login(user *repo.User) error {
	user.UserId = "283120389013489210"
	user.Email = "john@gmail.com"
	user.Username = "john"
	user.EmailUsername = "john"
	return nil
}
func TestLogin(t *testing.T) {
	mockRepo := &mockUserRepo{}

	ah := NewAuthHandler(mockRepo)

	jsonBody := `{
		"emailUsername": "john",
		"password": "password"
	}`

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(jsonBody))
	w := httptest.NewRecorder()

	ah.Login(w, req)

	resp := w.Result()

	fmt.Println(resp)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// Check response body
	body := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	if err != nil {
		t.Fatalf("can't unmarshal response body")
	}
	if body.Message != "Login Successful" {
		t.Fatalf("expected Login Successful, got %s", body)
	}
	if body.UserId == "" {
		t.Fatalf("expected UserId to not be empty, got %s", body)
	}

	// Check if cookie is set
	cookie := resp.Cookies()
	if len(cookie) == 0 {
		t.Fatal("expected JWT cookie to be set")
	}

	if cookie[0].Name != "jwt" {
		t.Fatalf("expected cookie name 'jwt', got %s", cookie[0].Name)
	}

	if cookie[0].Value != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG5AZ21haWwuY29tIiwidXNlcklkIjoiMjgzMTIwMzg5MDEzNDg5MjEwIiwidXNlcm5hbWUiOiJqb2huIn0.yVmb4FenrNcpNMpOdRgLrUtSNx8o2Ajv5vFMOO6EjPQ" {
		t.Fatalf("expected cookie to match, got %s", cookie[0].Value)
	}
}
