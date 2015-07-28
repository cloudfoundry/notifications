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
		repo = models.NewCampaignTypesRepository(fakes.NewIncrementingGUIDGenerator().Generate)
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

	Describe("List", func() {
		It("fetches a list of records from the database", func() {
			createdNotificationTypeOne, err := repo.Insert(conn, models.NotificationType{
				Name:        "notification-type-one",
				Description: "notification-type-one-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())

			createdNotificationTypeTwo, err := repo.Insert(conn, models.NotificationType{
				Name:        "notification-type-two",
				Description: "notification-type-two-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())

			returnNotificationTypeList, err := repo.List(conn, "some-sender-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(len(returnNotificationTypeList)).To(Equal(2))

			Expect(returnNotificationTypeList[0].ID).To(Equal(createdNotificationTypeOne.ID))
			Expect(returnNotificationTypeList[0].SenderID).To(Equal(createdNotificationTypeOne.SenderID))

			Expect(returnNotificationTypeList[1].ID).To(Equal(createdNotificationTypeTwo.ID))
			Expect(returnNotificationTypeList[1].SenderID).To(Equal(createdNotificationTypeTwo.SenderID))
		})

		It("fetches an empty list of records from the database if nothing has been inserted", func() {
			returnNotificationTypeList, err := repo.List(conn, "some-sender-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(len(returnNotificationTypeList)).To(Equal(0))
		})
	})

	Describe("Get", func() {
		It("fetches a record from the database", func() {
			notificationType, err := repo.Insert(conn, models.NotificationType{
				Name:        "notification-type",
				Description: "notification-type-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())

			returnNotificationType, err := repo.Get(conn, notificationType.ID)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnNotificationType).To(Equal(notificationType))
		})

		Context("failure cases", func() {
			It("fails to fetch the notification type given a non-existent notification_type_id", func() {
				_, err := repo.Insert(conn, models.NotificationType{
					Name:        "notification-type",
					Description: "notification-type-description",
					Critical:    false,
					TemplateID:  "some-template-id",
					SenderID:    "some-sender-id",
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = repo.Get(conn, "missing-notification-type-id")
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})
		})
	})
})
