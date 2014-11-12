package utilities

import "github.com/cloudfoundry-incubator/notifications/cf"

type OrganizationLoader struct {
	cloudController cf.CloudControllerInterface
}

type OrganizationLoaderInterface interface {
	Load(string, string) (cf.CloudControllerOrganization, error)
}

func NewOrganizationLoader(cloudController cf.CloudControllerInterface) OrganizationLoader {
	return OrganizationLoader{
		cloudController: cloudController,
	}
}

func (loader OrganizationLoader) Load(organizationGUID string, token string) (cf.CloudControllerOrganization, error) {
	organization, err := loader.cloudController.LoadOrganization(organizationGUID, token)
	if err != nil {
		return cf.CloudControllerOrganization{}, CCErrorFor(err)
	}

	return organization, nil
}
