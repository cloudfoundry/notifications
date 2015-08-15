package info

import "github.com/ryanmoran/stack"

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestLogging stack.Middleware
	RequestCounter stack.Middleware
}

func (r Routes) Register(m muxer) {
	m.Handle("GET", "/info", NewGetHandler(), r.RequestLogging, r.RequestCounter)
}
