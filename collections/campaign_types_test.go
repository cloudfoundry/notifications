package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/helpers"
	"github.com/cloudfoundry-incubator/notifications/models"

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

	Describe("Add", func() {
		var (
			campaignType collections.CampaignType
		)

		BeforeEach(func() {
			campaignType = collections.CampaignType{
				Name:        helpers.AddressOfString("My cool campaign type"),
				Description: helpers.AddressOfString("description"),
				Critical:    helpers.AddressOfBool(false),
				TemplateID:  helpers.AddressOfString(""),
				SenderID:    "mysender",
			}

			fakeCampaignTypesRepository.InsertCall.ReturnCampaignType.ID = "generated-id"
		})

		It("adds a campaign type to the collection", func() {
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "mysender",
				Name:     "some-sender",
				ClientID: "client_id",
			}

			returnedCampaignType, err := campaignTypesCollection.Add(fakeDatabaseConnection, campaignType, "client_id")
			Expect(err).NotTo(HaveOccurred())
			Expect(returnedCampaignType.ID).To(Equal("generated-id"))
			Expect(fakeCampaignTypesRepository.InsertCall.Connection).To(Equal(fakeDatabaseConnection))
			Expect(fakeCampaignTypesRepository.InsertCall.CampaignType.Name).To(Equal(*campaignType.Name))
			Expect(fakeCampaignTypesRepository.InsertCall.CampaignType.Description).To(Equal(*campaignType.Description))
			Expect(fakeCampaignTypesRepository.InsertCall.CampaignType.Critical).To(Equal(*campaignType.Critical))
			Expect(fakeCampaignTypesRepository.InsertCall.CampaignType.TemplateID).To(Equal(*campaignType.TemplateID))
			Expect(fakeCampaignTypesRepository.InsertCall.CampaignType.SenderID).To(Equal(campaignType.SenderID))
		})

		It("requires a name to be specified", func() {
			campaignType = collections.CampaignType{
				Description: helpers.AddressOfString("description"),
				Critical:    helpers.AddressOfBool(false),
				TemplateID:  helpers.AddressOfString(""),
				SenderID:    "mysender",
			}
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "mysender",
				Name:     "some-sender",
				ClientID: "client_id",
			}

			_, err := campaignTypesCollection.Add(fakeDatabaseConnection, campaignType, "client_id")
			Expect(err).To(MatchError(collections.NewValidationError("missing campaign type name")))
		})

		It("requires a description to be specified", func() {
			campaignType = collections.CampaignType{
				Name:       helpers.AddressOfString("some-campaign-type"),
				Critical:   helpers.AddressOfBool(false),
				TemplateID: helpers.AddressOfString(""),
				SenderID:   "mysender",
			}
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "mysender",
				Name:     "some-sender",
				ClientID: "client_id",
			}

			_, err := campaignTypesCollection.Add(fakeDatabaseConnection, campaignType, "client_id")
			Expect(err).To(MatchError(collections.NewValidationError("missing campaign type description")))
		})

		Context("failure cases", func() {
			It("generates a not found error when the sender does not exist", func() {
				fakeSendersRepository.GetCall.Err = models.RecordNotFoundError("sender with sender ID ROBOTS not found")

				_, err := campaignTypesCollection.Add(fakeDatabaseConnection, campaignType, "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Message: "Sender mysender not found",
					Err:     models.RecordNotFoundError("sender with sender ID ROBOTS not found"),
				}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "mismatch-client-id",
				}

				_, err := campaignTypesCollection.Add(fakeDatabaseConnection, campaignType, "some-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("Sender mysender not found")))
			})
		})
	})

	Describe("List", func() {
		It("retrieves a list of campaign types from the collection", func() {
			fakeCampaignTypesRepository.ListCall.ReturnCampaignTypeList = []models.CampaignType{
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
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
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

		It("retrieves an empty list of campaign types from the collection if no records have been added", func() {
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			returnedCampaignTypeList, err := campaignTypesCollection.List(fakeDatabaseConnection, "some-senderid", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(returnedCampaignTypeList)).To(Equal(0))
		})

		Context("failure cases", func() {
			It("validates that a sender id was specified", func() {
				_, err := campaignTypesCollection.List(fakeDatabaseConnection, "", "some-client-id")
				Expect(err).To(MatchError(collections.NewValidationError("missing sender id")))
			})

			It("validates that a client id was specified", func() {
				_, err := campaignTypesCollection.List(fakeDatabaseConnection, "some-sender-id", "")
				Expect(err).To(MatchError(collections.NewValidationError("missing client id")))
			})

			It("generates a not found error when the sender does not exist", func() {
				fakeSendersRepository.GetCall.Err = models.RecordNotFoundError("sender not found")

				_, err := campaignTypesCollection.List(fakeDatabaseConnection, "missing-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{
					Message: "sender not found",
					Err:     models.RecordNotFoundError("sender not found"),
				}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				fakeCampaignTypesRepository.ListCall.Err = models.RecordNotFoundError("sender not found")
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "mismatch-client-id",
				}

				_, err := campaignTypesCollection.List(fakeDatabaseConnection, "mismatch-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("sender not found")))
			})

			It("handles unexpected database errors", func() {
				fakeCampaignTypesRepository.ListCall.ReturnCampaignTypeList = []models.CampaignType{}
				fakeCampaignTypesRepository.ListCall.Err = errors.New("BOOM!")
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
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

	Describe("Get", func() {
		It("returns the ID if it is found", func() {
			fakeCampaignTypesRepository.GetReturn.CampaignType = models.CampaignType{
				ID:       "a-campaign-type-id",
				Name:     "typename",
				SenderID: "senderID",
			}
			fakeSendersRepository.GetCall.ReturnSender = models.Sender{
				ID:       "senderID",
				Name:     "I dont matter",
				ClientID: "some-client-id",
			}

			campaignType, err := campaignTypesCollection.Get(fakeDatabaseConnection, "a-campaign-type-id", "senderID", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(campaignType.Name).To(Equal(helpers.AddressOfString("typename")))
		})

		Context("failure cases", func() {
			It("validates that a campaign type id was specified", func() {
				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "", "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NewValidationError("missing campaign type id")))
			})

			It("validates that a sender id was specified", func() {
				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "", "some-client-id")
				Expect(err).To(MatchError(collections.NewValidationError("missing sender id")))
			})

			It("validates that a client id was specified", func() {
				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "some-sender-id", "")
				Expect(err).To(MatchError(collections.NewValidationError("missing client id")))
			})

			It("returns a not found error if the campaign type does not exist", func() {
				fakeCampaignTypesRepository.GetReturn.Err = models.RecordNotFoundError("campaign type not found")
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
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
				fakeCampaignTypesRepository.GetReturn.CampaignType = models.CampaignType{
					ID:       "some-campaign-type-id",
					Name:     "typename",
					SenderID: "some-sender-id",
				}
				fakeSendersRepository.GetCall.Err = models.RecordNotFoundError("sender not found")
				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "missing-sender-id", "some-client-id")
				Expect(err.(collections.NotFoundError).Message).To(Equal("sender some-campaign-type-id not found"))
			})

			It("returns a not found error if the campaign type does not belong to the sender", func() {
				fakeCampaignTypesRepository.GetReturn.CampaignType = models.CampaignType{
					ID:       "some-campaign-type-id",
					Name:     "typename",
					SenderID: "my-sender-id",
				}
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "someone-elses-sender-id",
					Name:     "some-sender",
					ClientID: "some-client-id",
				}
				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "someone-elses-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("campaign type some-campaign-type-id not found")))
			})

			It("returns a not found error if the sender does not belong to the client", func() {
				fakeCampaignTypesRepository.GetReturn.CampaignType = models.CampaignType{
					ID:       "some-campaign-type-id",
					Name:     "typename",
					SenderID: "my-sender-id",
				}
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
					ID:       "my-sender-id",
					Name:     "some-sender",
					ClientID: "client_id",
				}
				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "my-sender-id", "someone-elses-client-id")
				Expect(err).To(MatchError(collections.NewNotFoundError("sender my-sender-id not found")))
			})

			It("handles unexpected database errors from the senders repository", func() {
				fakeSendersRepository.GetCall.Err = errors.New("BOOM!")

				_, err := campaignTypesCollection.Get(fakeDatabaseConnection, "some-campaign-type-id", "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})

			It("handles unexpected database errors from the campaign types repository", func() {
				fakeCampaignTypesRepository.GetReturn.Err = errors.New("BOOM!")
				fakeSendersRepository.GetCall.ReturnSender = models.Sender{
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

	Describe("Update", func() {
		var (
			campaignType collections.CampaignType
		)

		BeforeEach(func() {
			campaignType = collections.CampaignType{
				ID:          "existing-campaign-type-id",
				Name:        helpers.AddressOfString("My cool campaign type"),
				Description: helpers.AddressOfString("description"),
				Critical:    helpers.AddressOfBool(false),
				TemplateID:  helpers.AddressOfString(""),
				SenderID:    "mysender",
			}
		})

		PIt("updates the name if the campaign type is found", func() {
		})

		PIt("updates the description if the campaign type is found", func() {
		})

		PIt("updates the critical flag if the campaign type is found", func() {
		})
	})
})
