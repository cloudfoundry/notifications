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

		It("fails to fetch the campaign type given a non-existent sender_id and name", func() {
			_, err := repo.GetBySenderIDAndName(conn, "another-sender-id", "some-campaign-type")
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
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
		})
	})
})
