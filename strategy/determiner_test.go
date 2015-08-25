package strategy_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/strategy"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Determiner", func() {
	var (
		determiner    strategy.Determiner
		userStrategy  *mocks.Strategy
		spaceStrategy *mocks.Strategy
		database      *mocks.Database
	)
	BeforeEach(func() {
		userStrategy = mocks.NewStrategy()
		spaceStrategy = mocks.NewStrategy()
		database = mocks.NewDatabase()
		determiner = strategy.NewStrategyDeterminer(userStrategy, spaceStrategy)
	})

	Context("when dispatching to a user", func() {
		It("determines the strategy and calls it", func() {
			err := determiner.Determine(database.Connection(), "some-uaa-host", gobble.NewJob(queue.CampaignJob{
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

			Expect(err).NotTo(HaveOccurred())
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
	})

	Context("when dispatching to a space", func() {
		It("determines the strategy and calls it", func() {
			err := determiner.Determine(database.Connection(), "some-uaa-host", gobble.NewJob(queue.CampaignJob{
				Campaign: collections.Campaign{
					ID:             "some-id",
					SendTo:         map[string]string{"space": "some-space-guid"},
					CampaignTypeID: "some-campaign-type-id",
					Text:           "some-text",
					HTML:           "<h1>my-html</h1>",
					Subject:        "The Best subject",
					TemplateID:     "some-template-id",
					ReplyTo:        "noreply@example.com",
					ClientID:       "some-client-id",
				},
			}))

			Expect(err).NotTo(HaveOccurred())
			Expect(spaceStrategy.DispatchCall.Receives.Dispatch).To(Equal(services.Dispatch{
				GUID:       "some-space-guid",
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
	})

	Context("when an error occurs", func() {
		Context("when the campaign cannot be unmarshalled", func() {
			It("returns the error", func() {
				err := determiner.Determine(database.Connection(), "some-uaa-host", gobble.NewJob("%%"))
				Expect(err).To(MatchError("json: cannot unmarshal string into Go value of type queue.CampaignJob"))
			})
		})

		Context("when dispatch errors", func() {
			It("returns the error", func() {
				spaceStrategy.DispatchCall.Returns.Error = errors.New("some error")
				err := determiner.Determine(database.Connection(), "some-uaa-host", gobble.NewJob(queue.CampaignJob{
					Campaign: collections.Campaign{
						SendTo: map[string]string{"space": "some-space-guid"},
					},
				}))
				Expect(err).To(MatchError(errors.New("some error")))
			})
		})

		Context("when the audience is not found", func() {
			It("returns an error", func() {
				err := determiner.Determine(database.Connection(), "some-uaa-host", gobble.NewJob(queue.CampaignJob{
					Campaign: collections.Campaign{
						SendTo: map[string]string{"some-audience": "wut"},
					},
				}))
				Expect(err).To(MatchError(strategy.NoStrategyError{errors.New("No strategy for the \"some-audience\" audience could be found")}))
			})
		})
	})
})
