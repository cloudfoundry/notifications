package handlers

import (
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/ryanmoran/stack"
)

type GetInfo struct{}

func NewGetInfo() GetInfo {
    return GetInfo{}
}

func (handler GetInfo) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("{}"))

    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.info",
    }).Log()
}
