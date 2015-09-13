package services

import "github.com/cloudfoundry-incubator/notifications/cf"

type SpaceLoader struct {
	cc cloudController
}

func NewSpaceLoader(cc cloudController) SpaceLoader {
	return SpaceLoader{
		cc: cc,
	}
}

func (loader SpaceLoader) Load(spaceGUID string, token string) (cf.CloudControllerSpace, error) {
	space, err := loader.cc.LoadSpace(spaceGUID, token)
	if err != nil {
		return cf.CloudControllerSpace{}, CCErrorFor(err)
	}

	return space, nil
}
