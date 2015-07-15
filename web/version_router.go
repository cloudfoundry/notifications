package web

import (
	"net/http"
	"strconv"
)

type VersionRouter map[int]http.Handler

func (vr VersionRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	versionHeader := req.Header.Get("X-NOTIFICATIONS-VERSION")
	if versionHeader == "" {
		versionHeader = "1"
	}

	version, err := strconv.ParseInt(versionHeader, 10, 64)
	if err != nil {
		http.NotFoundHandler().ServeHTTP(w, req)
		return
	}

	if handler, ok := vr[int(version)]; ok {
		handler.ServeHTTP(w, req)
		return
	}

	http.NotFoundHandler().ServeHTTP(w, req)
}
