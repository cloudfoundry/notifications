package models_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignTypesRepo", func() {
	var (
		repo models.CampaignTypesRepository
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
			campaignType := models.CampaignType{
				Name:        "some-campaign-type",
				Description: "some-campaign-type-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			}

			returnCampaignType, err := repo.Insert(conn, campaignType)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnCampaignType.ID).To(Equal("deadbeef-aabb-ccdd-eeff-001122334455"))
		})

		Context("failure cases", func() {
			It("passes along error messages from the database", func() {
				conn := fakes.NewConnection()
				conn.InsertCall.Err = errors.New("a useful database error message")

				campaignType := models.CampaignType{
					Name:        "some-campaign-type",
					Description: "some-campaign-type-description",
					Critical:    false,
					TemplateID:  "some-template-id",
					SenderID:    "some-sender-id",
				}

				_, err := repo.Insert(conn, campaignType)
				Expect(err).To(MatchError("a useful database error message"))
			})
		})
	})

	Describe("GetBySenderIDAndName", func() {
		It("fetches the campaign type given a sender_id and name", func() {
			createdCampaignType, err := repo.Insert(conn, models.CampaignType{
				Name:        "some-campaign-type",
				Description: "some-campaign-type-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())

			campaignType, err := repo.GetBySenderIDAndName(conn, "some-sender-id", "some-campaign-type")
			Expect(err).NotTo(HaveOccurred())

			Expect(campaignType.ID).To(Equal(createdCampaignType.ID))
		})

		Context("failure cases", func() {
			It("fails to fetch the campaign type given a non-existent sender_id and name", func() {
				_, err := repo.GetBySenderIDAndName(conn, "another-sender-id", "some-campaign-type")
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})

			It("passes along error messages from the database", func() {
				conn := fakes.NewConnection()
				conn.SelectOneCall.Errs = []error{errors.New("a useful database error message")}
				_, err := repo.GetBySenderIDAndName(conn, "some-sender-id", "some-campaign-type")
				Expect(err).To(MatchError("a useful database error message"))
			})
		})
	})

	Describe("List", func() {
		It("fetches a list of records from the database", func() {
			createdCampaignTypeOne, err := repo.Insert(conn, models.CampaignType{
				Name:        "campaign-type-one",
				Description: "campaign-type-one-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())

			createdCampaignTypeTwo, err := repo.Insert(conn, models.CampaignType{
				Name:        "campaign-type-two",
				Description: "campaign-type-two-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())

			returnCampaignTypeList, err := repo.List(conn, "some-sender-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(len(returnCampaignTypeList)).To(Equal(2))

			Expect(returnCampaignTypeList[0].ID).To(Equal(createdCampaignTypeOne.ID))
			Expect(returnCampaignTypeList[0].SenderID).To(Equal(createdCampaignTypeOne.SenderID))

			Expect(returnCampaignTypeList[1].ID).To(Equal(createdCampaignTypeTwo.ID))
			Expect(returnCampaignTypeList[1].SenderID).To(Equal(createdCampaignTypeTwo.SenderID))
		})

		It("fetches an empty list of records from the database if nothing has been inserted", func() {
			returnCampaignTypeList, err := repo.List(conn, "some-sender-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(len(returnCampaignTypeList)).To(Equal(0))
		})

		Context("failure cases", func() {
			It("returns errors", func() {
				conn := fakes.NewConnection()
				conn.SelectCall.Err = errors.New("BOOM!")
				_, err := repo.List(conn, "some-sender-id")
				Expect(err).To(MatchError("BOOM!"))
			})
		})
	})

	Describe("Get", func() {
		It("fetches a record from the database", func() {
			campaignType, err := repo.Insert(conn, models.CampaignType{
				Name:        "campaign-type",
				Description: "campaign-type-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())

			returnCampaignType, err := repo.Get(conn, campaignType.ID)
			Expect(err).NotTo(HaveOccurred())

			Expect(returnCampaignType).To(Equal(campaignType))
		})

		Context("failure cases", func() {
			It("fails to fetch the campaign type given a non-existent campaign_type_id", func() {
				_, err := repo.Insert(conn, models.CampaignType{
					Name:        "campaign-type",
					Description: "campaign-type-description",
					Critical:    false,
					TemplateID:  "some-template-id",
					SenderID:    "some-sender-id",
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = repo.Get(conn, "missing-campaign-type-id")
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})

			It("passes along error messages from the database", func() {
				conn := fakes.NewConnection()
				conn.SelectOneCall.Errs = []error{errors.New("a useful database error message")}
				_, err := repo.Get(conn, "campaign-type-id")
				Expect(err).To(MatchError("a useful database error message"))
			})

		})
	})

	Describe("Update", func() {
		It("Updates a campaign type in the database", func() {
			existingCampaignType, err := repo.Insert(conn, models.CampaignType{
				Name:        "campaign-type",
				Description: "campaign-type-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())

			returnCampaignType, err := repo.Update(conn, models.CampaignType{
				ID:          existingCampaignType.ID,
				Name:        "updated name",
				Description: "updated description",
				Critical:    false,
				TemplateID:  "updated template id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(returnCampaignType.ID).To(Equal(existingCampaignType.ID))
			Expect(returnCampaignType.Name).To(Equal("updated name"))

			campaignType, err := repo.Get(conn, existingCampaignType.ID)
			Expect(campaignType.Name).To(Equal("updated name"))
			Expect(campaignType.Description).To(Equal("updated description"))
			Expect(campaignType.Critical).To(Equal(false))
			Expect(campaignType.TemplateID).To(Equal("updated template id"))
		})

		Context("failure cases", func() {
			It("provides a helpful error when no records are updated", func() {
				_, err := repo.Update(conn, models.CampaignType{
					ID:          "I-do-not-exist",
					Name:        "updated name",
					Description: "updated description",
					Critical:    false,
					TemplateID:  "updated template id",
					SenderID:    "some-sender-id",
				})
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
				Expect(string(err.(models.RecordNotFoundError))).To(Equal("No records updated: Campaign type with id \"I-do-not-exist\" could not be found"))
			})

			It("passes along error messagers from the database", func() {
				conn := fakes.NewConnection()
				conn.UpdateCall.Err = errors.New("a database error")

				campaignType := models.CampaignType{
					ID:          "some id",
					Name:        "updated name",
					Description: "updated description",
					Critical:    false,
					TemplateID:  "updated template id",
					SenderID:    "some-sender-id",
				}

				_, err := repo.Update(conn, campaignType)
				Expect(err).To(MatchError("a database error"))
			})
		})
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			_, err := repo.Insert(conn, models.CampaignType{
				ID:          "my-campaign-id",
				Name:        "campaign-type",
				Description: "campaign-type-description",
				Critical:    false,
				TemplateID:  "some-template-id",
				SenderID:    "some-sender-id",
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("deletes the campaignType from the database", func() {
			err := repo.Delete(conn, models.CampaignType{ID: "my-campaign-id"})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.Get(conn, "my-campaign-id")
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
		})

		Context("when an error occurs", func() {
			It("returns the error", func() {
				conn := fakes.NewConnection()
				databaseError := errors.New("The database is not valid")
				conn.DeleteCall.Err = databaseError
				err := repo.Delete(conn, models.CampaignType{ID: "other-campaign-id"})
				Expect(err).To(MatchError(databaseError))
			})
		})
	})
})
