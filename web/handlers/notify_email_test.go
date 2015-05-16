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

var _ = Describe("NotifyEmail", func() {
	Describe("Execute", func() {
		var (
			handler     handlers.NotifyEmail
			writer      *httptest.ResponseRecorder
			errorWriter *fakes.ErrorWriter
			notify      *fakes.Notify
			context     stack.Context
			connection  *fakes.DBConn
			request     *http.Request
			strategy    *fakes.MailStrategy
		)

		BeforeEach(func() {
			errorWriter = fakes.NewErrorWriter()
			writer = httptest.NewRecorder()
			connection = fakes.NewDBConn()
			request = &http.Request{}
			strategy = fakes.NewMailStrategy()

			context = stack.NewContext()
			context.Set(handlers.VCAPRequestIDKey, "some-request-id")

			notify = fakes.NewNotify()
			handler = handlers.NewNotifyEmail(notify, errorWriter, strategy, nil)
		})

		Context("when notify.Execute returns a proper response", func() {
			It("writes that response", func() {
				notify.ExecuteCall.Response = []byte("whut")

				handler.Execute(writer, nil, nil, context)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(writer.Body.String()).To(Equal("whut"))
			})

			It("delegates to the Notify object with the correct arguments", func() {
				handler.Execute(writer, request, connection, context)

				Expect(notify.ExecuteCall.Args.Connection).To(Equal(connection))
				Expect(notify.ExecuteCall.Args.Request).To(Equal(request))
				Expect(notify.ExecuteCall.Args.Context).To(Equal(context))
				Expect(notify.ExecuteCall.Args.GUID).To(Equal(""))
				Expect(notify.ExecuteCall.Args.Strategy).To(Equal(strategy))
				Expect(notify.ExecuteCall.Args.Validator).To(BeAssignableToTypeOf(params.EmailValidator{}))
				Expect(notify.ExecuteCall.Args.VCAPRequestID).To(Equal("some-request-id"))
			})
		})

		Context("when notify.Execute errors", func() {
			It("propagates the error", func() {
				notify.ExecuteCall.Error = errors.New("Blambo!")
				err := handler.Execute(writer, nil, nil, context)

				Expect(err).To(Equal(notify.ExecuteCall.Error))

			})
		})
	})
})
