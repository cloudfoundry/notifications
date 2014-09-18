package handlers

import (
    "net/http"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/ryanmoran/stack"
)

type OptionsPreferences struct{}

func NewOptionsPreferences() OptionsPreferences {
    return OptionsPreferences{}
}

func (handler OptionsPreferences) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) {
    metrics.NewMetric("counter", map[string]interface{}{
        "name": "notifications.web.preferences.options",
    }).Log()

    w.WriteHeader(http.StatusNoContent)
}
