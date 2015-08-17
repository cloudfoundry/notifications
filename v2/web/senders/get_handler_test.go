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
		sendersCollection.GetCall.Returns.Sender = collections.Sender{
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

		Expect(sendersCollection.GetCall.Receives.Conn).To(Equal(database.Conn))
		Expect(database.ConnectionWasCalled).To(BeTrue())

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-sender-id",
			"name": "some-sender"
		}`))

		Expect(sendersCollection.GetCall.Receives.SenderID).To(Equal("some-sender-id"))
		Expect(sendersCollection.GetCall.Receives.ClientID).To(Equal("some-client-id"))
	})

	Context("failure cases", func() {
		It("returns a 301 when the sender id is missing", func() {
			var err error
			request, err = http.NewRequest("GET", "/senders/", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusMovedPermanently))
			headers := writer.Header()
			Expect(headers.Get("Location")).To(Equal("/senders"))
			Expect(writer.Body.String()).To(BeEmpty())
		})

		It("returns a 401 when the client id is missing", func() {
			context.Set("client_id", "")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusUnauthorized))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": [
					"missing client id"
				]
			}`))
		})

		It("returns a 404 when the sender does not exist", func() {
			sendersCollection.GetCall.Returns.Err = collections.NotFoundError{
				Err: errors.New("sender with id \"non-existent-sender-id\" not found"),
			}

			var err error
			request, err = http.NewRequest("GET", "/senders/non-existent-sender-id", nil)
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": [
					"sender with id \"non-existent-sender-id\" not found"
				]
			}`))
		})

		It("returns a 500 when the collection indicates a system error", func() {
			sendersCollection.GetCall.Returns.Err = errors.New("BOOM!")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": [
					"BOOM!"
				]
			}`))
		})
	})
})
