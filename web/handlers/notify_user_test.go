package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotifyUser", func() {
	Context("Execute", func() {
		var (
			handler     handlers.NotifyUser
			writer      *httptest.ResponseRecorder
			request     *http.Request
			notify      *fakes.Notify
			context     stack.Context
			connection  *fakes.Connection
			strategy    *fakes.Strategy
			errorWriter *fakes.ErrorWriter
		)

		BeforeEach(func() {
			writer = httptest.NewRecorder()
			request = &http.Request{URL: &url.URL{Path: "/users/user-123"}}
			strategy = fakes.NewStrategy()
			database := fakes.NewDatabase()
			connection = fakes.NewConnection()
			database.Conn = connection
			errorWriter = fakes.NewErrorWriter()

			context = stack.NewContext()
			context.Set("database", database)
			context.Set(handlers.VCAPRequestIDKey, "some-request-id")

			notify = fakes.NewNotify()
			handler = handlers.NewNotifyUser(notify, errorWriter, strategy)
		})

		Context("when notify.Execute returns a successful response", func() {
			It("returns the JSON representation of the response", func() {
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
				Expect(notify.ExecuteCall.Args.GUID).To(Equal("user-123"))
				Expect(notify.ExecuteCall.Args.Strategy).To(Equal(strategy))
				Expect(notify.ExecuteCall.Args.Validator).To(BeAssignableToTypeOf(handlers.GUIDValidator{}))
				Expect(notify.ExecuteCall.Args.VCAPRequestID).To(Equal("some-request-id"))
			})
		})

		Context("when notify.Execute returns an error", func() {
			It("propagates the error", func() {
				notify.ExecuteCall.Error = errors.New("BOOM!")
				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.Error).To(Equal(notify.ExecuteCall.Error))
			})
		})
	})
})
