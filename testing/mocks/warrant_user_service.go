package mocks

import "github.com/pivotal-cf-experimental/warrant"

type WarrantUserService struct {
	GetCall struct {
		Receives struct {
			GUID  string
			Token string
		}

		Returns struct {
			User  warrant.User
			Error error
		}
	}
}

func NewWarrantUserService() *WarrantUserService {
	return &WarrantUserService{}
}

func (s *WarrantUserService) Get(guid, token string) (warrant.User, error) {
	s.GetCall.Receives.GUID = guid
	s.GetCall.Receives.Token = token
	return s.GetCall.Returns.User, s.GetCall.Returns.Error
}
