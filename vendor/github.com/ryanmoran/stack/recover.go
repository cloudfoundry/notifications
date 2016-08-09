package stack

import (
    "log"
    "net/http"
    "runtime/debug"
)

type RecoverCallback func(http.ResponseWriter, *http.Request, Context, interface{})

var defaultRecoverCallback = RecoverCallback(func(w http.ResponseWriter, req *http.Request, context Context, err interface{}) {
    log.Println("[Recover]", err)
    log.Println("[Recover]", string(debug.Stack()))

    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("Internal Server Error"))
})

func Recover(w http.ResponseWriter, req *http.Request, callback *RecoverCallback, context Context) {
    err := recover()
    if err != nil {
        if callback != nil {
            (*callback)(w, req, context, err)
        } else {
            defaultRecoverCallback(w, req, context, err)
        }
    }
}
