package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignsCollection", func() {
	Describe("Create", func() {
		Context("when the audience is a user", func() {
			var (
				database          *fakes.Database
				enqueuer          *fakes.CampaignEnqueuer
				collection        collections.CampaignsCollection
				campaignTypesRepo *fakes.CampaignTypesRepository
				sendersRepo       *fakes.SendersRepository
			)

			BeforeEach(func() {
				database = fakes.NewDatabase()
				enqueuer = fakes.NewCampaignEnqueuer()
				campaignTypesRepo = fakes.NewCampaignTypesRepository()
				sendersRepo = fakes.NewSendersRepository()

				collection = collections.NewCampaignsCollection(enqueuer)
			})

			Context("enqueuing a campaignJob", func() {
				It("returns a campaignID after enqueuing the campaign with its type", func() {
					passedCampaign := collections.Campaign{
						SendTo:         map[string]string{"user": "some-guid"},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						ClientID:       "some-client-id",
					}

					campaign, err := collection.Create(database.Connection(), passedCampaign)
					Expect(err).NotTo(HaveOccurred())

					Expect(enqueuer.EnqueueCall.Receives.Campaign).To(Equal(collections.Campaign{
						ID:             "some-random-id",
						SendTo:         map[string]string{"user": "some-guid"},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						ClientID:       "some-client-id",
					}))
					Expect(enqueuer.EnqueueCall.Receives.JobType).To(Equal("campaign"))

					Expect(campaign.ID).To(Equal("some-random-id"))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			PIt("gets the template off of the campaign if the templateID is blank", func() {
			})

			Context("when an error fails", func() {
				Context("when enqueue fails", func() {
					It("returns the error to the caller", func() {
						passedCampaign := collections.Campaign{
							SendTo:         map[string]string{"user": "some-guid"},
							CampaignTypeID: "some-id",
							Text:           "some-test",
							HTML:           "no-html",
							Subject:        "some-subject",
							TemplateID:     "whoa-a-template-id",
							ReplyTo:        "nothing@example.com",
							ClientID:       "another-client-id",
						}
						enqueuer.EnqueueCall.Returns.Err = errors.New("enqueue failed")

						_, err := collection.Create(database.Connection(), passedCampaign)

						Expect(err).To(Equal(collections.PersistenceError{Err: errors.New("enqueue failed")}))
					})
				})
			})
		})
	})
})
