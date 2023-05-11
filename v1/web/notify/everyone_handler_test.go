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

var _ = Describe("EveryoneHandler", func() {
	Context("Execute", func() {
		var (
			handler     notify.EveryoneHandler
			writer      *httptest.ResponseRecorder
			request     *http.Request
			errorWriter *mocks.ErrorWriter
			notifyObj   *mocks.Notify
			context     stack.Context
			connection  *mocks.Connection
			strategy    *mocks.Strategy
		)

		BeforeEach(func() {
			errorWriter = mocks.NewErrorWriter()
			writer = httptest.NewRecorder()
			request = &http.Request{}
			strategy = mocks.NewStrategy()

			connection = mocks.NewConnection()
			database := mocks.NewDatabase()
			database.ConnectionCall.Returns.Connection = connection

			context = stack.NewContext()
			context.Set("database", database)
			context.Set(notify.VCAPRequestIDKey, "some-request-id")

			notifyObj = mocks.NewNotify()
			handler = notify.NewEveryoneHandler(notifyObj, errorWriter, strategy)
		})

		Context("when notifyObj.Execute returns a successful response", func() {
			It("returns the JSON representation of the response", func() {
				notifyObj.ExecuteCall.Returns.Response = []byte("hello")

				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(writer.Body.String()).To(Equal("hello"))
			})

			It("delegates to the notifyObj object with the correct arguments", func() {
				handler.ServeHTTP(writer, request, context)

				Expect(reflect.ValueOf(notifyObj.ExecuteCall.Receives.Connection).Pointer()).To(Equal(reflect.ValueOf(connection).Pointer()))
				Expect(notifyObj.ExecuteCall.Receives.Request).To(Equal(request))
				Expect(notifyObj.ExecuteCall.Receives.Context).To(Equal(context))
				Expect(notifyObj.ExecuteCall.Receives.GUID).To(Equal(""))
				Expect(notifyObj.ExecuteCall.Receives.Strategy).To(Equal(strategy))
				Expect(notifyObj.ExecuteCall.Receives.Validator).To(BeAssignableToTypeOf(notify.GUIDValidator{}))
				Expect(notifyObj.ExecuteCall.Receives.VCAPRequestID).To(Equal("some-request-id"))
			})
		})

		Context("when notifyObj.Execute returns an error", func() {
			It("propagates the error", func() {
				notifyObj.ExecuteCall.Returns.Error = errors.New("BOOM!")

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(Equal(notifyObj.ExecuteCall.Returns.Error))
			})
		})
	})
})
