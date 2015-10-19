package queue_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
	"github.com/go-gorp/gorp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignEnqueuer", func() {
	var (
		dbMap             *gorp.DbMap
		connection        *mocks.Connection
		gobbleQueue       *mocks.Queue
		gobbleInitializer *mocks.GobbleInitializer
		enqueuer          queue.CampaignEnqueuer
		campaign          collections.Campaign
	)

	BeforeEach(func() {
		gobbleQueue = mocks.NewQueue()
		gobbleInitializer = mocks.NewGobbleInitializer()

		dbMap = &gorp.DbMap{}
		connection = mocks.NewConnection()
		connection.GetDbMapCall.Returns.DbMap = dbMap
		database := mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = connection

		enqueuer = queue.NewCampaignEnqueuer(gobbleQueue, database, gobbleInitializer)
		campaign = collections.Campaign{
			ID: "27",
		}
	})

	Context("Enqueue", func() {
		It("puts a campaign on the queue", func() {
			err := enqueuer.Enqueue(campaign, "campaign")
			Expect(err).NotTo(HaveOccurred())

			Expect(gobbleQueue.EnqueueCall.Receives.Connection).To(Equal(connection))
			Expect(gobbleQueue.EnqueueCall.Receives.Jobs).To(HaveLen(1))
			Expect(gobbleQueue.EnqueueCall.Receives.Jobs[0]).To(Equal(gobble.NewJob(queue.CampaignJob{
				JobType:  "campaign",
				Campaign: campaign,
			})))

			isSamePtr := (gobbleInitializer.InitializeDBMapCall.Receives.DbMap == dbMap)
			Expect(isSamePtr).To(BeTrue())
		})

		Context("when an enqueuing occurs", func() {
			BeforeEach(func() {
				gobbleQueue.EnqueueCall.Returns.Error = errors.New("some-error")
			})

			It("returns an error", func() {
				err := enqueuer.Enqueue(campaign, "campaign")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("there was an error enqueuing the job: some-error"))
			})
		})
	})
})
