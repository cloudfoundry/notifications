package fakes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pivotal-cf-experimental/rainmaker/internal/documents"
)

func (fake *CloudController) CreateOrganization(w http.ResponseWriter, req *http.Request) {
	var document documents.CreateOrganizationRequest
	now := time.Now().UTC()
	err := json.NewDecoder(req.Body).Decode(&document)
	if err != nil {
		panic(err)
	}

	organization := NewOrganization(NewGUID("org"))
	organization.Name = document.Name
	organization.CreatedAt = now
	organization.UpdatedAt = now

	fake.Organizations.Add(organization)

	response, err := json.Marshal(organization)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
