package senders_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/senders"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateHandler", func() {
	var (
		handler           senders.CreateHandler
		sendersCollection *mocks.SendersCollection
		context           stack.Context
		writer            *httptest.ResponseRecorder
		request           *http.Request
		database          *mocks.Database
	)

	BeforeEach(func() {
		database = mocks.NewDatabase()
		context = stack.NewContext()
		context.Set("client_id", "some-client-id")
		context.Set("database", database)

		sendersCollection = mocks.NewSendersCollection()
		sendersCollection.SetCall.Returns.Sender = collections.Sender{
			ID:   "some-sender-id",
			Name: "some-sender",
		}

		writer = httptest.NewRecorder()

		requestBody, err := json.Marshal(map[string]string{
			"name": "some-sender",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err = http.NewRequest("POST", "/senders", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler = senders.NewCreateHandler(sendersCollection)
	})

	It("creates a sender", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(sendersCollection.SetCall.Receives.Conn).To(Equal(database.Conn))
		Expect(database.ConnectionWasCalled).To(BeTrue())
		Expect(sendersCollection.SetCall.Receives.Sender).To(Equal(collections.Sender{
			Name:     "some-sender",
			ClientID: "some-client-id",
		}))

		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-sender-id",
			"name": "some-sender"
		}`))
	})

	Context("failure cases", func() {
		It("returns a 400 when the JSON cannot be unmarshalled", func() {
			var err error
			request, err = http.NewRequest("POST", "/senders", strings.NewReader("%%%"))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusBadRequest))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": [
					"invalid json body"
				]
			}`))
		})

		It("returns a 422 when the request does not include a sender name", func() {
			var err error
			request, err = http.NewRequest("POST", "/senders", strings.NewReader("{}"))
			Expect(err).NotTo(HaveOccurred())

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(422))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": [
					"missing sender name"
				]
			}`))
		})

		It("returns a 401 when the request does not include a client id", func() {
			context.Set("client_id", "")

			handler.ServeHTTP(writer, request, context)
			Expect(writer.Code).To(Equal(http.StatusUnauthorized))
			Expect(writer.Body.String()).To(MatchJSON(`{
				"errors": [
					"missing client id"
				]
			}`))
		})

		It("returns a 500 when the collection indicates a system error", func() {
			sendersCollection.SetCall.Returns.Err = errors.New("BOOM!")

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
