package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ryanmoran/stack"
)

type GetInfo struct {
	version int
}

func NewGetInfo(version int) GetInfo {
	return GetInfo{
		version: version,
	}
}

func (handler GetInfo) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
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
