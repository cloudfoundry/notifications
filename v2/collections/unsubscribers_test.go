package collections_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UnsubscribersCollection", func() {
	var (
		unsubscribersRepository *mocks.UnsubscribersRepository
		connection              *mocks.Connection
		unsubscribersCollection collections.UnsubscribersCollection
		userFinder              *mocks.UserFinder
		campaignTypesRepository *mocks.CampaignTypesRepository
	)

	BeforeEach(func() {
		unsubscribersRepository = mocks.NewUnsubscribersRepository()
		connection = mocks.NewConnection()
		userFinder = mocks.NewUserFinder()
		campaignTypesRepository = mocks.NewCampaignTypesRepository()
		unsubscribersCollection = collections.NewUnsubscribersCollection(unsubscribersRepository, campaignTypesRepository, userFinder)
	})

	Describe("Set", func() {
		BeforeEach(func() {
			unsubscribersRepository.InsertCall.Returns.Unsubscriber = models.Unsubscriber{
				ID:             "some-id",
				CampaignTypeID: "some-campaign-type-id",
				UserGUID:       "some-user-guid",
			}
			userFinder.ExistsCall.Returns.Exists = true
		})

		Context("when the unsubscriber does not exist", func() {
			It("will insert an unsubscriber into the collection", func() {
				unsubscriber, err := unsubscribersCollection.Set(connection, collections.Unsubscriber{
					CampaignTypeID: "some-campaign-type-id",
					UserGUID:       "some-user-guid",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(userFinder.ExistsCall.Receives.GUID).To(Equal("some-user-guid"))
				Expect(unsubscribersRepository.InsertCall.Receives.Connection).To(Equal(connection))
				Expect(unsubscribersRepository.InsertCall.Receives.Unsubscriber).To(Equal(models.Unsubscriber{
					CampaignTypeID: "some-campaign-type-id",
					UserGUID:       "some-user-guid",
				}))

				Expect(unsubscriber).To(Equal(collections.Unsubscriber{
					ID:             "some-id",
					CampaignTypeID: "some-campaign-type-id",
					UserGUID:       "some-user-guid",
				}))
			})
		})

		Context("when an error occurs", func() {
			Context("when the userFinder errors", func() {
				It("returns the error", func() {
					userFinder.ExistsCall.Returns.Error = errors.New("some error")
					_, err := unsubscribersCollection.Set(connection, collections.Unsubscriber{
						CampaignTypeID: "some-campaign-type-id",
						UserGUID:       "some-weird-user",
					})
					Expect(err).To(MatchError("some error"))
				})
			})

			Describe("when the user does not exist", func() {
				It("returns a record not found error", func() {
					userFinder.ExistsCall.Returns.Exists = false
					_, err := unsubscribersCollection.Set(connection, collections.Unsubscriber{
						CampaignTypeID: "some-campaign-type-id",
						UserGUID:       "some-weird-user",
					})
					Expect(err).To(MatchError(collections.NotFoundError{errors.New(`User "some-weird-user" not found`)}))
				})
			})

			Describe("when the campaignType does not exist", func() {
				It("returns a record not found error", func() {
					campaignTypesRepository.GetCall.Returns.Error = models.RecordNotFoundError{errors.New("some-record-not-found-error")}
					_, err := unsubscribersCollection.Set(connection, collections.Unsubscriber{
						CampaignTypeID: "non-existent-campaign-type-id",
						UserGUID:       "some-user",
					})
					Expect(err).To(MatchError(collections.NotFoundError{models.RecordNotFoundError{errors.New(`some-record-not-found-error`)}}))
				})
			})

			Describe("when the campaignType is critical", func() {
				It("returns a permissions error", func() {
					campaignTypesRepository.GetCall.Returns.CampaignType = models.CampaignType{Critical: true}
					_, err := unsubscribersCollection.Set(connection, collections.Unsubscriber{
						CampaignTypeID: "some-critical-campaign-type",
						UserGUID:       "some-user",
					})
					Expect(err).To(MatchError(collections.PermissionsError{errors.New("Campaign type \"some-critical-campaign-type\" cannot be unsubscribed from")}))
				})
			})

			Describe("when the campaignTypes repository has any other error", func() {
				It("returns the error", func() {
					campaignTypesRepository.GetCall.Returns.Error = errors.New("some-error")
					_, err := unsubscribersCollection.Set(connection, collections.Unsubscriber{
						CampaignTypeID: "non-existent-campaign-type-id",
						UserGUID:       "some-user",
					})
					Expect(err).To(MatchError(errors.New("some-error")))
				})
			})

			Describe("when the unsubscribers repository has an error", func() {
				It("returns a persistence error", func() {
					unsubscribersRepository.InsertCall.Returns.Error = errors.New("some-other-error")
					_, err := unsubscribersCollection.Set(connection, collections.Unsubscriber{
						CampaignTypeID: "some-campaign-type-id",
						UserGUID:       "some-user-guid",
					})

					Expect(err).To(MatchError(ContainSubstring("some-other-error")))
					Expect(err).To(BeAssignableToTypeOf(collections.PersistenceError{}))
				})
			})
		})
	})

	Describe("Delete", func() {
		Context("when the unsubscriber exists", func() {
			It("will delete the unsubscriber from the collection", func() {
				err := unsubscribersCollection.Delete(connection, collections.Unsubscriber{
					CampaignTypeID: "some-campaign-type-id",
					UserGUID:       "some-user-guid",
				})
				Expect(err).NotTo(HaveOccurred())

				Expect(unsubscribersRepository.DeleteCall.Receives.Connection).To(Equal(connection))
				Expect(unsubscribersRepository.DeleteCall.Receives.Unsubscriber).To(Equal(models.Unsubscriber{
					CampaignTypeID: "some-campaign-type-id",
					UserGUID:       "some-user-guid",
				}))
			})
		})
	})
})
