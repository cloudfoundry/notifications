package handlers

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, object interface{}) {
	output, err := json.Marshal(object)
	if err != nil {
		panic(err) // No JSON we write into a response should ever panic
	}

	w.WriteHeader(status)
	w.Write(output)
}
