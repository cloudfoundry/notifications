package strategy_test

import (
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/strategy"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Determiner", func() {
	var (
		determiner   strategy.Determiner
		userStrategy *fakes.Strategy
		database     *fakes.Database
	)
	BeforeEach(func() {
		userStrategy = fakes.NewStrategy()
		database = fakes.NewDatabase()
		determiner = strategy.Determiner{
			UserStrategy: userStrategy,
		}
	})

	It("determines the strategy and calls it", func() {
		determiner.Determine(database.Connection(), "some-uaa-host", gobble.NewJob(queue.CampaignJob{
			Campaign: collections.Campaign{
				ID:             "some-id",
				SendTo:         map[string]string{"user": "some-user-guid"},
				CampaignTypeID: "some-campaign-type-id",
				Text:           "some-text",
				HTML:           "<h1>my-html</h1>",
				Subject:        "The Best subject",
				TemplateID:     "some-template-id",
				ReplyTo:        "noreply@example.com",
				ClientID:       "some-client-id",
			},
		}))

		Expect(userStrategy.DispatchCall.Receives.Dispatch).To(Equal(services.Dispatch{
			GUID:       "some-user-guid",
			UAAHost:    "some-uaa-host",
			Connection: database.Connection(),
			TemplateID: "some-template-id",
			Client: services.DispatchClient{
				ID: "some-client-id",
			},
			Message: services.DispatchMessage{
				To:      "",
				ReplyTo: "noreply@example.com",
				Subject: "The Best subject",
				Text:    "some-text",
				HTML: services.HTML{
					BodyContent: "<h1>my-html</h1>",
				},
			},
		}))
	})

	Context("when an error occurs", func() {
		Context("when the campaign cannot be unmarshalled", func() {
			PIt("returns the error", func() {
			})
		})
		Context("when dispatch errors", func() {
			PIt("returns the error", func() {
			})
		})
	})
})
