package util

import (
	"encoding/json"
	"net/http"
)

func JSONMarshaller(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	resultBytes := append([]byte(message+"\n"), bytes...)
	w.WriteHeader(statusCode)
	w.Write(resultBytes)
}
