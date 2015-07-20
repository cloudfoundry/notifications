package messages_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/v1/messages"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetHandler", func() {
	var (
		handler       messages.GetHandler
		errorWriter   *fakes.ErrorWriter
		writer        *httptest.ResponseRecorder
		request       *http.Request
		messageID     string
		err           error
		messageFinder *fakes.MessageFinder
		database      *fakes.Database
		context       stack.Context
	)

	BeforeEach(func() {
		errorWriter = fakes.NewErrorWriter()
		messageFinder = fakes.NewMessageFinder()
		writer = httptest.NewRecorder()
		messageID = "message-123"
		database = fakes.NewDatabase()
		context = stack.NewContext()
		context.Set("database", database)

		request, err = http.NewRequest("GET", "/messages/"+messageID, nil)
		if err != nil {
			panic(err)
		}

		handler = messages.NewGetHandler(messageFinder, errorWriter)
	})

	Describe("ServeHTTP", func() {
		It("Returns the status of the given message from the finder", func() {
			messageFinder.Messages[messageID] = services.Message{
				Status: "The generic status returned",
			}

			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.Bytes()).To(MatchJSON(`{
				"status": "The generic status returned"
			}`))

			Expect(messageFinder.FindCall.Arguments).To(ConsistOf([]interface{}{database, messageID}))
		})

		Context("When the finder errors", func() {
			It("Delegates to the error writer", func() {
				findError := errors.New("The finder returns a generic error")
				messageFinder.FindCall.Error = findError

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(findError))
			})
		})
	})
})
