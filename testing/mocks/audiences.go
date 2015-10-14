package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/horde"

type Audiences struct {
	GenerateAudiencesCall struct {
		Receives struct {
			Inputs []string
		}
		Returns struct {
			Audiences []horde.Audience
			Error     error
		}
	}
}

func NewAudiences() *Audiences {
	return &Audiences{}
}

func (a *Audiences) GenerateAudiences(inputs []string) ([]horde.Audience, error) {
	a.GenerateAudiencesCall.Receives.Inputs = inputs

	return a.GenerateAudiencesCall.Returns.Audiences, a.GenerateAudiencesCall.Returns.Error
}
