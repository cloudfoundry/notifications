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

var _ = Describe("NotifyOrganization", func() {
	Describe("Execute", func() {
		var handler handlers.NotifyOrganization
		var writer *httptest.ResponseRecorder
		var request *http.Request
		var notify *fakes.Notify

		BeforeEach(func() {
			var err error

			writer = httptest.NewRecorder()
			request, err = http.NewRequest("POST", "/organizations/org-001", nil)
			if err != nil {
				panic(err)
			}

			notify = fakes.NewNotify()
			fakeDatabase := fakes.NewDatabase()
			handler = handlers.NewNotifyOrganization(notify, nil, nil, fakeDatabase)
		})

		Context("when the notify.Execute returns a successful response", func() {
			It("returns the JSON representation of the response", func() {
				notify.Response = []byte("whatever")
				strategy := strategies.OrganizationStrategy{}

				handler.Execute(writer, request, nil, nil, strategy)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(notify.GUID).To(Equal("org-001"))

				body := string(writer.Body.Bytes())
				Expect(body).To(Equal("whatever"))
			})
		})

		Context("when the notify.Execute returns an error", func() {
			It("propagates the error", func() {
				notify.Error = errors.New("the error")
				strategy := strategies.OrganizationStrategy{}

				err := handler.Execute(writer, request, nil, nil, strategy)
				Expect(err).To(Equal(notify.Error))
			})
		})
	})
})
