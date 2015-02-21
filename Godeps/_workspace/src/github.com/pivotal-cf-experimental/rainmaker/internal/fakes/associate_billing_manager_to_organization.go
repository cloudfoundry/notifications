package fakes

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func (fake *CloudController) AssociateBillingManagerToOrganization(w http.ResponseWriter, req *http.Request) {
	r := regexp.MustCompile(`^/v2/organizations/(.*)/billing_managers/(.*)$`)
	matches := r.FindStringSubmatch(req.URL.Path)

	org, ok := fake.Organizations.Get(matches[1])
	if !ok {
		fake.NotFound(w)
		return
	}

	billingManager, ok := fake.Users.Get(matches[2])
	if !ok {
		fake.NotFound(w)
		return
	}

	org.BillingManagers.Add(billingManager)

	response, err := json.Marshal(org)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
