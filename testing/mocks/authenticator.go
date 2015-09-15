package mocks

import (
	"net/http"

	"github.com/ryanmoran/stack"
)

type Authenticator struct {
	ServeHTTPCall struct {
		Receives struct {
			Writer  http.ResponseWriter
			Request *http.Request
			Context stack.Context
		}
		Returns struct {
			Continue bool
		}
	}
}

func (a *Authenticator) ServeHTTP(writer http.ResponseWriter, request *http.Request, context stack.Context) bool {
	a.ServeHTTPCall.Receives.Writer = writer
	a.ServeHTTPCall.Receives.Request = request
	a.ServeHTTPCall.Receives.Context = context
	return a.ServeHTTPCall.Returns.Continue
}
