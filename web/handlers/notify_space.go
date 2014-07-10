package handlers

import "net/http"

type NotifySpace struct{}

func NewNotifySpace() NotifySpace {
    return NotifySpace{}
}

func (handler NotifySpace) ServeHTTP(w http.ResponseWriter, req *http.Request) {}
