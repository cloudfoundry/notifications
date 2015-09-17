package senders_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/senders"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListHandler", func() {
	var (
		handler           senders.ListHandler
		sendersCollection *mocks.SendersCollection
		context           stack.Context
		writer            *httptest.ResponseRecorder
		request           *http.Request
		conn              *mocks.Connection
		database          *mocks.Database
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		context = stack.NewContext()
		context.Set("client_id", "whatever")
		context.Set("database", database)

		sendersCollection = mocks.NewSendersCollection()

		sendersList := []collections.Sender{
			{
				ID:   "sender-id-one",
				Name: "first sender",
			},
			{
				ID:   "sender-id-two",
				Name: "second sender",
			},
		}
		sendersCollection.ListCall.Returns.SenderList = sendersList

		writer = httptest.NewRecorder()

		var err error
		request, err = http.NewRequest("GET", "/senders", nil)
		Expect(err).NotTo(HaveOccurred())

		handler = senders.NewListHandler(sendersCollection)
	})

	It("lists all senders", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(sendersCollection.ListCall.Receives.ClientID).To(Equal("whatever"))
		Expect(sendersCollection.ListCall.Receives.Connection).To(Equal(conn))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"senders": [
				{
					"id": "sender-id-one",
					"name": "first sender",
					"_links": {
						"self": {
							"href": "/senders/sender-id-one"
						},
						"campaign_types": {
							"href": "/senders/sender-id-one/campaign_types"
						},
						"campaigns": {
							"href": "/senders/sender-id-one/campaigns"
						}
					}
				},
				{
					"id": "sender-id-two",
					"name": "second sender",
					"_links": {
						"self": {
							"href": "/senders/sender-id-two"
						},
						"campaign_types": {
							"href": "/senders/sender-id-two/campaign_types"
						},
						"campaigns": {
							"href": "/senders/sender-id-two/campaigns"
						}
					}
				}
			],
			"_links": {
				"self": {
					"href": "/senders"
				}
			}
		}`))
	})

	Context("failure cases", func() {
		It("returns a 500 err when an unexpected error happens", func() {
			sendersCollection.ListCall.Returns.Error = errors.New("persistence error")

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": [
					"persistence error"
				]
			}`))
		})
	})
})
