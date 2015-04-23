package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"
)

type routeMatcher interface {
	Match(*http.Request, *mux.RouteMatch) bool
}

type RequestCounter struct {
	matcher routeMatcher
	logger  *log.Logger
}

func NewRequestCounter(matcher routeMatcher, logger *log.Logger) RequestCounter {
	return RequestCounter{
		matcher: matcher,
		logger:  logger,
	}
}

func (ware RequestCounter) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) bool {
	path := "UNKNOWN"
	var match mux.RouteMatch
	if ok := ware.matcher.Match(req, &match); ok {
		name := match.Route.GetName()
		path = convertNameToMetricPath(name)
	}

	m := metrics.NewMetric("counter", map[string]interface{}{
		"name": "notifications.web",
		"tags": map[string]string{
			"endpoint": req.Method + path,
		},
	})
	m.LogWith(ware.logger)

	return true
}

func convertNameToMetricPath(name string) string {
	parts := strings.SplitN(name, " ", -1)
	path := parts[1]

	path = strings.Replace(path, "{", ":", -1)
	path = strings.Replace(path, "}", "", -1)

	return path
}
