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

func (loader SpaceLoader) Load(guid TypedGUID, token string) (string, string, error) {
    if !guid.BelongsToSpace() {
        return "", "", nil
    }

    space, err := loader.cloudController.LoadSpace(guid.String(), token)
    if err != nil {
        return "", "", CCErrorFor(err)
    }

    org, err := loader.cloudController.LoadOrganization(space.OrganizationGuid, token)
    if err != nil {
        return "", "", CCErrorFor(err)
    }

    return space.Name, org.Name, nil
}
