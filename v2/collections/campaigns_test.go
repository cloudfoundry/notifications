package collections_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignsCollection", func() {
	var (
		startTime         time.Time
		conn              *mocks.Connection
		enqueuer          *mocks.CampaignEnqueuer
		collection        collections.CampaignsCollection
		campaignsRepo     *mocks.CampaignsRepository
		campaignTypesRepo *mocks.CampaignTypesRepository
		templatesRepo     *mocks.TemplatesRepository
		sendersRepo       *mocks.SendersRepository
		userFinder        *mocks.UserFinder
		spaceFinder       *mocks.SpaceFinder
		orgFinder         *mocks.OrgFinder
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()
		enqueuer = mocks.NewCampaignEnqueuer()
		campaignsRepo = mocks.NewCampaignsRepository()
		campaignTypesRepo = mocks.NewCampaignTypesRepository()
		templatesRepo = mocks.NewTemplatesRepository()
		sendersRepo = mocks.NewSendersRepository()
		userFinder = mocks.NewUserFinder()
		spaceFinder = mocks.NewSpaceFinder()
		orgFinder = mocks.NewOrgFinder()

		var err error
		startTime, err = time.Parse(time.RFC3339, "2015-09-01T12:34:56-07:00")
		Expect(err).NotTo(HaveOccurred())

		collection = collections.NewCampaignsCollection(enqueuer, campaignsRepo, campaignTypesRepo, templatesRepo, sendersRepo, userFinder, spaceFinder, orgFinder)
	})

	Describe("Create", func() {
		BeforeEach(func() {
			sendersRepo.GetCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}
		})

		Context("when the audience isn't a thing", func() {
			It("returns an error", func() {
				campaign := collections.Campaign{
					SendTo:         map[string][]string{"not a thing": {"some-thing-guid"}},
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					TemplateID:     "whoa-a-template-id",
					ReplyTo:        "nothing@example.com",
					SenderID:       "some-sender-id",
				}

				_, err := collection.Create(conn, campaign, "some-client-id", false)
				Expect(err).To(MatchError(collections.UnknownError{errors.New("The \"not a thing\" audience is not valid")}))
			})
		})

		Context("when the audience is an email", func() {
			Context("enqueuing a campaignJob", func() {
				BeforeEach(func() {
					campaignsRepo.InsertCall.Returns.Campaign = models.Campaign{
						ID:             "a-new-id",
						SendTo:         `{"emails":["test1@example.com","test2@example.com"]}`,
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
					}
				})

				It("returns a campaignID after enqueuing the campaign with its type", func() {
					campaign := collections.Campaign{
						SendTo:         map[string][]string{"emails": {"test1@example.com", "test2@example.com"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
						StartTime:      startTime,
					}

					enqueuedCampaign, err := collection.Create(conn, campaign, "some-client-id", false)
					Expect(err).NotTo(HaveOccurred())

					Expect(campaignsRepo.InsertCall.Receives.Connection).To(Equal(conn))
					Expect(campaignsRepo.InsertCall.Receives.Campaign).To(Equal(models.Campaign{
						SendTo:         `{"emails":["test1@example.com","test2@example.com"]}`,
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
						StartTime:      startTime,
					}))

					Expect(enqueuer.EnqueueCall.Receives.Campaign).To(Equal(collections.Campaign{
						ID:             "a-new-id",
						SendTo:         map[string][]string{"emails": {"test1@example.com", "test2@example.com"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
						ClientID:       "some-client-id",
						StartTime:      startTime,
					}))
					Expect(enqueuer.EnqueueCall.Receives.JobType).To(Equal("campaign"))

					Expect(enqueuedCampaign.ID).To(Equal("a-new-id"))
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

		Context("when the audience is a space", func() {
			BeforeEach(func() {
				spaceFinder.ExistsCall.Returns.Exists = true
				campaignsRepo.InsertCall.Returns.Campaign = models.Campaign{
					ID:             "a-new-id",
					SendTo:         `{"spaces":"some-space-guid"}`,
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					TemplateID:     "whoa-a-template-id",
					ReplyTo:        "nothing@example.com",
				}
			})

			Context("enqueuing a campaignJob", func() {
				It("returns a campaignID after enqueuing the campaign with its type", func() {
					campaign := collections.Campaign{
						SendTo:         map[string][]string{"spaces": {"some-space-guid"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
						StartTime:      startTime,
					}

					enqueuedCampaign, err := collection.Create(conn, campaign, "some-client-id", false)
					Expect(err).NotTo(HaveOccurred())

					Expect(enqueuer.EnqueueCall.Receives.Campaign).To(Equal(collections.Campaign{
						ID:             "a-new-id",
						SendTo:         map[string][]string{"spaces": {"some-space-guid"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
						ClientID:       "some-client-id",
						StartTime:      startTime,
					}))
					Expect(enqueuer.EnqueueCall.Receives.JobType).To(Equal("campaign"))

					Expect(enqueuedCampaign.ID).To(Equal("a-new-id"))
					Expect(err).NotTo(HaveOccurred())

					Expect(spaceFinder.ExistsCall.Receives.GUID).To(Equal("some-space-guid"))
				})
			})

			Context("when finding a space causes an error", func() {
				It("returns an error", func() {
					spaceFinder.ExistsCall.Returns.Error = errors.New("something bad happened")

					campaign := collections.Campaign{
						SendTo:         map[string][]string{"spaces": {"some-space-guid"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
					}

					_, err := collection.Create(conn, campaign, "some-client-id", false)
					Expect(err).To(MatchError(collections.UnknownError{errors.New("something bad happened")}))
				})
			})
		})

		Context("when the audience is an org", func() {
			BeforeEach(func() {
				orgFinder.ExistsCall.Returns.Exists = true
				campaignsRepo.InsertCall.Returns.Campaign = models.Campaign{
					ID:             "a-new-id",
					SendTo:         `{"orgs":"some-org-guid"}`,
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					TemplateID:     "whoa-a-template-id",
					ReplyTo:        "nothing@example.com",
				}
			})

			Context("enqueuing a campaignJob", func() {
				It("returns a campaignID after enqueuing the campaign with its type", func() {
					campaign := collections.Campaign{
						SendTo:         map[string][]string{"orgs": {"some-org-guid"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
						StartTime:      startTime,
					}

					enqueuedCampaign, err := collection.Create(conn, campaign, "some-client-id", false)
					Expect(err).NotTo(HaveOccurred())

					Expect(enqueuer.EnqueueCall.Receives.Campaign).To(Equal(collections.Campaign{
						ID:             "a-new-id",
						SendTo:         map[string][]string{"orgs": {"some-org-guid"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
						ClientID:       "some-client-id",
						StartTime:      startTime,
					}))
					Expect(enqueuer.EnqueueCall.Receives.JobType).To(Equal("campaign"))

					Expect(enqueuedCampaign.ID).To(Equal("a-new-id"))
					Expect(err).NotTo(HaveOccurred())

					Expect(orgFinder.ExistsCall.Receives.GUID).To(Equal("some-org-guid"))
				})
			})

			Context("when finding an org causes an error", func() {
				It("returns an error", func() {
					orgFinder.ExistsCall.Returns.Error = errors.New("something bad happened")

					campaign := collections.Campaign{
						SendTo:         map[string][]string{"orgs": {"some-org-guid"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
					}

					_, err := collection.Create(conn, campaign, "some-client-id", false)
					Expect(err).To(MatchError(collections.UnknownError{errors.New("something bad happened")}))
				})
			})
		})

		Context("when the audience is a user", func() {
			BeforeEach(func() {
				userFinder.ExistsCall.Returns.Exists = true
				campaignsRepo.InsertCall.Returns.Campaign = models.Campaign{
					ID:             "a-new-id",
					SendTo:         `{"users":"some-user-guid"}`,
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					TemplateID:     "whoa-a-template-id",
					ReplyTo:        "nothing@example.com",
				}
			})

			Context("enqueuing a campaignJob", func() {
				It("returns a campaignID after enqueuing the campaign with its type", func() {
					campaign := collections.Campaign{
						SendTo:         map[string][]string{"users": {"some-user-guid"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
						StartTime:      startTime,
					}

					enqueuedCampaign, err := collection.Create(conn, campaign, "some-client-id", false)
					Expect(err).NotTo(HaveOccurred())

					Expect(enqueuer.EnqueueCall.Receives.Campaign).To(Equal(collections.Campaign{
						ID:             "a-new-id",
						SendTo:         map[string][]string{"users": {"some-user-guid"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
						ClientID:       "some-client-id",
						StartTime:      startTime,
					}))
					Expect(enqueuer.EnqueueCall.Receives.JobType).To(Equal("campaign"))

					Expect(enqueuedCampaign.ID).To(Equal("a-new-id"))
					Expect(err).NotTo(HaveOccurred())

					Expect(userFinder.ExistsCall.Receives.GUID).To(Equal("some-user-guid"))

					Expect(sendersRepo.GetCall.Receives.SenderID).To(Equal("some-sender-id"))
					Expect(sendersRepo.GetCall.Receives.Connection).To(Equal(conn))
				})
			})

			It("gets the template off of the campaign type if the templateID is blank", func() {
				campaignTypesRepo.GetCall.Returns.CampaignType = models.CampaignType{
					TemplateID: "campaign-type-template-id",
				}

				campaign := collections.Campaign{
					SendTo:         map[string][]string{"users": {"some-guid"}},
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					TemplateID:     "",
					ReplyTo:        "nothing@example.com",
					SenderID:       "some-sender-id",
					StartTime:      startTime,
				}

				_, err := collection.Create(conn, campaign, "some-client-id", false)
				Expect(err).NotTo(HaveOccurred())

				Expect(campaignTypesRepo.GetCall.Receives.Connection).To(Equal(conn))
				Expect(campaignTypesRepo.GetCall.Receives.CampaignTypeID).To(Equal("some-id"))

				Expect(enqueuer.EnqueueCall.Receives.Campaign).To(Equal(collections.Campaign{
					ID:             "a-new-id",
					SendTo:         map[string][]string{"users": {"some-guid"}},
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					TemplateID:     "campaign-type-template-id",
					ReplyTo:        "nothing@example.com",
					SenderID:       "some-sender-id",
					ClientID:       "some-client-id",
					StartTime:      startTime,
				}))
			})

			It("uses the default template if neither the campaign nor the campaign type has one", func() {
				campaign := collections.Campaign{
					SendTo:         map[string][]string{"users": {"some-guid"}},
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					ReplyTo:        "nothing@example.com",
					SenderID:       "some-sender-id",
					StartTime:      startTime,
				}

				_, err := collection.Create(conn, campaign, "some-client-id", false)
				Expect(err).NotTo(HaveOccurred())

				Expect(enqueuer.EnqueueCall.Receives.Campaign).To(Equal(collections.Campaign{
					ID:             "a-new-id",
					SendTo:         map[string][]string{"users": {"some-guid"}},
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					TemplateID:     "default",
					ReplyTo:        "nothing@example.com",
					SenderID:       "some-sender-id",
					ClientID:       "some-client-id",
					StartTime:      startTime,
				}))
			})

			It("allows requestors with critical_notifications.write scope to send critical notifications", func() {
				campaignTypesRepo.GetCall.Returns.CampaignType = models.CampaignType{
					Critical: true,
				}

				campaign := collections.Campaign{
					SendTo:         map[string][]string{"users": {"some-guid"}},
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					TemplateID:     "some-template-id",
					ReplyTo:        "nothing@example.com",
					SenderID:       "some-sender-id",
					StartTime:      startTime,
				}

				_, err := collection.Create(conn, campaign, "some-client-id", true)
				Expect(err).NotTo(HaveOccurred())

				Expect(campaignTypesRepo.GetCall.Receives.Connection).To(Equal(conn))
				Expect(campaignTypesRepo.GetCall.Receives.CampaignTypeID).To(Equal("some-id"))

				Expect(enqueuer.EnqueueCall.Receives.Campaign).To(Equal(collections.Campaign{
					ID:             "a-new-id",
					SendTo:         map[string][]string{"users": {"some-guid"}},
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					TemplateID:     "some-template-id",
					ReplyTo:        "nothing@example.com",
					SenderID:       "some-sender-id",
					ClientID:       "some-client-id",
					StartTime:      startTime,
				}))
			})

			Context("when the user does not exist", func() {
				It("returns a not found error", func() {
					campaign := collections.Campaign{
						SendTo:         map[string][]string{"users": {"missing-user"}},
						CampaignTypeID: "some-id",
						Text:           "some-test",
						HTML:           "no-html",
						Subject:        "some-subject",
						TemplateID:     "whoa-a-template-id",
						ReplyTo:        "nothing@example.com",
						SenderID:       "some-sender-id",
					}

					userFinder.ExistsCall.Returns.Exists = false

					_, err := collection.Create(conn, campaign, "some-client-id", false)
					Expect(err).To(MatchError(collections.NotFoundError{errors.New("The user \"missing-user\" cannot be found")}))
				})
			})

			Context("when an error happens", func() {
				Context("when enqueue fails", func() {
					It("returns the error to the caller", func() {
						campaign := collections.Campaign{
							SendTo:         map[string][]string{"users": {"some-guid"}},
							CampaignTypeID: "some-id",
							Text:           "some-test",
							HTML:           "no-html",
							Subject:        "some-subject",
							TemplateID:     "whoa-a-template-id",
							ReplyTo:        "nothing@example.com",
							SenderID:       "some-sender-id",
						}
						enqueuer.EnqueueCall.Returns.Err = errors.New("enqueue failed")

						_, err := collection.Create(conn, campaign, "some-client-id", false)

						Expect(err).To(Equal(collections.PersistenceError{Err: errors.New("enqueue failed")}))
					})
				})

				Context("when inserting the campaign record fails", func() {
					It("returns the error", func() {
						campaign := collections.Campaign{
							SendTo:         map[string][]string{"users": {"some-guid"}},
							CampaignTypeID: "some-id",
							Text:           "some-test",
							HTML:           "no-html",
							Subject:        "some-subject",
							TemplateID:     "whoa-a-template-id",
							ReplyTo:        "nothing@example.com",
							SenderID:       "some-sender-id",
						}
						campaignsRepo.InsertCall.Returns.Error = errors.New("insert failed")

						_, err := collection.Create(conn, campaign, "some-client-id", false)

						Expect(err).To(Equal(collections.PersistenceError{Err: errors.New("insert failed")}))
					})
				})

				Context("when checking if the template exists", func() {
					var campaign collections.Campaign
					BeforeEach(func() {
						campaign = collections.Campaign{
							SendTo:         map[string][]string{"users": {"some-guid"}},
							CampaignTypeID: "some-id",
							Text:           "some-test",
							HTML:           "no-html",
							Subject:        "some-subject",
							TemplateID:     "error",
							ReplyTo:        "nothing@example.com",
							SenderID:       "some-sender-id",
						}
					})

					It("returns an error if the templateID is not found", func() {
						templatesRepo.GetCall.Returns.Error = models.RecordNotFoundError{}

						_, err := collection.Create(conn, campaign, "some-client-id", false)
						Expect(err).To(MatchError(collections.NotFoundError{models.RecordNotFoundError{}}))
					})

					It("returns a persistence error if there is some other error", func() {
						dbError := errors.New("the database is shutting off")
						templatesRepo.GetCall.Returns.Error = dbError

						_, err := collection.Create(conn, campaign, "some-client-id", false)
						Expect(err).To(MatchError(collections.PersistenceError{dbError}))
					})
				})

				Context("when checking if the sender exists", func() {
					var campaign collections.Campaign

					BeforeEach(func() {
						campaign = collections.Campaign{
							SendTo:         map[string][]string{"users": {"some-user-guid"}},
							CampaignTypeID: "some-id",
							Text:           "some-test",
							HTML:           "no-html",
							Subject:        "some-subject",
							TemplateID:     "whoa-a-template-id",
							ReplyTo:        "nothing@example.com",
							SenderID:       "missing-sender-id",
						}
					})

					It("returns an error if the senderID is not found", func() {
						sendersRepo.GetCall.Returns.Error = models.RecordNotFoundError{errors.New("sender not found")}

						_, err := collection.Create(conn, campaign, "some-client-id", false)
						Expect(err).To(MatchError(collections.NotFoundError{models.RecordNotFoundError{errors.New("sender not found")}}))
					})

					It("returns an error if the senderID belongs to a different client", func() {
						_, err := collection.Create(conn, campaign, "different-client-id", false)
						Expect(err).To(MatchError(collections.NotFoundError{errors.New("Sender with id \"missing-sender-id\" could not be found")}))
					})

					It("returns a persistence error if there is some other error", func() {
						dbError := errors.New("the database is shutting off")
						sendersRepo.GetCall.Returns.Error = dbError

						_, err := collection.Create(conn, campaign, "some-client-id", false)
						Expect(err).To(MatchError(collections.UnknownError{dbError}))
					})
				})

				Context("when checking if the campaign type exists", func() {
					var campaign collections.Campaign

					BeforeEach(func() {
						campaign = collections.Campaign{
							SendTo:         map[string][]string{"users": {"some-guid"}},
							CampaignTypeID: "some-id",
							Text:           "some-test",
							HTML:           "no-html",
							Subject:        "some-subject",
							TemplateID:     "error",
							ReplyTo:        "nothing@example.com",
							SenderID:       "some-sender-id",
						}
					})

					It("returns an error if the campaignTypeID is not found", func() {
						campaignTypesRepo.GetCall.Returns.Error = models.RecordNotFoundError{}

						_, err := collection.Create(conn, campaign, "some-client-id", false)
						Expect(err).To(MatchError(collections.NotFoundError{models.RecordNotFoundError{}}))
					})

					It("returns a persistence error if there is some other error", func() {
						dbError := errors.New("the database is shutting off")
						campaignTypesRepo.GetCall.Returns.Error = dbError

						_, err := collection.Create(conn, campaign, "some-client-id", false)
						Expect(err).To(MatchError(collections.PersistenceError{dbError}))
					})
				})

				Context("when sending critical notifications is not allowed", func() {
					var campaign collections.Campaign

					BeforeEach(func() {
						campaign = collections.Campaign{
							SendTo:         map[string][]string{"users": {"some-guid"}},
							CampaignTypeID: "some-id",
							Text:           "some-test",
							HTML:           "no-html",
							Subject:        "some-subject",
							TemplateID:     "error",
							ReplyTo:        "nothing@example.com",
							SenderID:       "some-sender-id",
						}

						campaignTypesRepo.GetCall.Returns.CampaignType = models.CampaignType{
							Critical: true,
						}
					})

					It("returns a permissions error", func() {
						_, err := collection.Create(conn, campaign, "some-client-id", false)
						Expect(err).To(MatchError(collections.PermissionsError{errors.New("Scope critical_notifications.write is required")}))
					})
				})
			})
		})
	})

	Context("Checking existence", func() {
		Context("when multiple audience types are provided", func() {
			var campaign collections.Campaign

			BeforeEach(func() {
				campaign = collections.Campaign{
					SendTo: map[string][]string{
						"users":  {"some-user-guid"},
						"spaces": {"some-space"},
						"orgs":   {"some-org"},
					},
					CampaignTypeID: "some-id",
					Text:           "some-test",
					HTML:           "no-html",
					Subject:        "some-subject",
					TemplateID:     "error",
					ReplyTo:        "nothing@example.com",
					SenderID:       "some-sender-id",
				}

				sendersRepo.GetCall.Returns.Sender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "some-client-id",
				}

				userFinder.ExistsCall.Returns.Exists = true
				spaceFinder.ExistsCall.Returns.Exists = true
				orgFinder.ExistsCall.Returns.Exists = true
			})

			It("checks existence on all of them", func() {
				_, err := collection.Create(conn, campaign, "some-client-id", false)
				Expect(err).NotTo(HaveOccurred())

				Expect(userFinder.ExistsCall.Receives.GUID).To(Equal("some-user-guid"))
				Expect(spaceFinder.ExistsCall.Receives.GUID).To(Equal("some-space"))
				Expect(orgFinder.ExistsCall.Receives.GUID).To(Equal("some-org"))
			})
		})
	})

	Describe("Get", func() {
		BeforeEach(func() {
			campaignsRepo.GetCall.Returns.Campaign = models.Campaign{
				ID:             "my-campaign-id",
				SendTo:         `{"users": ["some-guid"]}`,
				CampaignTypeID: "some-id",
				Text:           "some-text",
				HTML:           "no-html",
				Subject:        "some-subject",
				TemplateID:     "error",
				ReplyTo:        "nothing@example.com",
				SenderID:       "some-sender-id",
			}

			sendersRepo.GetCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}
		})

		It("returns the details about the campaign", func() {
			campaign, err := collection.Get(conn, "my-campaign-id", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(campaign.ID).To(Equal("my-campaign-id"))
			Expect(campaign.Text).To(Equal("some-text"))
		})

		Context("failure cases", func() {
			It("returns a not found error when the sender does not exist", func() {
				sendersRepo.GetCall.Returns.Error = models.RecordNotFoundError{errors.New("sender not found")}

				_, err := collection.Get(conn, "my-campaign-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{models.RecordNotFoundError{errors.New("sender not found")}}))
			})

			It("returns an unknown error if the senders repo returns an error", func() {
				sendersRepo.GetCall.Returns.Error = errors.New("i made a bad")

				_, err := collection.Get(conn, "my-campaign-id", "some-client-id")
				Expect(err).To(MatchError(collections.UnknownError{errors.New("i made a bad")}))
			})

			It("returns a not found error when the campaign belongs to a different client", func() {
				sendersRepo.GetCall.Returns.Sender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "other-client-id",
				}

				_, err := collection.Get(conn, "my-campaign-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{errors.New("Campaign with id \"my-campaign-id\" could not be found")}))
			})

			It("returns a not found error when the campaign does not exist", func() {
				campaignsRepo.GetCall.Returns.Error = models.RecordNotFoundError{errors.New("campaign not found")}

				_, err := collection.Get(conn, "missing-campaign-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{models.RecordNotFoundError{errors.New("campaign not found")}}))
			})

			It("returns an unknown error if the campaigns repo returns an error", func() {
				campaignsRepo.GetCall.Returns.Error = errors.New("my bad")

				_, err := collection.Get(conn, "my-campaign-id", "some-client-id")
				Expect(err).To(MatchError(collections.UnknownError{errors.New("my bad")}))
			})
		})
	})
})
