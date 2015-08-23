package mocks

import "github.com/cloudfoundry-incubator/notifications/cf"

type SpaceLoader struct {
	LoadCall struct {
		Receives struct {
			SpaceGUID string
			Token     string
		}
		Returns struct {
			Space cf.CloudControllerSpace
			Error error
		}
	}
}

func NewSpaceLoader() *SpaceLoader {
	return &SpaceLoader{}
}

func (sl *SpaceLoader) Load(spaceGUID, token string) (cf.CloudControllerSpace, error) {
	sl.LoadCall.Receives.SpaceGUID = spaceGUID
	sl.LoadCall.Receives.Token = token

	return sl.LoadCall.Returns.Space, sl.LoadCall.Returns.Error
}
