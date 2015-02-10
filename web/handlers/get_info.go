package handlers

import (
	"net/http"

	"github.com/ryanmoran/stack"
)

type GetInfo struct{}

func NewGetInfo() GetInfo {
	return GetInfo{}
}

func (handler GetInfo) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
