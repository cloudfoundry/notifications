package fakes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pivotal-cf-experimental/rainmaker/internal/documents"
)

func (fake *CloudController) CreateSpace(w http.ResponseWriter, req *http.Request) {
	var document documents.CreateSpaceRequest
	err := json.NewDecoder(req.Body).Decode(&document)
	if err != nil {
		panic(err)
	}
	now := time.Now().UTC()

	space := NewSpace(NewGUID("space"))
	space.Name = document.Name
	space.OrganizationGUID = document.OrganizationGUID
	space.CreatedAt = now
	space.UpdatedAt = now

	fake.Spaces.Add(space)

	response, err := json.Marshal(space)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
