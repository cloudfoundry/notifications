package rainmaker

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/pivotal-golang/rainmaker/internal/documents"
)

type Space struct {
	config                   Config
	GUID                     string
	URL                      string
	CreatedAt                time.Time
	UpdatedAt                time.Time
	Name                     string
	OrganizationGUID         string
	SpaceQuotaDefinitionGUID string
	OrganizationURL          string
	DevelopersURL            string
	ManagersURL              string
	AuditorsURL              string
	AppsURL                  string
	RoutesURL                string
	DomainsURL               string
	ServiceInstancesURL      string
	AppEventsURL             string
	EventsURL                string
	SecurityGroupsURL        string
	Developers               UsersList
}

func NewSpace(config Config, guid string) Space {
	return Space{
		config:     config,
		GUID:       guid,
		Developers: NewUsersList(config, NewRequestPlan("/v2/spaces/"+guid+"/developers", url.Values{})),
	}
}

func NewSpaceFromResponse(config Config, response documents.SpaceResponse) Space {
	space := NewSpace(config, response.Metadata.GUID)
	if response.Metadata.CreatedAt == nil {
		response.Metadata.CreatedAt = &time.Time{}
	}

	if response.Metadata.UpdatedAt == nil {
		response.Metadata.UpdatedAt = &time.Time{}
	}

	space.URL = response.Metadata.URL
	space.CreatedAt = *response.Metadata.CreatedAt
	space.UpdatedAt = *response.Metadata.UpdatedAt
	space.Name = response.Entity.Name
	space.OrganizationGUID = response.Entity.OrganizationGUID
	space.SpaceQuotaDefinitionGUID = response.Entity.SpaceQuotaDefinitionGUID
	space.OrganizationURL = response.Entity.OrganizationURL
	space.DevelopersURL = response.Entity.DevelopersURL
	space.ManagersURL = response.Entity.ManagersURL
	space.AuditorsURL = response.Entity.AuditorsURL
	space.AppsURL = response.Entity.AppsURL
	space.RoutesURL = response.Entity.RoutesURL
	space.DomainsURL = response.Entity.DomainsURL
	space.ServiceInstancesURL = response.Entity.ServiceInstancesURL
	space.AppEventsURL = response.Entity.AppEventsURL
	space.EventsURL = response.Entity.EventsURL
	space.SecurityGroupsURL = response.Entity.SecurityGroupsURL

	return space
}

func FetchSpace(config Config, path, token string) (Space, error) {
	_, body, err := NewClient(config).makeRequest(requestArguments{
		Method: "GET",
		Path:   path,
		Token:  token,
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return Space{}, err
	}

	var response documents.SpaceResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	return NewSpaceFromResponse(config, response), nil
}
