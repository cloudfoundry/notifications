package handlers

import "net/http"

type OptionsPreferences struct{}

func NewOptionsPreferences() OptionsPreferences {
    return OptionsPreferences{}
}

func (handler OptionsPreferences) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}
