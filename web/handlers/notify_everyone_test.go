package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotifyEveryone", func() {
	Context("Execute", func() {
		var (
			handler     handlers.NotifyEveryone
			writer      *httptest.ResponseRecorder
			request     *http.Request
			errorWriter *fakes.ErrorWriter
			notify      *fakes.Notify
			context     stack.Context
			connection  *fakes.DBConn
			strategy    *fakes.MailStrategy
		)

		BeforeEach(func() {
			errorWriter = fakes.NewErrorWriter()
			writer = httptest.NewRecorder()
			request = &http.Request{}
			strategy = fakes.NewMailStrategy()

			context = stack.NewContext()
			context.Set(handlers.VCAPRequestIDKey, "some-request-id")

			notify = fakes.NewNotify()
			handler = handlers.NewNotifyEveryone(notify, errorWriter, strategy, nil)
		})

		Context("when notify.Execute returns a successful response", func() {
			It("returns the JSON representation of the response", func() {
				notify.ExecuteCall.Response = []byte("hello")

				handler.Execute(writer, request, nil, context)

				Expect(writer.Code).To(Equal(http.StatusOK))
				body := string(writer.Body.Bytes())
				Expect(body).To(Equal("hello"))
			})

			It("delegates to the Notify object with the correct arguments", func() {
				handler.Execute(writer, request, connection, context)

				Expect(notify.ExecuteCall.Args.Connection).To(Equal(connection))
				Expect(notify.ExecuteCall.Args.Request).To(Equal(request))
				Expect(notify.ExecuteCall.Args.Context).To(Equal(context))
				Expect(notify.ExecuteCall.Args.GUID).To(Equal(""))
				Expect(notify.ExecuteCall.Args.Strategy).To(Equal(strategy))
				Expect(notify.ExecuteCall.Args.Validator).To(BeAssignableToTypeOf(params.GUIDValidator{}))
				Expect(notify.ExecuteCall.Args.VCAPRequestID).To(Equal("some-request-id"))
			})
		})

		Context("when notify.Execute returns an error", func() {
			It("propagates the error", func() {
				notify.ExecuteCall.Error = errors.New("BOOM!")

				err := handler.Execute(writer, request, nil, context)

				Expect(err).To(Equal(notify.ExecuteCall.Error))
			})
		})
	})
})
