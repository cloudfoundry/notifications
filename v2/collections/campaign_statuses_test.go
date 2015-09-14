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
		sendersRepository          *mocks.SendersRepository
		messagesRepository         *mocks.MessagesRepository
		conn                       *mocks.Connection
		campaignStatusesCollection collections.CampaignStatusesCollection
	)

	BeforeEach(func() {
		campaignsRepository = mocks.NewCampaignsRepository()
		sendersRepository = mocks.NewSendersRepository()
		messagesRepository = mocks.NewMessagesRepository()
		conn = mocks.NewConnection()

		campaignStatusesCollection = collections.NewCampaignStatusesCollection(campaignsRepository, sendersRepository, messagesRepository)
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
				SenderID:  "sender-id",
			}

			sendersRepository.GetCall.Returns.Sender = models.Sender{
				ID: "sender-id",
			}
		})

		It("returns the status of the campaign", func() {
			campaignStatus, err := campaignStatusesCollection.Get(conn, "campaign-id", "sender-id")
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

			Expect(sendersRepository.GetCall.Receives.Connection).To(Equal(conn))
			Expect(sendersRepository.GetCall.Receives.SenderID).To(Equal("sender-id"))
		})

		Context("when the campaign is not yet completed", func() {
			It("returns a transient status", func() {
				messagesRepository.CountByStatusCall.Returns.MessageCounts = models.MessageCounts{
					Total:     5,
					Retry:     1,
					Failed:    1,
					Delivered: 1,
					Queued:    2,
				}

				messagesRepository.MostRecentlyUpdatedByCampaignIDCall.Returns.Error = errors.New("sql: no rows")

				campaignStatus, err := campaignStatusesCollection.Get(conn, "campaign-id", "sender-id")
				Expect(err).NotTo(HaveOccurred())

				Expect(campaignStatus).To(Equal(collections.CampaignStatus{
					CampaignID:     "campaign-id",
					Status:         "sending",
					TotalMessages:  5,
					SentMessages:   1,
					FailedMessages: 1,
					RetryMessages:  1,
					QueuedMessages: 2,
					StartTime:      startTime,
					CompletedTime:  mysql.NullTime{},
				}))
			})
		})

		Context("when the campaign has not yet been processed", func() {
			It("returns a transient status", func() {
				messagesRepository.CountByStatusCall.Returns.MessageCounts = models.MessageCounts{
					Total:     0,
					Retry:     0,
					Failed:    0,
					Delivered: 0,
				}

				messagesRepository.MostRecentlyUpdatedByCampaignIDCall.Returns.Error = errors.New("sql: no rows")

				campaignStatus, err := campaignStatusesCollection.Get(conn, "campaign-id", "sender-id")
				Expect(err).NotTo(HaveOccurred())

				Expect(campaignStatus).To(Equal(collections.CampaignStatus{
					CampaignID:     "campaign-id",
					Status:         "sending",
					TotalMessages:  0,
					SentMessages:   0,
					FailedMessages: 0,
					RetryMessages:  0,
					StartTime:      startTime,
					CompletedTime:  mysql.NullTime{},
				}))
			})
		})

		Context("failure cases", func() {
			It("returns an error when the campaign cannot be found", func() {
				notFoundError := models.RecordNotFoundError{errors.New("not found")}
				campaignsRepository.GetCall.Returns.Error = notFoundError

				_, err := campaignStatusesCollection.Get(conn, "campaign-id", "sender-id")
				Expect(err).To(MatchError(collections.NotFoundError{notFoundError}))
			})

			It("returns an error when the campaign repo errors", func() {
				campaignsRepository.GetCall.Returns.Error = errors.New("connection error")

				_, err := campaignStatusesCollection.Get(conn, "campaign-id", "sender-id")
				Expect(err).To(MatchError(collections.UnknownError{errors.New("connection error")}))
			})

			It("returns an error when the sender cannot be found", func() {
				notFoundError := models.RecordNotFoundError{errors.New("not found")}
				sendersRepository.GetCall.Returns.Error = notFoundError

				_, err := campaignStatusesCollection.Get(conn, "campaign-id", "sender-id")
				Expect(err).To(MatchError(collections.NotFoundError{notFoundError}))
			})

			It("returns an error when the sender repo errors", func() {
				sendersRepository.GetCall.Returns.Error = errors.New("connection error")

				_, err := campaignStatusesCollection.Get(conn, "campaign-id", "sender-id")
				Expect(err).To(MatchError(collections.UnknownError{errors.New("connection error")}))
			})

			It("returns an error when the campaign does not match the sender", func() {
				campaignsRepository.GetCall.Returns.Campaign = models.Campaign{
					ID:        "campaign-id",
					StartTime: startTime,
					SenderID:  "different-sender-id",
				}

				_, err := campaignStatusesCollection.Get(conn, "campaign-id", "sender-id")
				Expect(err).To(MatchError(collections.NotFoundError{errors.New("Campaign with id \"campaign-id\" could not be found")}))
			})

			It("returns an error when the messages repo count call errors", func() {
				messagesRepository.CountByStatusCall.Returns.Error = errors.New("something bad")

				_, err := campaignStatusesCollection.Get(conn, "campaign-id", "sender-id")
				Expect(err).To(MatchError(collections.UnknownError{errors.New("something bad")}))
			})

			It("returns an error when the messages repo updated at errors", func() {
				messagesRepository.MostRecentlyUpdatedByCampaignIDCall.Returns.Error = errors.New("db went away")

				_, err := campaignStatusesCollection.Get(conn, "campaign-id", "sender-id")
				Expect(err).To(MatchError(collections.UnknownError{errors.New("db went away")}))
			})
		})
	})
})
