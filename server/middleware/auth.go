package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jwtKey = []byte("k8f9+2aV3b7XcQpL6eR1yT0uN4wZ5vQ2")
		cookie, err := r.Cookie("jwt")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ServerResponse{
				Status:  http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized),
				Error:   "missing auth token",
			})
			return
		}

		claims := jwt.MapClaims{}

		value := cookie.Value

		token, err := jwt.ParseWithClaims(value, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ServerResponse{
				Status:  http.StatusUnauthorized,
				Message: http.StatusText(http.StatusUnauthorized),
				Error:   "invalid token",
			})
			return // STOP the chain
		}

		// Optional: attach user info to context
		ctx := context.WithValue(r.Context(), "user", claims)
		r = r.WithContext(ctx)

		// Token OK → call next handler
		h.ServeHTTP(w, r)
	})
}
