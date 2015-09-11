package collections_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/go-sql-driver/mysql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignStatusesCollection", func() {
	var (
		campaignsRepository        *mocks.CampaignsRepository
		messagesRepository         *mocks.MessagesRepository
		conn                       *mocks.Connection
		campaignStatusesCollection collections.CampaignStatusesCollection
	)

	BeforeEach(func() {
		campaignsRepository = mocks.NewCampaignsRepository()
		messagesRepository = mocks.NewMessagesRepository()
		conn = mocks.NewConnection()

		campaignStatusesCollection = collections.NewCampaignStatusesCollection(campaignsRepository, messagesRepository)
	})

	Context("when a valid campaign is queried", func() {
		var startTime time.Time
		var updatedAtTime time.Time

		BeforeEach(func() {
			var err error

			messagesRepository.CountByStatusCall.Returns.MessageCounts = models.MessageCounts{
				Total:     2,
				Failed:    1,
				Delivered: 1,
			}

			updatedAtTime, err = time.Parse(time.RFC3339, "2015-09-01T12:45:56-07:00")
			Expect(err).NotTo(HaveOccurred())
			updatedAtTime = updatedAtTime.UTC()

			messagesRepository.MostRecentlyUpdatedByCampaignIDCall.Returns.Message = models.Message{
				CampaignID: "campaign-id",
				UpdatedAt:  updatedAtTime,
			}

			startTime, err = time.Parse(time.RFC3339, "2015-09-01T12:34:56-07:00")
			Expect(err).NotTo(HaveOccurred())
			startTime = startTime.UTC()

			campaignsRepository.GetCall.Returns.Campaign = models.Campaign{
				ID:        "campaign-id",
				StartTime: startTime,
			}
		})

		It("returns the status of the campaign", func() {
			campaignStatus, err := campaignStatusesCollection.Get(conn, "campaign-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(messagesRepository.CountByStatusCall.Receives.Connection).To(Equal(conn))
			Expect(messagesRepository.CountByStatusCall.Receives.CampaignIDList[0]).To(Equal("campaign-id"))

			Expect(messagesRepository.MostRecentlyUpdatedByCampaignIDCall.Receives.Connection).To(Equal(conn))
			Expect(messagesRepository.MostRecentlyUpdatedByCampaignIDCall.Receives.CampaignID).To(Equal("campaign-id"))

			Expect(campaignStatus).To(Equal(collections.CampaignStatus{
				CampaignID:     "campaign-id",
				Status:         "completed",
				TotalMessages:  2,
				SentMessages:   1,
				FailedMessages: 1,
				StartTime:      startTime,
				CompletedTime: mysql.NullTime{
					Time:  updatedAtTime,
					Valid: true,
				},
			}))

			Expect(campaignsRepository.GetCall.Receives.Connection).To(Equal(conn))
			Expect(campaignsRepository.GetCall.Receives.CampaignID).To(Equal("campaign-id"))
		})

		Context("when the campaign is not yet completed", func() {
			It("returns a transient status", func() {
				messagesRepository.CountByStatusCall.Returns.MessageCounts = models.MessageCounts{
					Total:     5,
					Retry:     1,
					Failed:    1,
					Delivered: 1,
				}

				messagesRepository.MostRecentlyUpdatedByCampaignIDCall.Returns.Error = errors.New("sql: no rows")

				campaignStatus, err := campaignStatusesCollection.Get(conn, "campaign-id")
				Expect(err).NotTo(HaveOccurred())

				Expect(campaignStatus).To(Equal(collections.CampaignStatus{
					CampaignID:     "campaign-id",
					Status:         "sending",
					TotalMessages:  5,
					SentMessages:   1,
					FailedMessages: 1,
					RetryMessages:  1,
					StartTime:      startTime,
					CompletedTime:  mysql.NullTime{},
				}))
			})
		})
	})
})
