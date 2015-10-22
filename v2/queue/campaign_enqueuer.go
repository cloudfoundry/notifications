package queue

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
)

type enqueuer interface {
	Enqueue(job *gobble.Job, conn gobble.ConnectionInterface) (*gobble.Job, error)
}

type DatabaseInterface interface {
	collections.DatabaseInterface
}

type CampaignJob struct {
	JobType  string
	Campaign collections.Campaign
}

type CampaignEnqueuer struct {
	gobbleQueue       enqueuer
	gobbleInitializer gobbleInitializer
	database          DatabaseInterface
}

func NewCampaignEnqueuer(queue enqueuer, database DatabaseInterface, gobbleInitializer gobbleInitializer) CampaignEnqueuer {
	return CampaignEnqueuer{
		gobbleQueue:       queue,
		database:          database,
		gobbleInitializer: gobbleInitializer,
	}
}

func (e CampaignEnqueuer) Enqueue(campaign collections.Campaign, jobType string) error {
	connection := e.database.Connection()
	e.gobbleInitializer.InitializeDBMap(connection.GetDbMap())
	job := gobble.NewJob(CampaignJob{
		JobType:  jobType,
		Campaign: campaign,
	})

	_, err := e.gobbleQueue.Enqueue(job, connection)
	if err != nil {
		return errors.New(fmt.Sprintf("there was an error enqueuing the job: %s", err))
	}

	return nil
}
