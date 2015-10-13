package queue

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type enqueuer interface {
	Enqueue(job gobble.Job) (gobble.Job, error)
}

type CampaignEnqueuer struct {
	gobbleQueue enqueuer
}

type CampaignJob struct {
	JobType  string
	Campaign collections.Campaign
}

func NewCampaignEnqueuer(queue enqueuer) CampaignEnqueuer {
	return CampaignEnqueuer{
		gobbleQueue: queue,
	}
}

func (e CampaignEnqueuer) Enqueue(campaign collections.Campaign, jobType string) error {
	_, err := e.gobbleQueue.Enqueue(gobble.NewJob(CampaignJob{
		JobType:  jobType,
		Campaign: campaign,
	}))
	if err != nil {
		return errors.New(fmt.Sprintf("there was an error enqueuing the job: %s", err))
	}

	return nil
}
