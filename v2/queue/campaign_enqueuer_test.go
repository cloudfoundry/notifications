package queue_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Enqueuer", func() {
	var (
		gobble   *fakes.Queue
		enqueuer queue.CampaignEnqueuer
		campaign collections.Campaign
	)

	BeforeEach(func() {
		gobble = fakes.NewQueue()
		enqueuer = queue.NewCampaignEnqueuer(gobble)
		campaign = collections.Campaign{
			ID: "27",
		}
	})

	Context("Enqueue", func() {
		It("puts a campaign on the queue", func() {
			err := enqueuer.Enqueue(campaign, "campaign")
			Expect(err).NotTo(HaveOccurred())
			Expect(gobble.Len()).To(Equal(1))

			var job queue.CampaignJob
			gobble.Jobs[1].Unmarshal(&job)
			Expect(job.Campaign.ID).To(Equal("27"))
			Expect(job.JobType).To(Equal("campaign"))
		})

		Context("when an enqueuing occurs", func() {
			BeforeEach(func() {
				gobble.EnqueueError = errors.New("some-error")
			})

			It("returns an error", func() {
				err := enqueuer.Enqueue(campaign, "campaign")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("there was an error enqueuing the job: some-error"))
				Expect(gobble.Len()).To(Equal(0))
			})
		})
	})
})
