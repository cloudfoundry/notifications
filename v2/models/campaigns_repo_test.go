package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignsRepository", func() {
	var (
		repo          models.CampaignsRepository
		connection    db.ConnectionInterface
		guidGenerator *mocks.GUIDGenerator
	)

	BeforeEach(func() {
		guid1 := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55})
		guid2 := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x56})
		guidGenerator = mocks.NewGUIDGenerator()
		guidGenerator.GenerateCall.Returns.GUIDs = []*uuid.UUID{&guid1, &guid2}

		repo = models.NewCampaignsRepository(guidGenerator.Generate)
		database := db.NewDatabase(sqlDB, db.Config{})
		helpers.TruncateTables(database)
		connection = database.Connection()
	})

	Describe("Get", func() {
		It("gets a campaign from the database", func() {
			campaign, err := repo.Set(connection, models.Campaign{
				SendTo:         `{"user": "user-123"}`,
				CampaignTypeID: "some-campaign-type-id",
				Text:           "come see our new stuff",
				HTML:           "<h1>New stuff</h1>",
				Subject:        "Cool New Stuff",
				TemplateID:     "random-template-id",
				ReplyTo:        "reply-to-address",
				SenderID:       "my-sender",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(campaign.ID).To(Equal("deadbeef-aabb-ccdd-eeff-001122334455"))

			retrievedCampaign, err := repo.Get(connection, campaign.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedCampaign).To(Equal(campaign))
		})
	})
})
