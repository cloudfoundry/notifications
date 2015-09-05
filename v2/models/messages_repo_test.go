package models_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Messages Repository", func() {
	var (
		repo  models.MessagesRepository
		conn  *db.Connection
		clock *mocks.Clock
	)

	BeforeEach(func() {
		database := models.NewDatabase(sqlDB, models.Config{})
		helpers.TruncateTables(db.NewDatabase(sqlDB, db.Config{}))
		conn = database.Connection().(*db.Connection)
		conn.AddTableWithName(models.Message{}, "messages")
		clock = mocks.NewClock()

		repo = models.NewMessagesRepository(clock)
	})

	Describe("CountByStatus", func() {
		BeforeEach(func() {
			err := conn.Insert(&models.Message{
				ID:         "message-id-123",
				CampaignID: "some-campaign-id",
				Status:     postal.StatusDelivered,
			})
			Expect(err).NotTo(HaveOccurred())

			err = conn.Insert(&models.Message{
				ID:         "message-id-1234",
				CampaignID: "some-campaign-id",
				Status:     postal.StatusFailed,
			})
			Expect(err).NotTo(HaveOccurred())

			err = conn.Insert(&models.Message{
				ID:         "message-id-456",
				CampaignID: "some-campaign-id",
				Status:     postal.StatusDelivered,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return the counts of each message status", func() {
			messageCounts, err := repo.CountByStatus(conn, "some-campaign-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(messageCounts.Failed).To(Equal(1))
			Expect(messageCounts.Delivered).To(Equal(2))
		})

		Context("when an error occurs", func() {
			It("should return an error", func() {
				connection := mocks.NewConnection()
				connection.SelectCall.Returns.Error = errors.New("some connection error")

				_, err := repo.CountByStatus(connection, "some-campaign-id")
				Expect(err).To(MatchError(errors.New("some connection error")))
			})
		})
	})

	Describe("MostRecentlyUpatedByCampaignID", func() {
		var anotherUpdatedAt time.Time
		BeforeEach(func() {
			var err error

			updatedAt, err := time.Parse(time.RFC3339, "2014-12-31T12:05:05+07:00")
			Expect(err).NotTo(HaveOccurred())
			updatedAt = updatedAt.UTC()

			err = conn.Insert(&models.Message{
				ID:         "message-id-123",
				CampaignID: "some-campaign-id",
				Status:     postal.StatusDelivered,
				UpdatedAt:  updatedAt,
			})
			Expect(err).NotTo(HaveOccurred())

			anotherUpdatedAt, err = time.Parse(time.RFC3339, "2014-12-31T12:06:05+07:00")
			Expect(err).NotTo(HaveOccurred())
			anotherUpdatedAt = anotherUpdatedAt.UTC()

			err = conn.Insert(&models.Message{
				ID:         "message-id-1234",
				CampaignID: "some-campaign-id",
				Status:     postal.StatusFailed,
				UpdatedAt:  anotherUpdatedAt,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return the most recently updated message", func() {
			message, err := repo.MostRecentlyUpdatedByCampaignID(conn, "some-campaign-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(message).To(Equal(models.Message{
				ID:         "message-id-1234",
				CampaignID: "some-campaign-id",
				Status:     "failed",
				UpdatedAt:  anotherUpdatedAt,
			}))
		})

		Context("when an error occurs", func() {
			It("returns an error", func() {
				connection := mocks.NewConnection()
				connection.SelectOneCall.Returns.Error = errors.New("some connection error")

				_, err := repo.MostRecentlyUpdatedByCampaignID(connection, "some-campaign-id")
				Expect(err).To(MatchError(errors.New("some connection error")))
			})
		})
	})

	Describe("Insert", func() {
		It("inserts a message into the database table", func() {
			message, err := repo.Insert(conn, models.Message{
				ID:         "some-message-id",
				Status:     "some-status",
				CampaignID: "some-campaign-id",
			})

			var msg models.Message
			err = conn.SelectOne(&msg, "SELECT * FROM `messages` WHERE `id` = ? AND `status` = ? AND `campaign_id` = ?", "some-message-id", "some-status", "some-campaign-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(msg).To(Equal(message))
		})

		Context("when an error occurs", func() {
			It("returns an error", func() {
				connection := mocks.NewConnection()
				connection.InsertCall.Returns.Error = errors.New("some connection error")

				_, err := repo.Insert(connection, models.Message{
					ID:         "some-message-id",
					Status:     "some-status",
					CampaignID: "some-campaign-id",
				})
				Expect(err).To(MatchError(errors.New("some connection error")))
			})
		})
	})

	Describe("Update", func() {
		var updatedAt time.Time

		BeforeEach(func() {
			err := conn.Insert(&models.Message{
				ID:         "some-message-id",
				Status:     "some-status",
				CampaignID: "some-campaign-id",
			})
			Expect(err).NotTo(HaveOccurred())

			updatedAt = time.Now().Truncate(time.Second).UTC()
			clock.NowCall.Returns.Time = updatedAt
		})

		It("updates an existing message in the database", func() {
			_, err := repo.Update(conn, models.Message{
				ID:         "some-message-id",
				Status:     "some-new-status",
				CampaignID: "some-campaign-id",
			})

			var msg models.Message
			err = conn.SelectOne(&msg, "SELECT * FROM `messages` WHERE `id` = ? AND `status` = ? AND `campaign_id` = ?", "some-message-id", "some-new-status", "some-campaign-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(msg).To(Equal(models.Message{
				ID:         "some-message-id",
				Status:     "some-new-status",
				CampaignID: "some-campaign-id",
				UpdatedAt:  updatedAt,
			}))
		})

		Context("when an error occurs", func() {
			It("returns an error", func() {
				connection := mocks.NewConnection()
				connection.UpdateCall.Returns.Error = errors.New("some update error")

				_, err := repo.Update(connection, models.Message{
					ID:         "some-message-id",
					Status:     "some-status",
					CampaignID: "some-campaign-id",
				})
				Expect(err).To(MatchError(errors.New("some update error")))
			})
		})
	})
})
