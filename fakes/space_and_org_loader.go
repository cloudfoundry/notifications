package fakes

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/postal"
)

type SpaceAndOrgLoader struct {
    LoadError    error
    Space        cf.CloudControllerSpace
    Organization cf.CloudControllerOrganization
}

func NewSpaceAndOrgLoader() *SpaceAndOrgLoader {
    return &SpaceAndOrgLoader{}
}

func (fake *SpaceAndOrgLoader) Load(postal.TypedGUID, string) (cf.CloudControllerSpace, cf.CloudControllerOrganization, error) {
    return fake.Space, fake.Organization, fake.LoadError
}
