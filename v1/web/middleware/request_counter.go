package middleware

import (
	"net/http"
	"strings"

	"fmt"

	"github.com/gorilla/mux"
	"github.com/rcrowley/go-metrics"
	"github.com/ryanmoran/stack"
)

type routeMatcher interface {
	Match(*http.Request, *mux.RouteMatch) bool
}

type RequestCounter struct {
	matcher routeMatcher
}

func NewRequestCounter(matcher routeMatcher) RequestCounter {
	return RequestCounter{
		matcher: matcher,
	}
}

func (ware RequestCounter) ServeHTTP(w http.ResponseWriter, req *http.Request, context stack.Context) bool {
	path := "UNKNOWN"
	var match mux.RouteMatch
	if ok := ware.matcher.Match(req, &match); ok {
		name := match.Route.GetName()
		path = convertNameToMetricPath(name)
	}

	mn := fmt.Sprintf("notifications.web.%s.%s", req.Method, path)
	metrics.GetOrRegisterCounter(mn, nil).Inc(1)

	return true
}

func convertNameToMetricPath(name string) string {
	parts := strings.SplitN(name, " ", -1)
	path := parts[1]

	path = strings.Replace(path, "{", ":", -1)
	path = strings.Replace(path, "}", "", -1)

	return path
}
