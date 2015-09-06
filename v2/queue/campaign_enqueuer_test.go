package queue_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Enqueuer", func() {
	var (
		gobbleQueue *mocks.Queue
		enqueuer    queue.CampaignEnqueuer
		campaign    collections.Campaign
	)

	BeforeEach(func() {
		gobbleQueue = mocks.NewQueue()
		enqueuer = queue.NewCampaignEnqueuer(gobbleQueue)
		campaign = collections.Campaign{
			ID: "27",
		}
	})

	Context("Enqueue", func() {
		It("puts a campaign on the queue", func() {
			err := enqueuer.Enqueue(campaign, "campaign")
			Expect(err).NotTo(HaveOccurred())

			Expect(gobbleQueue.EnqueueCall.Receives.Jobs).To(HaveLen(1))
			Expect(gobbleQueue.EnqueueCall.Receives.Jobs[0]).To(Equal(gobble.NewJob(queue.CampaignJob{
				JobType:  "campaign",
				Campaign: campaign,
			})))
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
