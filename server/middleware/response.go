package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ServerResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ResponseCapture struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (rc *ResponseCapture) WriteHeader(statusCode int) {
	rc.status = statusCode
}

func (rc *ResponseCapture) Write(b []byte) (int, error) {
	return rc.body.Write(b)
}

func StandardResponse(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := &ResponseCapture{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(rc, r)
		resp := ServerResponse{
			Status:  rc.status,
			Message: http.StatusText(rc.status),
			Data:    json.RawMessage(rc.body.Bytes()),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(rc.status)
		json.NewEncoder(w).Encode(resp)
	})
}
