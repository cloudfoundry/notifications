package middleware

import (
	"net/http"

	"github.com/ryanmoran/stack"
)

type CORS struct {
	origin string
}

func NewCORS(origin string) CORS {
	return CORS{
		origin: origin,
	}
}

func (ware CORS) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) bool {
	w.Header().Set("Access-Control-Allow-Origin", ware.origin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")

	return true
}
