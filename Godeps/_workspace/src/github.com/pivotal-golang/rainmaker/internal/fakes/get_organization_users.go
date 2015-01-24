package fakes

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func (fake *CloudController) GetOrganizationUsers(w http.ResponseWriter, req *http.Request) {
	r := regexp.MustCompile(`^/v2/organizations/(.*)/users$`)
	matches := r.FindStringSubmatch(req.URL.Path)

	query := req.URL.Query()
	pageNum := ParseInt(query.Get("page"), 1)
	perPage := ParseInt(query.Get("results-per-page"), 10)

	org, ok := fake.Organizations.Get(matches[1])
	if !ok {
		fake.NotFound(w)
		return
	}

	page := NewPage(org.Users, req.URL.Path, pageNum, perPage)
	response, err := json.Marshal(page)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
