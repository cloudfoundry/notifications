package mocks

import "github.com/pivotal-cf-experimental/rainmaker"

type RainmakerSpacesService struct {
	GetCall struct {
		Receives struct {
			GUID  string
			Token string
		}
		Returns struct {
			Space rainmaker.Space
			Error error
		}
	}
}

func NewRainmakerSpacesService() *RainmakerSpacesService {
	return &RainmakerSpacesService{}
}

func (s *RainmakerSpacesService) Get(guid, token string) (rainmaker.Space, error) {
	s.GetCall.Receives.GUID = guid
	s.GetCall.Receives.Token = token

	return s.GetCall.Returns.Space, s.GetCall.Returns.Error
}
