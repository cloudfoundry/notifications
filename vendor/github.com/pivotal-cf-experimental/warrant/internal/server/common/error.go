package common

import (
	"fmt"
	"net/http"
)

func Error(w http.ResponseWriter, status int, message, errorType string) {
	output := fmt.Sprintf(`{"message":"%s","error":"%s"}`, message, errorType)

	w.WriteHeader(status)
	w.Write([]byte(output))
}
