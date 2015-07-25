package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type Muxer struct {
	*mux.Router
}

func NewMuxer() Muxer {
	return Muxer{mux.NewRouter()}
}

func (m Muxer) Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware) {
	s := stack.NewStack(handler).Use(middleware...)
	m.Router.Handle(path, s).Methods(method).Name(fmt.Sprintf("%s %s", method, path))
}

func (m Muxer) Match(request *http.Request) http.Handler {
	match := &mux.RouteMatch{}
	ok := m.Router.Match(request, match)
	if !ok {
		return http.HandlerFunc(http.NotFound)
	}

	return match.Handler
}

func (m Muxer) GetRouter() *mux.Router {
	return m.Router
}
