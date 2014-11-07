package postal

import "github.com/cloudfoundry-incubator/notifications/cf"

type SpaceLoader struct {
    cloudController cf.CloudControllerInterface
}

type SpaceLoaderInterface interface {
    Load(string, string) (cf.CloudControllerSpace, error)
}

func NewSpaceLoader(cloudController cf.CloudControllerInterface) SpaceLoader {
    return SpaceLoader{
        cloudController: cloudController,
    }
}

func (loader SpaceLoader) Load(spaceGUID string, token string) (cf.CloudControllerSpace, error) {
    space, err := loader.cloudController.LoadSpace(spaceGUID, token)
    if err != nil {
        return cf.CloudControllerSpace{}, CCErrorFor(err)
    }

    return space, nil
}
