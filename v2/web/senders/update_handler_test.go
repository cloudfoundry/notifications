package senders_test

import (
	"bytes"
	"encoding/json"
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

var _ = Describe("UpdateHandler", func() {
	var (
		handler           senders.UpdateHandler
		sendersCollection *mocks.SendersCollection
		context           stack.Context
		writer            *httptest.ResponseRecorder
		database          *mocks.Database
		conn              *mocks.Connection
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		context = stack.NewContext()
		context.Set("client_id", "some-client-id")
		context.Set("database", database)

		sendersCollection = mocks.NewSendersCollection()
		sendersCollection.SetCall.Returns.Sender = collections.Sender{
			ID:   "some-sender-id",
			Name: "changed-sender",
		}

		writer = httptest.NewRecorder()

		handler = senders.NewUpdateHandler(sendersCollection)
	})

	It("updates a sender", func() {
		requestBody, err := json.Marshal(map[string]string{
			"name": "changed-sender",
		})
		Expect(err).NotTo(HaveOccurred())

		request, err := http.NewRequest("PUT", "/senders/some-sender-id", bytes.NewBuffer(requestBody))
		Expect(err).NotTo(HaveOccurred())

		handler.ServeHTTP(writer, request, context)

		Expect(sendersCollection.SetCall.Receives.Connection).To(Equal(conn))
		Expect(sendersCollection.SetCall.Receives.Sender).To(Equal(collections.Sender{
			ID:       "some-sender-id",
			Name:     "changed-sender",
			ClientID: "some-client-id",
		}))

		Expect(writer.Code).To(Equal(http.StatusOK))
		Expect(writer.Body.String()).To(MatchJSON(`{
			"id": "some-sender-id",
			"name": "changed-sender",
			"_links": {
				"self": {
					"href": "/senders/some-sender-id"
				},
				"campaign_types": {
					"href": "/senders/some-sender-id/campaign_types"
				},
				"campaigns": {
					"href": "/senders/some-sender-id/campaigns"
				}
			}
		}`))
	})

	Context("failure cases", func() {
		Context("when the sender cannot be got", func() {
			It("returns a 404 and a not found error", func() {
				sendersCollection.GetCall.Returns.Error = collections.NotFoundError{errors.New("Sender \"some-missing-sender-id\" does not exist.")}

				requestBody, err := json.Marshal(map[string]string{
					"name": "changed-sender",
				})
				Expect(err).NotTo(HaveOccurred())

				request, err := http.NewRequest("PUT", "/senders/some-missing-sender-id", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusNotFound))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["Sender \"some-missing-sender-id\" does not exist."]}`))

				Expect(sendersCollection.GetCall.Receives.Connection).To(Equal(conn))
				Expect(sendersCollection.GetCall.Receives.SenderID).To(Equal("some-missing-sender-id"))
				Expect(sendersCollection.GetCall.Receives.ClientID).To(Equal("some-client-id"))
			})

			Context("when the json is malformed", func() {
				It("returns an error", func() {
					requestBody := []byte("%%%")
					request, err := http.NewRequest("PUT", "/senders/some-missing-sender-id", bytes.NewBuffer(requestBody))
					Expect(err).NotTo(HaveOccurred())

					handler.ServeHTTP(writer, request, context)

					Expect(writer.Code).To(Equal(http.StatusBadRequest))
					Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["invalid json body"]}`))
				})

			})

			Context("when an unknown error occurs", func() {
				It("returns a 500 and an error", func() {
					sendersCollection.GetCall.Returns.Error = errors.New("some error we don't expect")

					requestBody, err := json.Marshal(map[string]string{
						"name": "changed-sender",
					})
					Expect(err).NotTo(HaveOccurred())

					request, err := http.NewRequest("PUT", "/senders/some-missing-sender-id", bytes.NewBuffer(requestBody))
					Expect(err).NotTo(HaveOccurred())

					handler.ServeHTTP(writer, request, context)

					Expect(writer.Code).To(Equal(http.StatusInternalServerError))
					Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["some error we don't expect"]}`))
				})
			})
		})

		Context("when the sender cannot be set", func() {
			It("returns a 500 with an error message", func() {
				sendersCollection.SetCall.Returns.Error = errors.New("something blew up on set")

				requestBody, err := json.Marshal(map[string]string{
					"name": "changed-sender",
				})
				Expect(err).NotTo(HaveOccurred())

				request, err := http.NewRequest("PUT", "/senders/some-missing-sender-id", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusInternalServerError))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["something blew up on set"]}`))
			})
		})

		Context("when the sender name conflicts", func() {
			It("returns a 422 with an error message", func() {
				sendersCollection.SetCall.Returns.Error = collections.DuplicateRecordError{errors.New("duplicate record")}

				requestBody, err := json.Marshal(map[string]string{
					"name": "changed-sender",
				})
				Expect(err).NotTo(HaveOccurred())

				request, err := http.NewRequest("PUT", "/senders/some-sender-id", bytes.NewBuffer(requestBody))
				Expect(err).NotTo(HaveOccurred())

				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(422))
				Expect(writer.Body.String()).To(MatchJSON(`{"errors": ["duplicate record"]}`))
			})
		})
	})
})
