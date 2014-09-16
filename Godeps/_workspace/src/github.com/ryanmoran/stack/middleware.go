package stack

import "net/http"

type Middleware interface {
    ServeHTTP(http.ResponseWriter, *http.Request, Context) bool
}
