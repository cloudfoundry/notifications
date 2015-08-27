package mocks

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type CampaignEnqueuer struct {
	EnqueueCall struct {
		Receives struct {
			Campaign collections.Campaign
			JobType  string
		}
		Returns struct {
			Err error
		}
	}
}

func NewCampaignEnqueuer() *CampaignEnqueuer {
	return &CampaignEnqueuer{}
}

func (e *CampaignEnqueuer) Enqueue(campaign collections.Campaign, jobType string) error {
	e.EnqueueCall.Receives.Campaign = campaign
	e.EnqueueCall.Receives.JobType = jobType

	return e.EnqueueCall.Returns.Err
}
