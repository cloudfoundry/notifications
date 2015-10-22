package v2_test

import (
	"bytes"
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/postal/v2"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/horde"
	"github.com/cloudfoundry-incubator/notifications/v2/queue"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignJobProcessor", func() {
	var (
		processor                   v2.CampaignJobProcessor
		database                    *mocks.Database
		connection                  *mocks.Connection
		enqueuer                    *mocks.V2Enqueuer
		users, orgs, emails, spaces *mocks.Audiences
		buffer                      *bytes.Buffer
		logger                      lager.Logger
	)

	BeforeEach(func() {
		database = mocks.NewDatabase()
		connection = mocks.NewConnection()
		database.ConnectionCall.Returns.Connection = connection

		enqueuer = mocks.NewV2Enqueuer()
		emails = mocks.NewAudiences()
		spaces = mocks.NewAudiences()
		orgs = mocks.NewAudiences()
		users = mocks.NewAudiences()
		processor = v2.NewCampaignJobProcessor(notify.EmailFormatter{},
			notify.HTMLExtractor{}, emails, spaces, orgs, users, enqueuer)
		buffer = bytes.NewBuffer([]byte{})
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))
	})

	Context("when the audience is users", func() {
		It("enqueues a job based on the users audience", func() {
			users.GenerateAudiencesCall.Returns.Audiences = []horde.Audience{
				{
					Users: []horde.User{
						{GUID: "some-user-guid"},
						{GUID: "some-other-user-guid"},
					},
					Endorsement: "some endorsement",
				},
			}

			err := processor.Process(database.Connection(), "some-uaa-host", *gobble.NewJob(queue.CampaignJob{
				Campaign: collections.Campaign{
					ID: "some-id",
					SendTo: map[string][]string{
						"users": {"some-user-guid", "some-other-user-guid"},
					},
					CampaignTypeID: "some-campaign-type-id",
					Text:           "some-text",
					HTML:           "<h1>my-html</h1>",
					Subject:        "The Best subject",
					TemplateID:     "some-template-id",
					ReplyTo:        "noreply@example.com",
					ClientID:       "some-client-id",
				},
			}), logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(users.GenerateAudiencesCall.Receives.Inputs).To(Equal([]string{
				"some-user-guid",
				"some-other-user-guid",
			}))

			Expect(enqueuer.EnqueueCall.Receives.Connection).To(Equal(connection))
			Expect(enqueuer.EnqueueCall.Receives.Users).To(ConsistOf([]queue.User{
				{GUID: "some-user-guid", Endorsement: "some endorsement"},
				{GUID: "some-other-user-guid", Endorsement: "some endorsement"},
			}))
			Expect(enqueuer.EnqueueCall.Receives.Options).To(Equal(queue.Options{
				ReplyTo:           "noreply@example.com",
				Subject:           "The Best subject",
				KindDescription:   "",
				SourceDescription: "",
				Text:              "some-text",
				HTML: queue.HTML{
					BodyContent:    "<h1>my-html</h1>",
					BodyAttributes: "",
					Head:           "",
					Doctype:        "",
				},
				KindID:      "",
				To:          "",
				Role:        "",
				Endorsement: "",
				TemplateID:  "some-template-id",
			}))
			Expect(enqueuer.EnqueueCall.Receives.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(enqueuer.EnqueueCall.Receives.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(enqueuer.EnqueueCall.Receives.Client).To(Equal("some-client-id"))
			Expect(enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("some-uaa-host"))
			Expect(enqueuer.EnqueueCall.Receives.Scope).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Receives.VCAPRequestID).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Receives.RequestReceived).To(Equal(time.Time{}))
			Expect(enqueuer.EnqueueCall.Receives.CampaignID).To(Equal("some-id"))
		})
	})

	Context("when the audience is emails", func() {
		It("enqueues a job based on the emails audience", func() {
			emails.GenerateAudiencesCall.Returns.Audiences = []horde.Audience{
				{
					Users: []horde.User{
						{Email: "some-user@example.com"},
						{Email: "some-other-user@example.com"},
						{Email: "some-user@example.com"},
					},
					Endorsement: "some endorsement",
				},
			}

			err := processor.Process(database.Connection(), "some-uaa-host", *gobble.NewJob(queue.CampaignJob{
				Campaign: collections.Campaign{
					ID: "some-id",
					SendTo: map[string][]string{
						"emails": {
							"some-user@example.com",
							"some-other-user@example.com",
							"some-user@example.com",
						},
					},
					CampaignTypeID: "some-campaign-type-id",
					Text:           "some-text",
					HTML:           "<h1>my-html</h1>",
					Subject:        "The Best subject",
					TemplateID:     "some-template-id",
					ReplyTo:        "noreply@example.com",
					ClientID:       "some-client-id",
				},
			}), logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(emails.GenerateAudiencesCall.Receives.Inputs).To(Equal([]string{
				"some-user@example.com",
				"some-other-user@example.com",
				"some-user@example.com",
			}))

			Expect(enqueuer.EnqueueCall.Receives.Connection).To(Equal(connection))
			Expect(enqueuer.EnqueueCall.Receives.Users).To(ConsistOf([]queue.User{
				{Email: "some-user@example.com", Endorsement: "some endorsement"},
				{Email: "some-other-user@example.com", Endorsement: "some endorsement"},
			}))
			Expect(enqueuer.EnqueueCall.Receives.Options).To(Equal(queue.Options{
				ReplyTo:           "noreply@example.com",
				Subject:           "The Best subject",
				KindDescription:   "",
				SourceDescription: "",
				Text:              "some-text",
				HTML: queue.HTML{
					BodyContent:    "<h1>my-html</h1>",
					BodyAttributes: "",
					Head:           "",
					Doctype:        "",
				},
				KindID:      "",
				To:          "",
				Role:        "",
				Endorsement: "",
				TemplateID:  "some-template-id",
			}))
			Expect(enqueuer.EnqueueCall.Receives.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(enqueuer.EnqueueCall.Receives.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(enqueuer.EnqueueCall.Receives.Client).To(Equal("some-client-id"))
			Expect(enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("some-uaa-host"))
			Expect(enqueuer.EnqueueCall.Receives.Scope).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Receives.VCAPRequestID).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Receives.RequestReceived).To(Equal(time.Time{}))
			Expect(enqueuer.EnqueueCall.Receives.CampaignID).To(Equal("some-id"))
		})
	})

	Context("when the audience is spaces", func() {
		It("enqueues jobs based on the spaces audience", func() {
			spaces.GenerateAudiencesCall.Returns.Audiences = []horde.Audience{
				{
					Users: []horde.User{
						{GUID: "some-user-guid-for-space"},
					},
					Endorsement: "some endorsement",
				},
				{
					Users: []horde.User{
						{GUID: "some-other-user-guid-for-space"},
					},
					Endorsement: "some endorsement",
				},
			}

			err := processor.Process(database.Connection(), "some-uaa-host", *gobble.NewJob(queue.CampaignJob{
				Campaign: collections.Campaign{
					ID: "some-id",
					SendTo: map[string][]string{
						"spaces": {"some-space-guid", "some-other-space-guid"},
					},
					CampaignTypeID: "some-campaign-type-id",
					Text:           "some-text",
					HTML:           "<h1>my-html</h1>",
					Subject:        "The Best subject",
					TemplateID:     "some-template-id",
					ReplyTo:        "noreply@example.com",
					ClientID:       "some-client-id",
				},
			}), logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(spaces.GenerateAudiencesCall.Receives.Inputs).To(Equal([]string{
				"some-space-guid",
				"some-other-space-guid",
			}))
			Expect(spaces.GenerateAudiencesCall.Receives.Logger).To(Equal(logger))

			Expect(enqueuer.EnqueueCall.Receives.Connection).To(Equal(connection))
			Expect(enqueuer.EnqueueCall.Receives.Users).To(ConsistOf([]queue.User{
				{GUID: "some-user-guid-for-space", Endorsement: "some endorsement"},
				{GUID: "some-other-user-guid-for-space", Endorsement: "some endorsement"},
			}))
			Expect(enqueuer.EnqueueCall.Receives.Options).To(Equal(queue.Options{
				ReplyTo:           "noreply@example.com",
				Subject:           "The Best subject",
				KindDescription:   "",
				SourceDescription: "",
				Text:              "some-text",
				HTML: queue.HTML{
					BodyContent:    "<h1>my-html</h1>",
					BodyAttributes: "",
					Head:           "",
					Doctype:        "",
				},
				KindID:      "",
				To:          "",
				Role:        "",
				Endorsement: "",
				TemplateID:  "some-template-id",
			}))
			Expect(enqueuer.EnqueueCall.Receives.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(enqueuer.EnqueueCall.Receives.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(enqueuer.EnqueueCall.Receives.Client).To(Equal("some-client-id"))
			Expect(enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("some-uaa-host"))
			Expect(enqueuer.EnqueueCall.Receives.Scope).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Receives.VCAPRequestID).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Receives.RequestReceived).To(Equal(time.Time{}))
			Expect(enqueuer.EnqueueCall.Receives.CampaignID).To(Equal("some-id"))
		})
	})

	Context("when the audience is organizations", func() {
		It("enqueues jobs based on the organizations audience", func() {
			orgs.GenerateAudiencesCall.Returns.Audiences = []horde.Audience{
				{
					Users: []horde.User{
						{GUID: "some-user-guid-for-org"},
					},
					Endorsement: "some endorsement",
				},
				{
					Users: []horde.User{
						{GUID: "some-other-user-guid-for-org"},
					},
					Endorsement: "some endorsement",
				},
			}

			err := processor.Process(database.Connection(), "some-uaa-host", *gobble.NewJob(queue.CampaignJob{
				Campaign: collections.Campaign{
					ID: "some-id",
					SendTo: map[string][]string{
						"orgs": {"some-org-guid", "some-other-org-guid"},
					},
					CampaignTypeID: "some-campaign-type-id",
					Text:           "some-text",
					HTML:           "<h1>my-html</h1>",
					Subject:        "The Best subject",
					TemplateID:     "some-template-id",
					ReplyTo:        "noreply@example.com",
					ClientID:       "some-client-id",
				},
			}), logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(orgs.GenerateAudiencesCall.Receives.Inputs).To(Equal([]string{
				"some-org-guid",
				"some-other-org-guid",
			}))

			Expect(enqueuer.EnqueueCall.Receives.Connection).To(Equal(connection))
			Expect(enqueuer.EnqueueCall.Receives.Users).To(ConsistOf([]queue.User{
				{GUID: "some-user-guid-for-org", Endorsement: "some endorsement"},
				{GUID: "some-other-user-guid-for-org", Endorsement: "some endorsement"},
			}))
			Expect(enqueuer.EnqueueCall.Receives.Options).To(Equal(queue.Options{
				ReplyTo:           "noreply@example.com",
				Subject:           "The Best subject",
				KindDescription:   "",
				SourceDescription: "",
				Text:              "some-text",
				HTML: queue.HTML{
					BodyContent:    "<h1>my-html</h1>",
					BodyAttributes: "",
					Head:           "",
					Doctype:        "",
				},
				KindID:      "",
				To:          "",
				Role:        "",
				Endorsement: "",
				TemplateID:  "some-template-id",
			}))
			Expect(enqueuer.EnqueueCall.Receives.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(enqueuer.EnqueueCall.Receives.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(enqueuer.EnqueueCall.Receives.Client).To(Equal("some-client-id"))
			Expect(enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("some-uaa-host"))
			Expect(enqueuer.EnqueueCall.Receives.Scope).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Receives.VCAPRequestID).To(Equal(""))
			Expect(enqueuer.EnqueueCall.Receives.RequestReceived).To(Equal(time.Time{}))
			Expect(enqueuer.EnqueueCall.Receives.CampaignID).To(Equal("some-id"))
		})
	})

	Context("when there are multiple audience types", func() {
		BeforeEach(func() {
			orgs.GenerateAudiencesCall.Returns.Audiences = []horde.Audience{
				{
					Users: []horde.User{
						{GUID: "some-user-guid-for-org"},
					},
					Endorsement: "some-org endorsement",
				},
				{
					Users: []horde.User{
						{GUID: "some-other-user-guid-for-org"},
					},
					Endorsement: "some-other-org endorsement",
				},
			}

			spaces.GenerateAudiencesCall.Returns.Audiences = []horde.Audience{
				{
					Users: []horde.User{
						{GUID: "some-user-guid-for-space"},
					},
					Endorsement: "some-space endorsement",
				},
				{
					Users: []horde.User{
						{GUID: "some-other-user-guid-for-space"},
					},
					Endorsement: "some-other-space endorsement",
				},
			}

			users.GenerateAudiencesCall.Returns.Audiences = []horde.Audience{
				{
					Users: []horde.User{
						{GUID: "some-user-guid"},
						{GUID: "some-other-user-guid"},
					},
					Endorsement: "some users endorsement",
				},
			}

			emails.GenerateAudiencesCall.Returns.Audiences = []horde.Audience{
				{
					Users: []horde.User{
						{Email: "some-user@example.com"},
						{Email: "some-other-user@example.com"},
					},
					Endorsement: "some emails endorsement",
				},
			}
		})

		It("enqueues jobs based on the audiences", func() {
			err := processor.Process(database.Connection(), "some-uaa-host", *gobble.NewJob(queue.CampaignJob{
				Campaign: collections.Campaign{
					ID: "some-id",
					SendTo: map[string][]string{
						"orgs":   {"some-org-guid", "some-other-org-guid"},
						"spaces": {"some-space-guid", "some-other-space-guid"},
						"users":  {"some-user-guid", "some-other-user-guid"},
						"emails": {"some-user@example.com", "some-other-user@example.com"},
					},
					CampaignTypeID: "some-campaign-type-id",
					Text:           "some-text",
					HTML:           "<h1>my-html</h1>",
					Subject:        "The Best subject",
					TemplateID:     "some-template-id",
					ReplyTo:        "noreply@example.com",
					ClientID:       "some-client-id",
				},
			}), logger)
			Expect(err).NotTo(HaveOccurred())

			Expect(orgs.GenerateAudiencesCall.Receives.Inputs).To(Equal([]string{
				"some-org-guid",
				"some-other-org-guid",
			}))
			Expect(spaces.GenerateAudiencesCall.Receives.Inputs).To(Equal([]string{
				"some-space-guid",
				"some-other-space-guid",
			}))
			Expect(users.GenerateAudiencesCall.Receives.Inputs).To(Equal([]string{
				"some-user-guid",
				"some-other-user-guid",
			}))
			Expect(emails.GenerateAudiencesCall.Receives.Inputs).To(Equal([]string{
				"some-user@example.com",
				"some-other-user@example.com",
			}))

			Expect(enqueuer.EnqueueCall.Receives.Connection).To(Equal(connection))
			Expect(enqueuer.EnqueueCall.Receives.Users).To(ConsistOf([]queue.User{
				{GUID: "some-user-guid-for-org", Endorsement: "some-org endorsement"},
				{GUID: "some-other-user-guid-for-org", Endorsement: "some-other-org endorsement"},
				{GUID: "some-user-guid-for-space", Endorsement: "some-space endorsement"},
				{GUID: "some-other-user-guid-for-space", Endorsement: "some-other-space endorsement"},
				{GUID: "some-user-guid", Endorsement: "some users endorsement"},
				{GUID: "some-other-user-guid", Endorsement: "some users endorsement"},
				{Email: "some-user@example.com", Endorsement: "some emails endorsement"},
				{Email: "some-other-user@example.com", Endorsement: "some emails endorsement"},
			}))
			Expect(enqueuer.EnqueueCall.Receives.Options).To(Equal(queue.Options{
				ReplyTo: "noreply@example.com",
				Subject: "The Best subject",
				Text:    "some-text",
				HTML: queue.HTML{
					BodyContent: "<h1>my-html</h1>",
				},
				Endorsement: "",
				TemplateID:  "some-template-id",
			}))
			Expect(enqueuer.EnqueueCall.Receives.Client).To(Equal("some-client-id"))
			Expect(enqueuer.EnqueueCall.Receives.UAAHost).To(Equal("some-uaa-host"))
			Expect(enqueuer.EnqueueCall.Receives.CampaignID).To(Equal("some-id"))
		})
	})

	Context("when an error occurs", func() {
		Context("when the campaign cannot be unmarshalled", func() {
			It("returns the error", func() {
				err := processor.Process(database.Connection(), "some-uaa-host", *gobble.NewJob("%%"), logger)
				Expect(err).To(MatchError("json: cannot unmarshal string into Go value of type queue.CampaignJob"))
			})
		})

		Context("when the audience is not found", func() {
			It("returns an error", func() {
				err := processor.Process(database.Connection(), "some-uaa-host", *gobble.NewJob(queue.CampaignJob{
					Campaign: collections.Campaign{
						SendTo: map[string][]string{"some-audience": {"wut"}},
					},
				}), logger)
				Expect(err).To(MatchError(v2.NoAudienceError{errors.New("generator for \"some-audience\" audience could not be found")}))
			})
		})

		Context("when the HTML extractor fails", func() {
			It("returns an error", func() {
				htmlExtractor := mocks.NewHTMLExtractor()
				htmlExtractor.ExtractCall.Returns.Error = errors.New("some extraction error")
				processor = v2.NewCampaignJobProcessor(notify.EmailFormatter{},
					htmlExtractor, emails, spaces, orgs, users, enqueuer)

				err := processor.Process(database.Connection(), "some-uaa-host", *gobble.NewJob(queue.CampaignJob{
					Campaign: collections.Campaign{
						SendTo: map[string][]string{"spaces": {"some-space-guid"}},
					},
				}), logger)
				Expect(err).To(MatchError(errors.New("some extraction error")))
			})
		})

		Context("when the generator fails to generate an audience", func() {
			It("returns an error", func() {
				emails.GenerateAudiencesCall.Returns.Error = errors.New("emails failure")

				err := processor.Process(database.Connection(), "some-uaa-host", *gobble.NewJob(queue.CampaignJob{
					Campaign: collections.Campaign{
						SendTo: map[string][]string{"emails": {"wut@example.com"}},
					},
				}), logger)
				Expect(err).To(MatchError(errors.New("emails failure")))
			})
		})
	})
})
