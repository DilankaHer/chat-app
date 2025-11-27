package middleware

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jwtKey = []byte("k8f9+2aV3b7XcQpL6eR1yT0uN4wZ5vQ2")
		fmt.Println(r.Cookies())
		cookie, err := r.Cookie("jwt")
		if err != nil {
			fmt.Println("missing auth token")
			// util.WriteError(w, http.StatusUnauthorized, "missing auth token")
			return
		}

		claims := jwt.MapClaims{}

		value := cookie.Value
		method := jwt.GetSigningMethod(value)
		fmt.Println("Method", method)

		token, err := jwt.ParseWithClaims(value, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})
		if err != nil {
			fmt.Println(err)
			return 
		}

		fmt.Println(token.Claims)
	})
}