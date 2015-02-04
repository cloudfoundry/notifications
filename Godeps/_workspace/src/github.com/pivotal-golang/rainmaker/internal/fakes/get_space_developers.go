package fakes

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func (fake *CloudController) GetSpaceDevelopers(w http.ResponseWriter, req *http.Request) {
	r := regexp.MustCompile(`^/v2/spaces/(.*)/developers$`)
	matches := r.FindStringSubmatch(req.URL.Path)

	query := req.URL.Query()
	pageNum := parseInt(query.Get("page"), 1)
	perPage := parseInt(query.Get("results-per-page"), 10)

	space, ok := fake.Spaces.Get(matches[1])
	if !ok {
		fake.NotFound(w)
		return
	}

	page := NewPage(space.Developers, req.URL.Path, pageNum, perPage)
	response, err := json.Marshal(page)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
