package senders_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/senders"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetHandler", func() {
	var (
		handler           senders.GetHandler
		sendersCollection *fakes.SendersCollection
		context           stack.Context
		writer            *httptest.ResponseRecorder
		request           *http.Request
		database          *fakes.Database
	)

	BeforeEach(func() {
		database = fakes.NewDatabase()
		context = stack.NewContext()
		context.Set("client_id", "some-client-id")
		context.Set("database", database)

		sendersCollection = fakes.NewSendersCollection()
		sendersCollection.GetCall.ReturnSender = collections.Sender{
			ID:   "some-sender-id",
			Name: "some-sender",
		}

		writer = httptest.NewRecorder()

		var err error
		request, err = http.NewRequest("GET", "/senders/some-sender-id", nil)
		Expect(err).NotTo(HaveOccurred())

		handler = senders.NewGetHandler(sendersCollection)
	})

	It("gets a sender", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(sendersCollection.GetCall.Conn).To(Equal(database.Conn))
		Expect(database.ConnectionWasCalled).To(BeTrue())

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-sender-id",
			"name": "some-sender"
		}`))

		Expect(sendersCollection.GetCall.SenderID).To(Equal("some-sender-id"))
		Expect(sendersCollection.GetCall.ClientID).To(Equal("some-client-id"))
	})

	Context("failure cases", func() {
		It("returns a 422 when the URL does not include a sender_id", func() {
			sendersCollection.GetCall.Err = collections.ValidationError{
				Err: errors.New("missing sender id"),
			}

			var err error
			request, err = http.NewRequest("GET", "/senders/", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "missing sender id"
			}`))
		})

		It("returns a 404 when the sender does not exist", func() {
			sendersCollection.GetCall.Err = collections.NotFoundError{
				Err: errors.New("sender not found"),
			}

			var err error
			request, err = http.NewRequest("GET", "/senders/non-existent-sender-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "sender not found"
			}`))
		})

		It("returns a 500 when the collection indicates a system error", func() {
			sendersCollection.GetCall.Err = errors.New("BOOM!")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"error": "BOOM!"
			}`))
		})
	})
})
