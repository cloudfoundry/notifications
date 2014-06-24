package middleware

import "net/http"

type Middleware interface {
    ServeHTTP(http.ResponseWriter, *http.Request) bool
}
