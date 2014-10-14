package postal

import "github.com/cloudfoundry-incubator/notifications/cf"

type SpaceLoader struct {
    cloudController cf.CloudControllerInterface
}

func NewSpaceLoader(cloudController cf.CloudControllerInterface) SpaceLoader {
    return SpaceLoader{
        cloudController: cloudController,
    }
}

func (loader SpaceLoader) Load(guid TypedGUID, token string) (cf.CloudControllerSpace, cf.CloudControllerOrganization, error) {
    if !guid.BelongsToSpace() {
        return loader.Error(nil)
    }

    space, err := loader.cloudController.LoadSpace(guid.String(), token)
    if err != nil {
        return loader.Error(CCErrorFor(err))
    }

    org, err := loader.cloudController.LoadOrganization(space.OrganizationGUID, token)
    if err != nil {
        return loader.Error(CCErrorFor(err))
    }

    return space, org, nil
}

func (loader SpaceLoader) Error(err error) (cf.CloudControllerSpace, cf.CloudControllerOrganization, error) {
    return cf.CloudControllerSpace{}, cf.CloudControllerOrganization{}, err
}
