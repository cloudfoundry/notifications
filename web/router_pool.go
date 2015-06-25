package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouterPool() *RouterPool {
	return &RouterPool{}
}

type RouterPool struct {
	routers []MatchableRouter
}

func (rp *RouterPool) Add(router MatchableRouter) {
	rp.routers = append(rp.routers, router)
}

func (rp *RouterPool) AddMux(router *mux.Router) {
	rp.Add(MuxToMatchableRouter(router))
}

func (rp *RouterPool) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, router := range rp.routers {
		if match := router.Match(req); match != nil {
			match.ServeHTTP(w, req)
			return
		}
	}

	http.NotFound(w, req)
}

type MatchableRouter interface {
	Match(request *http.Request) http.Handler
}

func MuxToMatchableRouter(router *mux.Router) MatchableRouter {
	return &matchableMuxRouter{router}
}

type matchableMuxRouter struct {
	router *mux.Router
}

func (m *matchableMuxRouter) Match(req *http.Request) http.Handler {
	routeMatch := &mux.RouteMatch{}
	if m.router.Match(req, routeMatch) {
		return routeMatch.Handler
	}

	return nil
}
