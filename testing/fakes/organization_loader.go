package fakes

import "github.com/cloudfoundry-incubator/notifications/cf"

type OrganizationLoader struct {
	LoadError    error
	Organization cf.CloudControllerOrganization
}

func NewOrganizationLoader() *OrganizationLoader {
	return &OrganizationLoader{}
}

func (fake *OrganizationLoader) Load(organizationGUID, token string) (cf.CloudControllerOrganization, error) {
	return fake.Organization, fake.LoadError
}
