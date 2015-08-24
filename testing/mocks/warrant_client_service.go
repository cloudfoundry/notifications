package mocks

type WarrantClientService struct {
	GetTokenCall struct {
		Receives struct {
			ID     string
			Secret string
		}

		Returns struct {
			Token string
			Error error
		}
	}
}

func NewWarrantClientService() *WarrantClientService {
	return &WarrantClientService{}
}

func (s *WarrantClientService) GetToken(id, secret string) (string, error) {
	s.GetTokenCall.Receives.ID = id
	s.GetTokenCall.Receives.Secret = secret

	return s.GetTokenCall.Returns.Token, s.GetTokenCall.Returns.Error
}
