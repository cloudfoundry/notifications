package info

import "github.com/ryanmoran/stack"

type muxer interface {
	Handle(method, path string, handler stack.Handler, middleware ...stack.Middleware)
}

type Routes struct {
	RequestCounter stack.Middleware
	RequestLogging stack.Middleware
}

func (r Routes) Register(m muxer) {
	m.Handle("GET", "/info", NewGetHandler(), r.RequestLogging, r.RequestCounter)
}
