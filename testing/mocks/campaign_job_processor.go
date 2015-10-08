package mocks

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
)

type CampaignJobProcessor struct {
	ProcessCall struct {
		Receives struct {
			Connection services.ConnectionInterface
			UAAHost    string
			Job        gobble.Job
		}

		Returns struct {
			Error error
		}

		WasCalled bool
	}
}

func NewCampaignJobProcessor() *CampaignJobProcessor {
	return &CampaignJobProcessor{}
}

func (p *CampaignJobProcessor) Process(conn services.ConnectionInterface, uaaHost string, job gobble.Job) error {
	p.ProcessCall.Receives.Connection = conn
	p.ProcessCall.Receives.UAAHost = uaaHost
	p.ProcessCall.Receives.Job = job
	p.ProcessCall.WasCalled = true

	return p.ProcessCall.Returns.Error
}
