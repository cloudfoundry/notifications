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

var _ = Describe("DeleteHandler", func() {
	var (
		context           stack.Context
		request           *http.Request
		writer            *httptest.ResponseRecorder
		sendersCollection *mocks.SendersCollection
		conn              *mocks.Connection
		handler           senders.DeleteHandler
	)

	BeforeEach(func() {
		conn = mocks.NewConnection()
		database := mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		context = stack.NewContext()
		context.Set("database", database)
		context.Set("client_id", "some-client-id")

		var err error
		request, err = http.NewRequest("DELETE", "/senders/some-sender-id", nil)
		Expect(err).NotTo(HaveOccurred())
		writer = httptest.NewRecorder()

		sendersCollection = mocks.NewSendersCollection()
		handler = senders.NewDeleteHandler(sendersCollection)
	})

	It("deletes the sender", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusNoContent))
		Expect(writer.Body.String()).To(BeEmpty())

		Expect(sendersCollection.DeleteCall.Receives.Connection).To(Equal(conn))
		Expect(sendersCollection.DeleteCall.Receives.SenderID).To(Equal("some-sender-id"))
		Expect(sendersCollection.DeleteCall.Receives.ClientID).To(Equal("some-client-id"))
	})

	Context("failure cases", func() {
		It("returns a 404 with an error message if the sender does not exist", func() {
			sendersCollection.DeleteCall.Returns.Error = collections.NotFoundError{errors.New("not found")}

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusNotFound))
			Expect(writer.Body).To(MatchJSON(`{
				"errors": ["not found"]
			}`))
		})

		It("returns a 500 with an error message if the collection errors", func() {
			sendersCollection.DeleteCall.Returns.Error = errors.New("the database is irregularly shaped")

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusInternalServerError))
			Expect(writer.Body).To(MatchJSON(`{
				"errors": ["the database is irregularly shaped"]
			}`))
		})
	})
})
