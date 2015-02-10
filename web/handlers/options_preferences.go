package handlers

import (
	"net/http"

	"github.com/ryanmoran/stack"
)

type OptionsPreferences struct{}

func NewOptionsPreferences() OptionsPreferences {
	return OptionsPreferences{}
}

func (handler OptionsPreferences) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	w.WriteHeader(http.StatusNoContent)
}
