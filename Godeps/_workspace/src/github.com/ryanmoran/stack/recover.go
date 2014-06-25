package stack

import (
    "log"
    "net/http"
    "runtime/debug"
)

func Recover(w http.ResponseWriter, req *http.Request) {
    err := recover()
    if err != nil {
        log.Println("[Recover]", err)
        log.Println("[Recover]", string(debug.Stack()))
    }
}
