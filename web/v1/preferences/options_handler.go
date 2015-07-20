package preferences

import (
	"net/http"

	"github.com/ryanmoran/stack"
)

type OptionsHandler struct{}

func NewOptionsHandler() OptionsHandler {
	return OptionsHandler{}
}

func (h OptionsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
	w.WriteHeader(http.StatusNoContent)
}
