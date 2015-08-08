package fakes

import "github.com/cloudfoundry-incubator/notifications/cf"

type OrganizationLoader struct {
	LoadCall struct {
		Receives struct {
			OrganizationGUID string
			Token            string
		}
		Returns struct {
			Organization cf.CloudControllerOrganization
			Error        error
		}
	}
}

func NewOrganizationLoader() *OrganizationLoader {
	return &OrganizationLoader{}
}

func (ol *OrganizationLoader) Load(organizationGUID, token string) (cf.CloudControllerOrganization, error) {
	ol.LoadCall.Receives.OrganizationGUID = organizationGUID
	ol.LoadCall.Receives.Token = token

	return ol.LoadCall.Returns.Organization, ol.LoadCall.Returns.Error
}
