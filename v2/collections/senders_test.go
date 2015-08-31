package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SendersCollection", func() {
	var (
		sendersCollection       collections.SendersCollection
		sendersRepository       *mocks.SendersRepository
		campaignTypesRepository *mocks.CampaignTypesRepository
		conn                    *mocks.Connection
	)

	BeforeEach(func() {
		sendersRepository = mocks.NewSendersRepository()
		campaignTypesRepository = mocks.NewCampaignTypesRepository()

		sendersCollection = collections.NewSendersCollection(sendersRepository, campaignTypesRepository)
		conn = mocks.NewConnection()
	})

	Describe("Set", func() {
		It("adds a sender to the collection", func() {
			sendersRepository.InsertCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			sender, err := sendersCollection.Set(conn, collections.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(sender).To(Equal(collections.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))

			Expect(sendersRepository.InsertCall.Receives.Connection).To(Equal(conn))
			Expect(sendersRepository.InsertCall.Receives.Sender).To(Equal(models.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))
		})

		It("will idempotently add duplicates", func() {
			sendersRepository.InsertCall.Returns.Sender = models.Sender{}
			sendersRepository.InsertCall.Returns.Error = models.DuplicateRecordError{}
			sendersRepository.GetByClientIDAndNameCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}

			sender, err := sendersCollection.Set(conn, collections.Sender{
				Name:     "some-sender",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(sender).To(Equal(collections.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))
			Expect(sendersRepository.GetByClientIDAndNameCall.Receives.Connection).To(Equal(conn))
			Expect(sendersRepository.GetByClientIDAndNameCall.Receives.ClientID).To(Equal("some-client-id"))
			Expect(sendersRepository.GetByClientIDAndNameCall.Receives.Name).To(Equal("some-sender"))
		})

		It("errors if the updated sender conflicts with an existing one", func() {
			sendersRepository.UpdateCall.Returns.Error = models.DuplicateRecordError{}
			_, err := sendersCollection.Set(conn, collections.Sender{
				ID:       "some-sender-id",
				Name:     "existing-sender",
				ClientID: "some-client-id",
			})
			Expect(err).To(MatchError(collections.DuplicateRecordError{models.DuplicateRecordError{}}))
		})

		It("updates a sender if an ID is supplied", func() {
			sendersRepository.UpdateCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "changed-sender",
				ClientID: "some-client-id",
			}
			sender, err := sendersCollection.Set(conn, collections.Sender{
				ID:       "some-sender-id",
				Name:     "changed-sender",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(sender).To(Equal(collections.Sender{
				ID:       "some-sender-id",
				Name:     "changed-sender",
				ClientID: "some-client-id",
			}))

			Expect(sendersRepository.UpdateCall.Receives.Connection).To(Equal(conn))
			Expect(sendersRepository.UpdateCall.Receives.Sender).To(Equal(models.Sender{
				ID:       "some-sender-id",
				Name:     "changed-sender",
				ClientID: "some-client-id",
			}))
		})

		Context("failure cases", func() {
			Context("when inserting", func() {
				It("handles unexpected database errors", func() {
					sendersRepository.InsertCall.Returns.Sender = models.Sender{}
					sendersRepository.InsertCall.Returns.Error = errors.New("BOOM!")

					_, err := sendersCollection.Set(conn, collections.Sender{
						Name:     "some-sender",
						ClientID: "some-client-id",
					})
					Expect(err).To(MatchError(collections.PersistenceError{
						Err: errors.New("BOOM!"),
					}))
				})

				It("returns a persistence error when the sender cannot be found by client id and name", func() {
					sendersRepository.InsertCall.Returns.Sender = models.Sender{}
					sendersRepository.InsertCall.Returns.Error = models.DuplicateRecordError{}

					sendersRepository.GetByClientIDAndNameCall.Returns.Sender = models.Sender{}
					sendersRepository.GetByClientIDAndNameCall.Returns.Error = errors.New("BOOM!")

					_, err := sendersCollection.Set(conn, collections.Sender{
						Name:     "some-sender",
						ClientID: "some-client-id",
					})
					Expect(err).To(MatchError(collections.PersistenceError{
						Err: errors.New("BOOM!"),
					}))
				})
			})

			Context("when updating", func() {
				It("handles unexpected database errors", func() {
					sendersRepository.UpdateCall.Returns.Sender = models.Sender{}
					sendersRepository.UpdateCall.Returns.Error = errors.New("BOOM!")

					_, err := sendersCollection.Set(conn, collections.Sender{
						ID:       "some-sender-id",
						Name:     "some-sender",
						ClientID: "some-client-id",
					})

					Expect(err).To(MatchError(collections.PersistenceError{
						Err: errors.New("BOOM!"),
					}))
				})
			})
		})
	})

	Describe("List", func() {
		It("will list all senders in the collection", func() {
			sendersRepository.ListCall.Returns.Senders = []models.Sender{
				{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "some-client-id",
				},
			}

			senders, err := sendersCollection.List(conn, "some-client-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(len(senders)).To(Equal(1))
			Expect(senders[0].ID).To(Equal("some-sender-id"))

			Expect(sendersRepository.ListCall.Receives.Connection).To(Equal(conn))
			Expect(sendersRepository.ListCall.Receives.ClientID).To(Equal("some-client-id"))
		})

		Context("failure cases", func() {
			It("handles unexpected database errors", func() {
				sendersRepository.ListCall.Returns.Error = errors.New("BOOM!")

				_, err := sendersCollection.List(conn, "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})
		})
	})

	Describe("Get", func() {
		BeforeEach(func() {
			sendersRepository.GetCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}
		})

		It("will retrieve a sender from the collection", func() {
			sender, err := sendersCollection.Get(conn, "some-sender-id", "some-client-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(sender).To(Equal(collections.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))

			Expect(sendersRepository.GetCall.Receives.Connection).To(Equal(conn))
			Expect(sendersRepository.GetCall.Receives.SenderID).To(Equal("some-sender-id"))
		})

		Context("failure cases", func() {
			It("generates a not found error when the sender does not exist", func() {
				recordNotFoundError := models.NewRecordNotFoundError("sender not found")
				sendersRepository.GetCall.Returns.Error = recordNotFoundError

				_, err := sendersCollection.Get(conn, "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{recordNotFoundError}))
			})

			It("generates a not found error when the sender belongs to a different client", func() {
				_, err := sendersCollection.Get(conn, "some-sender-id", "mismatch-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{errors.New(`Sender with id "some-sender-id" could not be found`)}))
			})

			It("handles unexpected database errors", func() {
				sendersRepository.GetCall.Returns.Sender = models.Sender{}
				sendersRepository.GetCall.Returns.Error = errors.New("BOOM!")

				_, err := sendersCollection.Get(conn, "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.PersistenceError{
					Err: errors.New("BOOM!"),
				}))
			})
		})
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			sendersRepository.GetCall.Returns.Sender = models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}
		})

		It("deletes the sender", func() {
			err := sendersCollection.Delete(conn, "some-sender-id", "some-client-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(sendersRepository.GetCall.Receives.Connection).To(Equal(conn))
			Expect(sendersRepository.GetCall.Receives.SenderID).To(Equal("some-sender-id"))

			Expect(sendersRepository.DeleteCall.Receives.Sender).To(Equal(models.Sender{
				ID:       "some-sender-id",
				Name:     "some-sender",
				ClientID: "some-client-id",
			}))
		})

		It("deletes the associated campaign types", func() {
			campaignTypesRepository.ListCall.Returns.CampaignTypeList = []models.CampaignType{
				{
					ID:       "some-campaign-type-id",
					SenderID: "some-sender-id",
				},
			}

			err := sendersCollection.Delete(conn, "some-sender-id", "some-client-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(campaignTypesRepository.ListCall.Receives.Connection).To(Equal(conn))
			Expect(campaignTypesRepository.ListCall.Receives.SenderID).To(Equal("some-sender-id"))

			Expect(campaignTypesRepository.DeleteCall.Receives.Connection).To(Equal(conn))
			Expect(campaignTypesRepository.DeleteCall.Receives.CampaignType).To(Equal(models.CampaignType{
				ID:       "some-campaign-type-id",
				SenderID: "some-sender-id",
			}))
		})

		Context("failure cases", func() {
			It("returns a NotFoundError when the sender does not exist", func() {
				recordNotFoundError := models.RecordNotFoundError{errors.New("not found")}
				sendersRepository.GetCall.Returns.Error = recordNotFoundError

				err := sendersCollection.Delete(conn, "missing-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{recordNotFoundError}))
			})

			It("returns an UnknownError when the repo Get call fails", func() {
				someError := errors.New("unfounded")
				sendersRepository.GetCall.Returns.Error = someError

				err := sendersCollection.Delete(conn, "missing-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.UnknownError{someError}))
			})

			It("returns an UnknownError when the repo Delete call fails", func() {
				someError := errors.New("unfounded")
				sendersRepository.DeleteCall.Returns.Error = someError

				err := sendersCollection.Delete(conn, "missing-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.UnknownError{someError}))
			})

			It("returns a NotFoundError when the sender does not belong to the client", func() {
				sendersRepository.GetCall.Returns.Sender = models.Sender{
					ID:       "some-sender-id",
					Name:     "some-sender",
					ClientID: "other-client-id",
				}

				err := sendersCollection.Delete(conn, "some-sender-id", "some-client-id")
				Expect(err).To(MatchError(collections.NotFoundError{errors.New("Sender with id \"some-sender-id\" could not be found")}))
			})

			It("returns an UnknownError when getting the campaign types returns an error", func() {
				campaignTypesRepository.ListCall.Returns.Error = errors.New("error")

				err := sendersCollection.Delete(conn, "some-sender-id", "some-client-id")

				Expect(err).To(MatchError(collections.UnknownError{errors.New("error")}))
			})

			It("returns an UnknownError when deleting the campaign types returns an error", func() {
				campaignTypesRepository.ListCall.Returns.CampaignTypeList = []models.CampaignType{
					{
						ID:       "some-campaign-type-id",
						SenderID: "some-sender-id",
					},
				}
				campaignTypesRepository.DeleteCall.Returns.Error = errors.New("error")

				err := sendersCollection.Delete(conn, "some-sender-id", "some-client-id")

				Expect(err).To(MatchError(collections.UnknownError{errors.New("error")}))
			})
		})
	})
})
