package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignTypesCollection", func() {
	var (
		campaignTypesCollection     collections.CampaignTypesCollection
		fakeCampaignTypesRepository *fakes.CampaignTypesRepository
		fakeSendersRepository       *fakes.SendersRepository
		fakeDatabaseConnection      *fakes.Connection
	)

	BeforeEach(func() {
		fakeCampaignTypesRepository = fakes.NewCampaignTypesRepository()
		fakeSendersRepository = fakes.NewSendersRepository()

		campaignTypesCollection = collections.NewCampaignTypesCollection(fakeCampaignTypesRepository, fakeSendersRepository)
		fakeDatabaseConnection = fakes.NewConnection()
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			fakeCampaignTypesRepository.GetCall.Returns.CampaignType = models.CampaignType{
				ID:          "some-campaign-id",
				Name:        "My cool campaign type",
				Description: "description",
				Critical:    false,
				TemplateID:  "",
				SenderID:    "mysender",
			}

			fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
				ID:       "mysender",
				Name:     "some-name",
				ClientID: "some-client-random-id",
			}
		})

		Context("when the clientID and senderID are valid", func() {
			It("deletes the given campaign", func() {
				Expect(campaignTypesCollection.Delete(fakeDatabaseConnection, "some-campaign-id", "mysender", "some-client-random-id")).To(Succeed())

				Expect(fakeSendersRepository.GetCall.Receives.Conn).To(Equal(fakeDatabaseConnection))
				Expect(fakeSendersRepository.GetCall.Receives.SenderID).To(Equal("mysender"))

				Expect(fakeCampaignTypesRepository.GetCall.Receives.Connection).To(Equal(fakeDatabaseConnection))
				Expect(fakeCampaignTypesRepository.GetCall.Receives.CampaignTypeID).To(Equal("some-campaign-id"))

				Expect(fakeCampaignTypesRepository.DeleteCall.Receives.CampaignType.ID).To(Equal("some-campaign-id"))
				Expect(fakeCampaignTypesRepository.DeleteCall.Receives.Connection).To(Equal(fakeDatabaseConnection))
			})
		})

		Context("when an error occurs", func() {
			Context("when the sender does not match the client ID", func() {
				It("returns an error", func() {
					err := campaignTypesCollection.Delete(fakeDatabaseConnection, "some-campaign-id", "mysender", "some-other-client-id")
					Expect(err).To(MatchError(collections.NewNotFoundError("sender mysender not found")))
				})
			})

			Context("when the campaign type does not belong to the sender", func() {
				It("returns an error", func() {
					fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
						ID:       "othersender",
						Name:     "some-name",
						ClientID: "some-other-client-id",
					}

					err := campaignTypesCollection.Delete(fakeDatabaseConnection, "some-campaign-id", "othersender", "some-other-client-id")
					Expect(err).To(MatchError(collections.NewNotFoundError("campaign type some-campaign-id not found")))
				})
			})

			Context("when the campaign type does not exist", func() {
				It("returns an error", func() {
					originalError := models.RecordNotFoundError("record not found")
					fakeCampaignTypesRepository.GetCall.Returns.CampaignType = models.CampaignType{}
					fakeCampaignTypesRepository.GetCall.Returns.Err = originalError

					err := campaignTypesCollection.Delete(fakeDatabaseConnection, "some-bad-campaign-id", "mysender", "some-client-random-id")
					Expect(err).To(MatchError(collections.NewNotFoundErrorWithOriginalError("campaign type some-bad-campaign-id not found", originalError)))
				})
			})

			Context("when the sender returns an error", func() {
				It("returns an error", func() {
					originalError := models.RecordNotFoundError("record not found")
					fakeSendersRepository.GetCall.Returns.Sender = models.Sender{}
					fakeSendersRepository.GetCall.Returns.Err = originalError
					err := campaignTypesCollection.Delete(fakeDatabaseConnection, "some-campaign-id", "not-found", "")
					Expect(err).To(MatchError(collections.NewNotFoundErrorWithOriginalError("sender not-found not found", originalError)))
				})
			})

			Context("when the database connection returns some other error while deleting", func() {
				It("returns the error", func() {
					fakeCampaignTypesRepository.DeleteCall.Returns.Err = errors.New("indeletable")

					err := campaignTypesCollection.Delete(fakeDatabaseConnection, "some-campaign-id", "mysender", "some-client-random-id")
					Expect(err).To(MatchError("indeletable"))
				})
			})

			Context("when the database connection returns some other error while getting sender", func() {
				It("returns the error", func() {
					fakeSendersRepository.GetCall.Returns.Sender = models.Sender{}
					fakeSendersRepository.GetCall.Returns.Err = errors.New("nope")

					err := campaignTypesCollection.Delete(fakeDatabaseConnection, "some-campaign-id", "not-found", "")
					Expect(err).To(MatchError(collections.PersistenceError{Err: errors.New("nope")}))
				})
			})

			Context("when the database connection returns some other error while getting campaign type", func() {
				It("returns the error", func() {
					fakeCampaignTypesRepository.GetCall.Returns.Err = errors.New("undeletable")

					err := campaignTypesCollection.Delete(fakeDatabaseConnection, "some-campaign-id", "mysender", "some-client-random-id")
					Expect(err).To(MatchError(collections.PersistenceError{Err: errors.New("undeletable")}))
				})
			})
		})
	})

	Describe("Set", func() {
		var (
			campaignType collections.CampaignType
		)

		BeforeEach(func() {
			campaignType = collections.CampaignType{
				Name:        "My cool campaign type",
				Description: "description",
				Critical:    false,
				TemplateID:  "",
				SenderID:    "mysender",
			}

			fakeCampaignTypesRepository.InsertCall.Returns.CampaignType.ID = "generated-id"
		})

		It("sets a new campaign type within the collection", func() {
			fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
				ID:       "mysender",
				Name:     "some-sender",
				ClientID: "client-id",
			}

			returnedCampaignType, err := campaignTypesCollection.Set(fakeDatabaseConnection, campaignType, "client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCampaignType.ID).To(Equal("generated-id"))
			Expect(fakeCampaignTypesRepository.InsertCall.Receives.Connection).To(Equal(fakeDatabaseConnection))
			Expect(fakeCampaignTypesRepository.InsertCall.Receives.CampaignType).To(Equal(models.CampaignType{
				ID:          "",
				Name:        "My cool campaign type",
				Description: "description",
				Critical:    false,
				TemplateID:  "",
				SenderID:    "mysender",
			}))
		})

		It("sets an existing campaign type within the collection", func() {
			fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
				ID:       "mysender",
				Name:     "some-sender",
				ClientID: "client-id",
			}
			campaignType.ID = "existing-campaign-type-id"
			fakeCampaignTypesRepository.UpdateCall.Returns.CampaignType.ID = "existing-campaign-type-id"

			returnedCampaignType, err := campaignTypesCollection.Set(fakeDatabaseConnection, campaignType, "client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCampaignType.ID).To(Equal("existing-campaign-type-id"))
			Expect(fakeCampaignTypesRepository.UpdateCall.Receives.Connection).To(Equal(fakeDatabaseConnection))
			Expect(fakeCampaignTypesRepository.UpdateCall.Receives.CampaignType).To(Equal(models.CampaignType{
				ID:          "existing-campaign-type-id",
				Name:        "My cool campaign type",
				Description: "description",
				Critical:    false,
				TemplateID:  "",
				SenderID:    "mysender",
			}))
		})

		Context("failure cases", func() {
			It("generates a not found error when the sender does not exist", func() {
				fakeSendersRepository.GetCall.Returns.Err = models.RecordNotFoundError("sender with sender ID ROBOTS not found")

				_, err := campaignTypesCollection.Set(fakeDatabaseConnection, campaignType, "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Message: "sender mysender not found",
					Err:     models.RecordNotFoundError("sender with sender ID ROBOTS not found"),
				}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "mismatch-client-id",
				}

				_, err := campaignTypesCollection.Set(fakeDatabaseConnection, campaignType, "some-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("sender mysender not found")))
			})
		})
	})

	Describe("Get", func() {
		It("returns the ID if it is found", func() {
			fakeCampaignTypesRepository.GetCall.Returns.CampaignType = models.CampaignType{
				ID:       "a-campaign-type-id",
				Name:     "typename",
				SenderID: "senderID",
			}
			fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
				ID:       "senderID",
				Name:     "I dont matter",
				ClientID: "some-client-id",
			}

			campaignType, err := campaignTypesCollection.Get(fakeDatabaseConnection, "a-campaign-type-id", "senderID", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(campaignType.Name).To(Equal("typename"))
		})

		Context("failure cases", func() {
			It("returns a not found error if the campaign type does not exist", func() {
				fakeCampaignTypesRepository.GetCall.Returns.Err = models.RecordNotFoundError("campaign type not found")
				fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
					ID:       "some-sender-id",
					Name:     "I dont matter",
					ClientID: "some-client-id",
				}
				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "missing-campaign-type-id", "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Message: "campaign type missing-campaign-type-id not found",
					Err:     models.RecordNotFoundError("campaign type not found"),
				}))
			})

			It("returns a not found error if the sender does not exist", func() {
				fakeCampaignTypesRepository.GetCall.Returns.CampaignType = models.CampaignType{
					ID:       "some-campaign-type-id",
					Name:     "typename",
					SenderID: "some-sender-id",
				}
				fakeSendersRepository.GetCall.Returns.Err = models.RecordNotFoundError("sender not found")
				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "missing-sender-id", "some-client-id")
				Expect(err.(collections.NotFoundError).Message).To(Equal("sender missing-sender-id not found"))
			})

			It("returns a not found error if the campaign type does not belong to the sender", func() {
				fakeCampaignTypesRepository.GetCall.Returns.CampaignType = models.CampaignType{
					ID:       "some-campaign-type-id",
					Name:     "typename",
					SenderID: "my-sender-id",
				}
				fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
					ID:       "someone-elses-sender-id",
					Name:     "some-sender",
					ClientID: "some-client-id",
				}
				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "someone-elses-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("campaign type some-campaign-type-id not found")))
			})

			It("returns a not found error if the sender does not belong to the client", func() {
				fakeCampaignTypesRepository.GetCall.Returns.CampaignType = models.CampaignType{
					ID:       "some-campaign-type-id",
					Name:     "typename",
					SenderID: "my-sender-id",
				}
				fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
					ID:       "my-sender-id",
					Name:     "some-sender",
					ClientID: "client_id",
				}
				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "my-sender-id", "someone-elses-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("sender my-sender-id not found")))
			})

			It("handles unexpected database errors from the senders repository", func() {
				fakeSendersRepository.GetCall.Returns.Err = errors.New("BOOM!")

				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})

			It("handles unexpected database errors from the campaign types repository", func() {
				fakeCampaignTypesRepository.GetCall.Returns.Err = errors.New("BOOM!")
				fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "some-client-id",
				}

				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})
		})
	})

	Describe("List", func() {
		It("retrieves a list of campaign types from the collection", func() {
			fakeCampaignTypesRepository.ListCall.Returns.CampaignTypeList = []models.CampaignType{
				{
					ID:          "campaign-type-id-one",
					Name:        "first-campaign-type",
					Description: "first-campaign-type-description",
					Critical:    false,
					TemplateID:  "",
					SenderID:    "some-sender-id",
				},
				{
					ID:          "campaign-type-id-two",
					Name:        "second-campaign-type",
					Description: "second-campaign-type-description",
					Critical:    true,
					TemplateID:  "",
					SenderID:    "some-sender-id",
				},
			}
			fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			returnedCampaignTypeList, err := campaignTypesCollection.List(fakeDatabaseConnection, "some-sender-id", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(returnedCampaignTypeList)).To(Equal(2))

			Expect(returnedCampaignTypeList[0].ID).To(Equal("campaign-type-id-one"))
			Expect(returnedCampaignTypeList[0].SenderID).To(Equal("some-sender-id"))

			Expect(returnedCampaignTypeList[1].ID).To(Equal("campaign-type-id-two"))
			Expect(returnedCampaignTypeList[1].SenderID).To(Equal("some-sender-id"))
		})

		It("retrieves an empty list of campaign types from the collection if no records have been Seted", func() {
			fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			returnedCampaignTypeList, err := campaignTypesCollection.List(fakeDatabaseConnection, "some-senderid", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(returnedCampaignTypeList)).To(Equal(0))
		})

		Context("failure cases", func() {
			It("generates a not found error when the sender does not exist", func() {
				fakeSendersRepository.GetCall.Returns.Err = models.RecordNotFoundError("sender not found")

				_, err := campaignTypesCollection.List(fakeDatabaseConnection, "missing-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Message: "sender missing-sender-id not found",
					Err:     models.RecordNotFoundError("sender not found"),
				}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "mismatch-client-id",
				}

				_, err := campaignTypesCollection.List(fakeDatabaseConnection, "mismatch-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("sender mismatch-sender-id not found")))
			})

			It("handles unexpected database errors", func() {
				fakeCampaignTypesRepository.ListCall.Returns.CampaignTypeList = []models.CampaignType{}
				fakeCampaignTypesRepository.ListCall.Returns.Err = errors.New("BOOM!")
				fakeSendersRepository.GetCall.Returns.Sender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "some-client-id",
				}

				_, err := campaignTypesCollection.List(fakeDatabaseConnection, "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})
		})
	})
})
