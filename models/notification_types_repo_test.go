package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationTypesRepo", func() {
	var (
		repo models.NotificationTypesRepository
		conn models.ConnectionInterface
	)

	BeforeEach(func() {
		TruncateTables()
		repo = models.NewNotificationTypesRepository(fakes.NewIncrementingGUIDGenerator().Generate)
		db := models.NewDatabase(sqlDB, models.Config{})
		db.Setup()
		conn = db.Connection()
	})

	Describe("Insert", func() {
		It("inserts the record into the database", func() {
			notificationType := models.NotificationType{
				Name:        "some-notification-type",
				Description: "some-notification-type-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			}

			returnNotificationType, err := repo.Insert(conn, notificationType)
			Expect(err).NotTo(HaveOccurred())
			Expect(returnNotificationType.ID).To(Equal("deadbeef-aabb-ccdd-eeff-001122334455"))
		})
	})

	Describe("GetBySenderIDAndName", func() {
		It("fetches the notification type given a sender_id and name", func() {
			createdNotificationType, err := repo.Insert(conn, models.NotificationType{
				Name:        "some-notification-type",
				Description: "some-notification-type-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())

			notificationType, err := repo.GetBySenderIDAndName(conn, "some-sender-id", "some-notification-type")
			Expect(err).NotTo(HaveOccurred())
			Expect(notificationType.ID).To(Equal(createdNotificationType.ID))
		})

		It("fails to fetch the notification type given a non-existent sender_id and name", func() {
			_, err := repo.GetBySenderIDAndName(conn, "another-sender-id", "some-notification-type")
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
		})
	})
})
