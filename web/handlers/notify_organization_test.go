package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/params"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotifyOrganization", func() {
	Describe("Execute", func() {
		var (
			handler    handlers.NotifyOrganization
			writer     *httptest.ResponseRecorder
			request    *http.Request
			notify     *fakes.Notify
			context    stack.Context
			connection *fakes.DBConn
			strategy   *fakes.Strategy
		)

		BeforeEach(func() {
			writer = httptest.NewRecorder()
			request = &http.Request{URL: &url.URL{Path: "/organizations/org-001"}}
			strategy = fakes.NewStrategy()

			context = stack.NewContext()
			context.Set(handlers.VCAPRequestIDKey, "some-request-id")

			notify = fakes.NewNotify()
			handler = handlers.NewNotifyOrganization(notify, nil, strategy, nil)
		})

		Context("when the notify.Execute returns a successful response", func() {
			It("returns the JSON representation of the response", func() {
				notify.ExecuteCall.Response = []byte("whatever")

				handler.Execute(writer, request, nil, context)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(writer.Body.String()).To(Equal("whatever"))
			})

			It("delegates to the Notify object with the correct arguments", func() {
				handler.Execute(writer, request, connection, context)

				Expect(notify.ExecuteCall.Args.Connection).To(Equal(connection))
				Expect(notify.ExecuteCall.Args.Request).To(Equal(request))
				Expect(notify.ExecuteCall.Args.Context).To(Equal(context))
				Expect(notify.ExecuteCall.Args.GUID).To(Equal("org-001"))
				Expect(notify.ExecuteCall.Args.Strategy).To(Equal(strategy))
				Expect(notify.ExecuteCall.Args.Validator).To(BeAssignableToTypeOf(params.GUIDValidator{}))
				Expect(notify.ExecuteCall.Args.VCAPRequestID).To(Equal("some-request-id"))
			})
		})

		Context("when the notify.Execute returns an error", func() {
			It("propagates the error", func() {
				notify.ExecuteCall.Error = errors.New("the error")

				err := handler.Execute(writer, request, nil, context)
				Expect(err).To(Equal(notify.ExecuteCall.Error))
			})
		})
	})
})
