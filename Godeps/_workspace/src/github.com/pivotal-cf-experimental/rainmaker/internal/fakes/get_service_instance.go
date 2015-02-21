package fakes

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (fake *CloudController) GetServiceInstance(w http.ResponseWriter, req *http.Request) {
	instanceGUID := strings.TrimPrefix(req.URL.Path, "/v2/service_instances/")

	instance, ok := fake.ServiceInstances.Get(instanceGUID)
	if !ok {
		fake.NotFound(w)
	}

	response, err := json.Marshal(instance)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
