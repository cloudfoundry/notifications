package info

import (
	"encoding/json"
	"net/http"

	"github.com/ryanmoran/stack"
)

type GetHandler struct {
	version int
}

func NewGetHandler(version int) GetHandler {
	return GetHandler{
		version: version,
	}
}

func (handler GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	w.WriteHeader(http.StatusOK)
	output, err := json.Marshal(map[string]interface{}{
		"version": handler.version,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(output)
}
