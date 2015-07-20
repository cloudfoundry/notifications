package notify_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/v1/notify"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EveryoneHandler", func() {
	Context("Execute", func() {
		var (
			handler     notify.EveryoneHandler
			writer      *httptest.ResponseRecorder
			request     *http.Request
			errorWriter *fakes.ErrorWriter
			notifyObj   *fakes.Notify
			context     stack.Context
			connection  *fakes.Connection
			strategy    *fakes.Strategy
		)

		BeforeEach(func() {
			errorWriter = fakes.NewErrorWriter()
			writer = httptest.NewRecorder()
			request = &http.Request{}
			strategy = fakes.NewStrategy()
			connection = fakes.NewConnection()
			database := fakes.NewDatabase()
			database.Conn = connection

			context = stack.NewContext()
			context.Set("database", database)
			context.Set(notify.VCAPRequestIDKey, "some-request-id")

			notifyObj = fakes.NewNotify()
			handler = notify.NewEveryoneHandler(notifyObj, errorWriter, strategy)
		})

		Context("when notifyObj.Execute returns a successful response", func() {
			It("returns the JSON representation of the response", func() {
				notifyObj.ExecuteCall.Response = []byte("hello")

				handler.ServeHTTP(writer, request, context)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(writer.Body.String()).To(Equal("hello"))
			})

			It("delegates to the notifyObj object with the correct arguments", func() {
				handler.ServeHTTP(writer, request, context)

				Expect(reflect.ValueOf(notifyObj.ExecuteCall.Args.Connection).Pointer()).To(Equal(reflect.ValueOf(connection).Pointer()))
				Expect(notifyObj.ExecuteCall.Args.Request).To(Equal(request))
				Expect(notifyObj.ExecuteCall.Args.Context).To(Equal(context))
				Expect(notifyObj.ExecuteCall.Args.GUID).To(Equal(""))
				Expect(notifyObj.ExecuteCall.Args.Strategy).To(Equal(strategy))
				Expect(notifyObj.ExecuteCall.Args.Validator).To(BeAssignableToTypeOf(notify.GUIDValidator{}))
				Expect(notifyObj.ExecuteCall.Args.VCAPRequestID).To(Equal("some-request-id"))
			})
		})

		Context("when notifyObj.Execute returns an error", func() {
			It("propagates the error", func() {
				notifyObj.ExecuteCall.Error = errors.New("BOOM!")

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(notifyObj.ExecuteCall.Error))
			})
		})
	})
})
