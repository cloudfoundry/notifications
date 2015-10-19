package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/v2/horde"
	"github.com/pivotal-golang/lager"
)

type Audiences struct {
	GenerateAudiencesCall struct {
		Receives struct {
			Inputs []string
			Logger lager.Logger
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

func (a *Audiences) GenerateAudiences(inputs []string, logger lager.Logger) ([]horde.Audience, error) {
	a.GenerateAudiencesCall.Receives.Inputs = inputs
	a.GenerateAudiencesCall.Receives.Logger = logger

	return a.GenerateAudiencesCall.Returns.Audiences, a.GenerateAudiencesCall.Returns.Error
}
