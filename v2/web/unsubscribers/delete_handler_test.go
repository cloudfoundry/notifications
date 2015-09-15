package unsubscribers_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/unsubscribers"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteHandler", func() {
	var (
		handler                 unsubscribers.DeleteHandler
		writer                  *httptest.ResponseRecorder
		request                 *http.Request
		context                 stack.Context
		unsubscribersCollection *mocks.UnsubscribersCollection
		database                *mocks.Database
		connection              *mocks.Connection
	)

	BeforeEach(func() {
		var err error

		unsubscribersCollection = mocks.NewUnsubscribersCollection()
		handler = unsubscribers.NewDeleteHandler(unsubscribersCollection)

		connection = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = connection

		context = stack.NewContext()
		context.Set("database", database)

		writer = httptest.NewRecorder()

		request, err = http.NewRequest("DELETE", "/senders/some-sender-id/campaign_types/some-campaign-type-id/unsubscribers/some-user-guid", nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("removes the unsubscribe for the user", func() {
		handler.ServeHTTP(writer, request, context)
		Expect(writer.Code).To(Equal(http.StatusNoContent))
		Expect(writer.Body.String()).To(BeEmpty())

		Expect(unsubscribersCollection.DeleteCall.Receives.Unsubscriber).To(Equal(collections.Unsubscriber{
			CampaignTypeID: "some-campaign-type-id",
			UserGUID:       "some-user-guid",
		}))
		Expect(unsubscribersCollection.DeleteCall.Receives.Connection).To(Equal(connection))
	})

})
