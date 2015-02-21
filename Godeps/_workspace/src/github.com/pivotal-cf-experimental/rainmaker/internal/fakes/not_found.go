package fakes

import (
	"encoding/json"
	"net/http"
)

func (fake *CloudController) NotFound(w http.ResponseWriter) {
	errorBody, err := json.Marshal(map[string]interface{}{
		"code":        10000,
		"description": "Unknown request",
		"error_code":  "CF-NotFound",
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write(errorBody)
}
