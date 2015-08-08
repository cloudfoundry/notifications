package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmailHandler", func() {
	Describe("ServeHTTP", func() {
		var (
			handler     notify.EmailHandler
			writer      *httptest.ResponseRecorder
			errorWriter *fakes.ErrorWriter
			notifyObj   *fakes.Notify
			context     stack.Context
			connection  *fakes.Connection
			request     *http.Request
			strategy    *fakes.Strategy
			database    *fakes.Database
		)

		BeforeEach(func() {
			errorWriter = fakes.NewErrorWriter()
			writer = httptest.NewRecorder()
			database = fakes.NewDatabase()
			connection = fakes.NewConnection()
			database.Conn = connection
			request = &http.Request{}
			strategy = fakes.NewStrategy()

			context = stack.NewContext()
			context.Set(notify.VCAPRequestIDKey, "some-request-id")
			context.Set("database", database)

			notifyObj = fakes.NewNotify()
			handler = notify.NewEmailHandler(notifyObj, errorWriter, strategy)
		})

		Context("when notifyObj.Execute returns a proper response", func() {
			It("writes that response", func() {
				notifyObj.ExecuteCall.Response = []byte("whut")

				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(writer.Body.String()).To(Equal("whut"))
			})

			It("delegates to the notifyObj object with the correct arguments", func() {
				handler.ServeHTTP(writer, request, context)

				Expect(reflect.ValueOf(notifyObj.ExecuteCall.Args.Connection).Pointer()).To(Equal(reflect.ValueOf(connection).Pointer()))
				Expect(notifyObj.ExecuteCall.Args.Request).To(Equal(request))
				Expect(notifyObj.ExecuteCall.Args.Context).To(Equal(context))
				Expect(notifyObj.ExecuteCall.Args.GUID).To(Equal(""))
				Expect(notifyObj.ExecuteCall.Args.Strategy).To(Equal(strategy))
				Expect(notifyObj.ExecuteCall.Args.Validator).To(BeAssignableToTypeOf(notify.EmailValidator{}))
				Expect(notifyObj.ExecuteCall.Args.VCAPRequestID).To(Equal("some-request-id"))
			})
		})

		Context("when notifyObj.Execute errors", func() {
			It("propagates the error", func() {
				notifyObj.ExecuteCall.Error = errors.New("Blambo!")

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(notifyObj.ExecuteCall.Error))
			})
		})
	})
})
