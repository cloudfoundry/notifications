package mocks

import "github.com/pivotal-cf-experimental/rainmaker"

type RainmakerOrganizationsService struct {
	GetCall struct {
		Receives struct {
			GUID  string
			Token string
		}
		Returns struct {
			Organization rainmaker.Organization
			Error        error
		}
	}
}

func NewRainmakerOrganizationsService() *RainmakerOrganizationsService {
	return &RainmakerOrganizationsService{}
}

func (s *RainmakerOrganizationsService) Get(guid, token string) (rainmaker.Organization, error) {
	s.GetCall.Receives.GUID = guid
	s.GetCall.Receives.Token = token

	return s.GetCall.Returns.Organization, s.GetCall.Returns.Error
}
