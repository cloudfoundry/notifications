package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"

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
			strategy    *fakes.Strategy
			database    *fakes.Database
		)

		BeforeEach(func() {
			errorWriter = fakes.NewErrorWriter()
			writer = httptest.NewRecorder()
			database = fakes.NewDatabase()
			connection = fakes.NewDBConn()
			database.Conn = connection
			request = &http.Request{}
			strategy = fakes.NewStrategy()

			context = stack.NewContext()
			context.Set(handlers.VCAPRequestIDKey, "some-request-id")
			context.Set("database", database)

			notify = fakes.NewNotify()
			handler = handlers.NewNotifyEmail(notify, errorWriter, strategy)
		})

		Context("when notify.Execute returns a proper response", func() {
			It("writes that response", func() {
				notify.ExecuteCall.Response = []byte("whut")

				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(writer.Body.String()).To(Equal("whut"))
			})

			It("delegates to the Notify object with the correct arguments", func() {
				handler.ServeHTTP(writer, request, context)

				Expect(reflect.ValueOf(notify.ExecuteCall.Args.Connection).Pointer()).To(Equal(reflect.ValueOf(connection).Pointer()))
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

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(notify.ExecuteCall.Error))
			})
		})
	})
})
