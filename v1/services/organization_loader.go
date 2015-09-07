package services

import "github.com/cloudfoundry-incubator/notifications/cf"

type OrganizationLoader struct {
	cloudController cf.CloudControllerInterface
}

func NewOrganizationLoader(cloudController cf.CloudControllerInterface) OrganizationLoader {
	return OrganizationLoader{
		cloudController: cloudController,
	}
}

func (loader OrganizationLoader) Load(orgGUID string, token string) (cf.CloudControllerOrganization, error) {
	organization, err := loader.cloudController.LoadOrganization(orgGUID, token)
	if err != nil {
		return cf.CloudControllerOrganization{}, CCErrorFor(err)
	}

	return organization, nil
}
