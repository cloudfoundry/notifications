package root

import "github.com/ryanmoran/stack"

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestLogging stack.Middleware
}

func (r Routes) Register(m muxer) {
	m.Handle("GET", "/", NewGetHandler(), r.RequestLogging)
}
