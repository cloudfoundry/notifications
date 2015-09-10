package postal_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Campaign Status Updater", func() {
	var (
		campaignStatusUpdater postal.CampaignStatusUpdater
		messagesRepo          *mocks.MessagesRepository
		campaignsRepo         *mocks.CampaignsRepository
		pollInterval          time.Duration
		database              *mocks.Database
	)

	BeforeEach(func() {
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = mocks.NewConnection()
		messagesRepo = mocks.NewMessagesRepository()
		campaignsRepo = mocks.NewCampaignsRepository()
		pollInterval = 1 * time.Millisecond

		campaignStatusUpdater = postal.NewCampaignStatusUpdater(database.Connection(),
			messagesRepo,
			campaignsRepo,
			pollInterval)
	})

	Describe("Run", func() {
		It("starts the worker to gather campaign ids", func() {
			campaignStatusUpdater.Run()

			Eventually(func() int {
				return len(campaignsRepo.ListSendingCampaignsCall.Invocations)
			}).Should(BeNumerically(">=", 2))

			call1 := campaignsRepo.ListSendingCampaignsCall.Invocations[0]
			call2 := campaignsRepo.ListSendingCampaignsCall.Invocations[1]
			Expect(call2).To(BeTemporally(">", call1.Add(pollInterval)))
		})
	})

	Describe("Update", func() {
		Context("when all the messages are final", func() {
			BeforeEach(func() {
				campaignsRepo.ListSendingCampaignsCall.Returns.Campaigns = []models.Campaign{
					{ID: "some-great-campaign-id"},
					{ID: "another-great-campaign-id"},
				}

				messagesRepo.CountByStatusCall.Returns.MessageCounts = models.MessageCounts{
					Failed:    1,
					Delivered: 2,
				}
			})

			It("marks the campaign as completed with its statistics", func() {
				campaignStatusUpdater.Update()

				Expect(campaignsRepo.ListSendingCampaignsCall.Invocations).To(HaveLen(1))
				Expect(campaignsRepo.ListSendingCampaignsCall.Receives.Connection).To(Equal(database.Connection()))

				Expect(messagesRepo.CountByStatusCall.Receives.Connection).To(Equal(database.Connection()))
				Expect(messagesRepo.CountByStatusCall.Receives.CampaignIDList).To(HaveLen(2))
				Expect(messagesRepo.CountByStatusCall.Receives.CampaignIDList).To(ConsistOf([]string{
					"some-great-campaign-id",
					"another-great-campaign-id",
				}))

				Expect(campaignsRepo.UpdateCall.Receives.Connection).To(Equal(database.Connection()))
				Expect(campaignsRepo.UpdateCall.Receives.CampaignList).To(ConsistOf(models.Campaign{
					ID:             "some-great-campaign-id",
					Status:         "completed",
					TotalMessages:  3,
					SentMessages:   2,
					FailedMessages: 1,
				}, models.Campaign{
					ID:             "another-great-campaign-id",
					Status:         "completed",
					TotalMessages:  3,
					SentMessages:   2,
					FailedMessages: 1,
				}))
			})
		})
	})
})
