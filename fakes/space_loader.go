package fakes

import "github.com/cloudfoundry-incubator/notifications/cf"

type SpaceLoader struct {
    LoadError error
    Space     cf.CloudControllerSpace
}

func NewSpaceLoader() *SpaceLoader {
    return &SpaceLoader{}
}

func (fake *SpaceLoader) Load(spaceGUID, token string) (cf.CloudControllerSpace, error) {
    return fake.Space, fake.LoadError
}
