package senders_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/senders"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListHandler", func() {
	var (
		handler           senders.ListHandler
		sendersCollection *fakes.SendersCollection
		context           stack.Context
		writer            *httptest.ResponseRecorder
		request           *http.Request
		database          *fakes.Database
	)

	BeforeEach(func() {
		context = stack.NewContext()
		context.Set("client_id", "whatever")

		database = fakes.NewDatabase()
		context.Set("database", database)

		sendersCollection = fakes.NewSendersCollection()

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
		Expect(sendersCollection.ListCall.Receives.Conn).To(Equal(database.Connection()))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"senders": [
				{
					"id": "sender-id-one",
					"name": "first sender"
				},
				{
					"id": "sender-id-two",
					"name": "second sender"
				}
			]
		}`))
	})

	Context("failure cases", func() {
		It("returns a 500 err when an unexpected error happens", func() {
			sendersCollection.ListCall.Returns.Err = errors.New("persistence error")

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
