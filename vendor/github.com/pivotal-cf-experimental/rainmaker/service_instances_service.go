package rainmaker

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf-experimental/rainmaker/internal/documents"
)

type ServiceInstancesService struct {
	config Config
}

func NewServiceInstancesService(config Config) *ServiceInstancesService {
	return &ServiceInstancesService{
		config: config,
	}
}

func (service *ServiceInstancesService) Create(name, planGUID, spaceGUID, token string) (ServiceInstance, error) {
	_, body, err := NewClient(service.config).makeRequest(requestArguments{
		Method: "POST",
		Path:   "/v2/service_instances",
		Body: documents.CreateServiceInstanceRequest{
			Name:      name,
			PlanGUID:  planGUID,
			SpaceGUID: spaceGUID,
		},
		Token: token,
		AcceptableStatusCodes: []int{http.StatusCreated},
	})
	if err != nil {
		return ServiceInstance{}, err
	}

	var response documents.ServiceInstanceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	return newServiceInstanceFromResponse(service.config, response), nil
}

func (service *ServiceInstancesService) Get(instanceGUID, token string) (ServiceInstance, error) {
	_, body, err := NewClient(service.config).makeRequest(requestArguments{
		Method: "GET",
		Path:   "/v2/service_instances/" + instanceGUID,
		Token:  token,
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return ServiceInstance{}, err
	}

	var response documents.ServiceInstanceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}
	return newServiceInstanceFromResponse(service.config, response), nil
}
