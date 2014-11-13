package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotifyUser", func() {
	Context("Execute", func() {
		var handler handlers.NotifyUser
		var writer *httptest.ResponseRecorder
		var request *http.Request
		var errorWriter *fakes.ErrorWriter
		var notify *fakes.Notify

		BeforeEach(func() {
			var err error
			errorWriter = fakes.NewErrorWriter()
			writer = httptest.NewRecorder()
			request, err = http.NewRequest("POST", "/users/user-123", nil)
			if err != nil {
				panic(err)
			}

			notify = fakes.NewNotify()
			fakeDatabase := fakes.NewDatabase()
			handler = handlers.NewNotifyUser(notify, errorWriter, nil, fakeDatabase)
		})

		Context("when notify.Execute returns a successful response", func() {
			It("returns the JSON representation of the response", func() {
				notify.Response = []byte("whut")

				handler.Execute(writer, request, nil, nil, strategies.UserStrategy{})

				Expect(writer.Code).To(Equal(http.StatusOK))
				body := string(writer.Body.Bytes())
				Expect(body).To(Equal("whut"))

				Expect(notify.GUID.String()).To(Equal("user-123"))
			})
		})

		Context("when notify.Execute returns an error", func() {
			It("propagates the error", func() {
				notify.Error = errors.New("BOOM!")

				err := handler.Execute(writer, request, nil, nil, strategies.UserStrategy{})

				Expect(err).To(Equal(notify.Error))
			})
		})
	})
})
