package info

import (
	"encoding/json"
	"net/http"

	"github.com/ryanmoran/stack"
)

type GetHandler struct {
}

func NewGetHandler() GetHandler {
	return GetHandler{}
}

func (handler GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	w.WriteHeader(http.StatusOK)
	output, err := json.Marshal(map[string]interface{}{
		"version": 2,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(output)
}
