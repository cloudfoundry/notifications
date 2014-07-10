package stack

import (
    "log"
    "net/http"
    "runtime/debug"
)

type RecoverCallback func(http.ResponseWriter, *http.Request, interface{})

var defaultRecoverCallback = RecoverCallback(func(w http.ResponseWriter, req *http.Request, err interface{}) {
    log.Println("[Recover]", err)
    log.Println("[Recover]", string(debug.Stack()))

    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("Internal Server Error"))
})

func Recover(w http.ResponseWriter, req *http.Request, callback *RecoverCallback) {
    err := recover()
    if err != nil {
        if callback != nil {
            (*callback)(w, req, err)
        } else {
            defaultRecoverCallback(w, req, err)
        }
    }
}
