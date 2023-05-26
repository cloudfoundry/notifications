package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmailHandler", func() {
	Describe("ServeHTTP", func() {
		var (
			handler     notify.EmailHandler
			writer      *httptest.ResponseRecorder
			errorWriter *mocks.ErrorWriter
			notifyObj   *mocks.Notify
			context     stack.Context
			connection  *mocks.Connection
			request     *http.Request
			strategy    *mocks.Strategy
			database    *mocks.Database
		)

		BeforeEach(func() {
			errorWriter = mocks.NewErrorWriter()
			writer = httptest.NewRecorder()

			connection = mocks.NewConnection()
			database = mocks.NewDatabase()
			database.ConnectionCall.Returns.Connection = connection

			request = &http.Request{}
			strategy = mocks.NewStrategy()

			context = stack.NewContext()
			context.Set(notify.VCAPRequestIDKey, "some-request-id")
			context.Set("database", database)

			notifyObj = mocks.NewNotify()
			handler = notify.NewEmailHandler(notifyObj, errorWriter, strategy)
		})

		Context("when notifyObj.Execute returns a proper response", func() {
			It("writes that response", func() {
				notifyObj.ExecuteCall.Returns.Response = []byte("whut")

				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(writer.Body.String()).To(Equal("whut"))
			})

			It("delegates to the notifyObj object with the correct arguments", func() {
				handler.ServeHTTP(writer, request, context)

				Expect(reflect.ValueOf(notifyObj.ExecuteCall.Receives.Connection).Pointer()).To(Equal(reflect.ValueOf(connection).Pointer()))
				Expect(notifyObj.ExecuteCall.Receives.Request).To(Equal(request))
				Expect(notifyObj.ExecuteCall.Receives.Context).To(Equal(context))
				Expect(notifyObj.ExecuteCall.Receives.GUID).To(Equal(""))
				Expect(notifyObj.ExecuteCall.Receives.Strategy).To(Equal(strategy))
				Expect(notifyObj.ExecuteCall.Receives.Validator).To(BeAssignableToTypeOf(notify.EmailValidator{}))
				Expect(notifyObj.ExecuteCall.Receives.VCAPRequestID).To(Equal("some-request-id"))
			})
		})

		Context("when notifyObj.Execute errors", func() {
			It("propagates the error", func() {
				notifyObj.ExecuteCall.Returns.Error = errors.New("Blambo!")

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(notifyObj.ExecuteCall.Returns.Error))
			})
		})
	})
})
