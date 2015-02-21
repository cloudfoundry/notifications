package fakes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pivotal-cf-experimental/rainmaker/internal/documents"
)

func (fake *CloudController) CreateServiceInstance(w http.ResponseWriter, req *http.Request) {
	var document documents.CreateServiceInstanceRequest
	err := json.NewDecoder(req.Body).Decode(&document)
	if err != nil {
		panic(err)
	}

	now := time.Now().UTC()
	instance := NewServiceInstance(NewGUID("service-instance"))
	instance.Name = document.Name
	instance.PlanGUID = document.PlanGUID
	instance.SpaceGUID = document.SpaceGUID
	instance.CreatedAt = now
	instance.UpdatedAt = now

	fake.ServiceInstances.Add(instance)

	response, err := json.Marshal(instance)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
