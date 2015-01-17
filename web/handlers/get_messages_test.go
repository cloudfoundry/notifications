package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetMessages", func() {
	var handler handlers.GetMessages
	var errorWriter *fakes.ErrorWriter
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var messageID string
	var err error
	var messageFinder *fakes.MessageFinder

	BeforeEach(func() {
		errorWriter = fakes.NewErrorWriter()
		messageFinder = fakes.NewMessageFinder()
		handler = handlers.NewGetMessages(messageFinder, errorWriter)
		writer = httptest.NewRecorder()
		messageID = "message-123"

		request, err = http.NewRequest("GET", "/messages/"+messageID, nil)
		if err != nil {
			panic(err)
		}

	})

	Describe("ServeHTTP", func() {
		It("Returns the status of the given message from the finder", func() {
			messageFinder.Messages[messageID] = services.Message{
				Status: "The generic status returned",
			}

			handler.ServeHTTP(writer, request, nil)

			Expect(writer.Code).To(Equal(http.StatusOK))
			Expect(writer.Body.Bytes()).To(MatchJSON(`{
				"status": "The generic status returned"
			}`))
		})

		Context("When the finder errors", func() {
			It("Delegates to the error writer", func() {
				findError := errors.New("The finder returns a generic error")
				messageFinder.FindError = findError
				handler.ServeHTTP(writer, request, nil)
				Expect(errorWriter.Error).To(Equal(findError))
			})
		})
	})
})
