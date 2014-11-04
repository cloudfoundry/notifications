package postal

import "github.com/cloudfoundry-incubator/notifications/cf"

type SpaceAndOrgLoader struct {
    cloudController cf.CloudControllerInterface
}

func NewSpaceAndOrgLoader(cloudController cf.CloudControllerInterface) SpaceAndOrgLoader {
    return SpaceAndOrgLoader{
        cloudController: cloudController,
    }
}

func (loader SpaceAndOrgLoader) Load(guid TypedGUID, token string) (cf.CloudControllerSpace, cf.CloudControllerOrganization, error) {
    var space cf.CloudControllerSpace
    var err error
    organizationGUID := guid.String()

    if !guid.BelongsToSpace() && !guid.BelongsToOrganization() {
        return loader.Error(nil)
    }

    if guid.BelongsToSpace() {
        space, err = loader.cloudController.LoadSpace(guid.String(), token)
        if err != nil {
            return loader.Error(CCErrorFor(err))
        }
        organizationGUID = space.OrganizationGUID
    }

    org, err := loader.cloudController.LoadOrganization(organizationGUID, token)
    if err != nil {
        return loader.Error(CCErrorFor(err))
    }

    return space, org, nil
}

func (loader SpaceAndOrgLoader) Error(err error) (cf.CloudControllerSpace, cf.CloudControllerOrganization, error) {
    return cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, err
}
