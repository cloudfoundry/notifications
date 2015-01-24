package rainmaker

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pivotal-golang/rainmaker/internal/documents"
)

type User struct {
	GUID                           string
	URL                            string
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
	Admin                          bool
	Active                         bool
	DefaultSpaceGUID               string
	SpacesURL                      string
	OrganizationsURL               string
	ManagedOrganizationsURL        string
	BillingManagedOrganizationsURL string
	AuditedOrganizationsURL        string
	ManagedSpacesURL               string
	AuditedSpacesURL               string
}

func NewUserFromResponse(response documents.UserResponse) User {
	if response.Metadata.CreatedAt == nil {
		response.Metadata.CreatedAt = &time.Time{}
	}

	if response.Metadata.UpdatedAt == nil {
		response.Metadata.UpdatedAt = &time.Time{}
	}

	return User{
		GUID:                           response.Metadata.GUID,
		URL:                            response.Metadata.URL,
		CreatedAt:                      *response.Metadata.CreatedAt,
		UpdatedAt:                      *response.Metadata.UpdatedAt,
		Admin:                          response.Entity.Admin,
		Active:                         response.Entity.Active,
		DefaultSpaceGUID:               response.Entity.DefaultSpaceGUID,
		SpacesURL:                      response.Entity.SpacesURL,
		OrganizationsURL:               response.Entity.OrganizationsURL,
		ManagedOrganizationsURL:        response.Entity.ManagedOrganizationsURL,
		BillingManagedOrganizationsURL: response.Entity.BillingManagedOrganizationsURL,
		AuditedOrganizationsURL:        response.Entity.AuditedOrganizationsURL,
		ManagedSpacesURL:               response.Entity.ManagedSpacesURL,
		AuditedSpacesURL:               response.Entity.AuditedSpacesURL,
	}
}

func FetchUser(config Config, path, token string) (User, error) {
	_, body, err := NewClient(config).makeRequest(requestArguments{
		Method: "GET",
		Path:   path,
		Token:  token,
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return User{}, err
	}

	var response documents.UserResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	return NewUserFromResponse(response), nil
}
