package models_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UnsubscribersRepo", func() {
	var (
		repo          models.UnsubscribersRepository
		conn          db.ConnectionInterface
		guidGenerator *mocks.GUIDGenerator
	)

	BeforeEach(func() {
		database := db.NewDatabase(sqlDB, db.Config{})
		helpers.TruncateTables(database)

		guidGenerator = mocks.NewGUIDGenerator()
		guidGenerator.GenerateCall.Returns.GUIDs = []string{"first-random-guid", "second-random-guid"}

		repo = models.NewUnsubscribersRepository(guidGenerator.Generate)
		conn = database.Connection()
	})

	Describe("Insert", func() {
		It("returns the inserted record", func() {
			createdUnsubscriber, err := repo.Insert(conn, models.Unsubscriber{
				CampaignTypeID: "some-campaign-type-id",
				UserGUID:       "some-user-guid",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(createdUnsubscriber.ID).To(Equal("first-random-guid"))
			Expect(createdUnsubscriber.CampaignTypeID).To(Equal("some-campaign-type-id"))
			Expect(createdUnsubscriber.UserGUID).To(Equal("some-user-guid"))
		})

		Context("when an error occurs", func() {
			Context("when the guid generator errors", func() {
				It("returns an error", func() {
					guidGenerator.GenerateCall.Returns.Error = errors.New("some-guid-error")

					_, err := repo.Insert(conn, models.Unsubscriber{
						CampaignTypeID: "some-campaign-type-id",
						UserGUID:       "some-user-guid",
					})
					Expect(err).To(MatchError(errors.New("some-guid-error")))
				})
			})

			Context("when inserting a database record errors", func() {
				It("returns an error", func() {
					connection := mocks.NewConnection()
					connection.InsertCall.Returns.Error = errors.New("some other error")

					_, err := repo.Insert(connection, models.Unsubscriber{
						CampaignTypeID: "some-campaign-type-id",
						UserGUID:       "some-user-guid",
					})
					Expect(err).To(MatchError(errors.New("some other error")))
				})
			})
		})
	})

	Describe("Get", func() {
		Context("when an unsubscriber record exists with the given campaign_type_id and user_guid", func() {
			It("returns the unsubscriber", func() {
				createdUnsubscriber, err := repo.Insert(conn, models.Unsubscriber{
					CampaignTypeID: "some-campaign-type-id",
					UserGUID:       "some-user-guid",
				})
				Expect(err).NotTo(HaveOccurred())

				gottenUnsubscriber, err := repo.Get(conn, "some-user-guid", "some-campaign-type-id")
				Expect(err).NotTo(HaveOccurred())
				Expect(gottenUnsubscriber).To(Equal(createdUnsubscriber))
			})
		})

		Context("when an error occurs", func() {
			Context("when an unsubscriber does not exist", func() {
				It("returns a RecordNotFound error", func() {
					_, err := repo.Get(conn, "some-user-guid", "some-campaign-type-id")
					Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError{}))
				})
			})

			Context("when an unknown error happens", func() {
				It("returns the error", func() {
					connection := mocks.NewConnection()
					connection.SelectOneCall.Returns.Error = errors.New("some other error")

					_, err := repo.Get(connection, "some-user-guid", "some-campaign-type-id")
					Expect(err).To(MatchError(errors.New("some other error")))
				})
			})
		})
	})

	Describe("Delete", func() {
		It("deletes the specified record", func() {
			_, err := repo.Insert(conn, models.Unsubscriber{
				CampaignTypeID: "some-campaign-type-id",
				UserGUID:       "some-user-guid",
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.Get(conn, "some-user-guid", "some-campaign-type-id")
			Expect(err).NotTo(HaveOccurred())

			err = repo.Delete(conn, models.Unsubscriber{
				CampaignTypeID: "some-campaign-type-id",
				UserGUID:       "some-user-guid",
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.Get(conn, "some-user-guid", "some-campaign-type-id")
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError{}))
		})

		It("does not return an error if the user is not unsubscribed", func() {
			err := repo.Delete(conn, models.Unsubscriber{
				CampaignTypeID: "some-campaign-type-id",
				UserGUID:       "some-user-guid",
			})
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when an unknown error happens", func() {
			It("returns the error", func() {
				connection := mocks.NewConnection()
				connection.ExecCall.Returns.Error = errors.New("some other error")

				err := repo.Delete(connection, models.Unsubscriber{
					CampaignTypeID: "some-campaign-type-id",
					UserGUID:       "some-user-guid",
				})
				Expect(err).To(MatchError(errors.New("some other error")))
			})
		})
	})
})
