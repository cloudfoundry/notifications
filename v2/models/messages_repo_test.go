package models_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessagesRepository", func() {
	var (
		repo          models.MessagesRepository
		conn          *db.Connection
		clock         *mocks.Clock
		guidGenerator *mocks.IDGenerator
	)

	BeforeEach(func() {
		database := models.NewDatabase(sqlDB, models.Config{})
		helpers.TruncateTables(db.NewDatabase(sqlDB, db.Config{}))
		conn = database.Connection().(*db.Connection)
		conn.AddTableWithName(models.Message{}, "messages")
		clock = mocks.NewClock()
		guidGenerator = mocks.NewIDGenerator()
		guidGenerator.GenerateCall.Returns.IDs = []string{"random-guid-1"}

		repo = models.NewMessagesRepository(clock, guidGenerator.Generate)
	})

	Describe("CountByStatus", func() {
		BeforeEach(func() {
			err := conn.Insert(&models.Message{
				ID:         "random-guid-1",
				CampaignID: "some-campaign-id",
				Status:     common.StatusDelivered,
				UpdatedAt:  time.Now().UTC().Truncate(time.Second),
			})
			Expect(err).NotTo(HaveOccurred())

			err = conn.Insert(&models.Message{
				ID:         "random-guid-2",
				CampaignID: "some-campaign-id",
				Status:     common.StatusFailed,
				UpdatedAt:  time.Now().UTC().Truncate(time.Second),
			})
			Expect(err).NotTo(HaveOccurred())

			err = conn.Insert(&models.Message{
				ID:         "random-guid-3",
				CampaignID: "some-campaign-id",
				Status:     common.StatusDelivered,
				UpdatedAt:  time.Now().UTC().Truncate(time.Second),
			})
			Expect(err).NotTo(HaveOccurred())

			err = conn.Insert(&models.Message{
				ID:         "random-guid-4",
				CampaignID: "some-campaign-id",
				Status:     common.StatusRetry,
				UpdatedAt:  time.Now().UTC().Truncate(time.Second),
			})
			Expect(err).NotTo(HaveOccurred())

			err = conn.Insert(&models.Message{
				ID:         "random-guid-5",
				CampaignID: "some-campaign-id",
				Status:     common.StatusQueued,
				UpdatedAt:  time.Now().UTC().Truncate(time.Second),
			})
			Expect(err).NotTo(HaveOccurred())

			err = conn.Insert(&models.Message{
				ID:         "random-guid-6",
				CampaignID: "some-campaign-id",
				Status:     common.StatusUndeliverable,
				UpdatedAt:  time.Now().UTC().Truncate(time.Second),
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return the counts of each message status", func() {
			messageCounts, err := repo.CountByStatus(conn, "some-campaign-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(messageCounts).To(Equal(models.MessageCounts{
				Total:         6,
				Retry:         1,
				Failed:        1,
				Delivered:     2,
				Queued:        1,
				Undeliverable: 1,
			}))
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

	Describe("MostRecentlyUpdatedByCampaignID", func() {
		var anotherUpdatedAt time.Time
		BeforeEach(func() {
			updatedAt, err := time.Parse(time.RFC3339, "2014-12-31T12:05:05+07:00")
			Expect(err).NotTo(HaveOccurred())
			updatedAt = updatedAt.UTC()

			err = conn.Insert(&models.Message{
				ID:         "random-guid-1",
				CampaignID: "some-campaign-id",
				Status:     common.StatusDelivered,
				UpdatedAt:  updatedAt,
			})
			Expect(err).NotTo(HaveOccurred())

			anotherUpdatedAt, err = time.Parse(time.RFC3339, "2014-12-31T12:06:05+07:00")
			Expect(err).NotTo(HaveOccurred())
			anotherUpdatedAt = anotherUpdatedAt.UTC()

			err = conn.Insert(&models.Message{
				ID:         "random-guid-2",
				CampaignID: "some-campaign-id",
				Status:     common.StatusFailed,
				UpdatedAt:  anotherUpdatedAt,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return the most recently updated message", func() {
			message, err := repo.MostRecentlyUpdatedByCampaignID(conn, "some-campaign-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(message).To(Equal(models.Message{
				ID:         "random-guid-2",
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
			clock.NowCall.Returns.Time = time.Now().UTC().Truncate(time.Second)
			message, err := repo.Insert(conn, models.Message{
				Status:     "some-status",
				CampaignID: "some-campaign-id",
			})
			Expect(err).NotTo(HaveOccurred())

			var msg models.Message
			err = conn.SelectOne(&msg, "SELECT * FROM `messages` WHERE `id` = ? AND `status` = ? AND `campaign_id` = ?", "random-guid-1", "some-status", "some-campaign-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(msg).To(Equal(message))
		})

		Context("when an error occurs", func() {
			It("returns an error", func() {
				connection := mocks.NewConnection()
				connection.InsertCall.Returns.Error = errors.New("some connection error")

				_, err := repo.Insert(connection, models.Message{
					Status:     "some-status",
					CampaignID: "some-campaign-id",
				})
				Expect(err).To(MatchError(errors.New("some connection error")))
			})

			It("returns an error when the guid generator errors", func() {
				guidGenerator.GenerateCall.Returns.Error = errors.New("some guid error")

				_, err := repo.Insert(conn, models.Message{
					Status:     "some-status",
					CampaignID: "some-campaign-id",
				})
				Expect(err).To(MatchError(errors.New("some guid error")))
			})
		})
	})

	Describe("Update", func() {
		var updatedAt time.Time

		BeforeEach(func() {
			err := conn.Insert(&models.Message{
				ID:         "random-guid-1",
				Status:     "some-status",
				CampaignID: "some-campaign-id",
				UpdatedAt:  time.Now().Add(-30 * time.Second).UTC().Truncate(time.Second),
			})
			Expect(err).NotTo(HaveOccurred())

			updatedAt = time.Now().Truncate(time.Second).UTC()
			clock.NowCall.Returns.Time = updatedAt
		})

		It("updates an existing message in the database", func() {
			_, err := repo.Update(conn, models.Message{
				ID:         "random-guid-1",
				Status:     "some-new-status",
				CampaignID: "some-campaign-id",
			})

			var msg models.Message
			err = conn.SelectOne(&msg, "SELECT * FROM `messages` WHERE `id` = ? AND `status` = ? AND `campaign_id` = ?", "random-guid-1", "some-new-status", "some-campaign-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(msg).To(Equal(models.Message{
				ID:         "random-guid-1",
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
					Status:     "some-status",
					CampaignID: "some-campaign-id",
				})
				Expect(err).To(MatchError(errors.New("some update error")))
			})
		})
	})
})
