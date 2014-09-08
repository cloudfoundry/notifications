package middleware

import "net/http"

type CORS struct{}

func NewCORS() CORS {
    return CORS{}
}

func (ware CORS) ServeHTTP(w http.ResponseWriter, req *http.Request) bool {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET")
    w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization")

    return true
}
