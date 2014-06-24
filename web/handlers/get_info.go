package handlers

import "net/http"

type GetInfo struct{}

func NewGetInfo() GetInfo {
    return GetInfo{}
}

func (handler GetInfo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{}`))
}
