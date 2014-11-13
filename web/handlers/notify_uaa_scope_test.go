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

var _ = Describe("NotifyUAAScope", func() {
	Describe("Execute", func() {
		var notify *fakes.Notify
		var handler handlers.NotifyUAAScope
		var writer *httptest.ResponseRecorder
		var request *http.Request

		BeforeEach(func() {
			var err error

			writer = httptest.NewRecorder()
			request, err = http.NewRequest("POST", "/uaa_scopes/great.scope", nil)
			if err != nil {
				panic(err)
			}

			notify = fakes.NewNotify()
			fakeDatabase := fakes.NewDatabase()
			handler = handlers.NewNotifyUAAScope(notify, nil, nil, fakeDatabase)
		})

		Context("when the notify.Execute returns a successful response", func() {
			It("returns the JSON representation of the response", func() {
				notify.Response = []byte("whatever")
				strategy := strategies.UAAScopeStrategy{}

				handler.Execute(writer, request, nil, nil, strategy)

				Expect(writer.Code).To(Equal(http.StatusOK))
				Expect(notify.GUID.String()).To(Equal("great.scope"))

				body := string(writer.Body.Bytes())
				Expect(body).To(Equal("whatever"))
			})
		})

		Context("when notify.Execute returns an error", func() {
			It("Propagates the error", func() {
				notify.Error = errors.New("the error")
				strategy := strategies.UAAScopeStrategy{}

				err := handler.Execute(writer, request, nil, nil, strategy)
				Expect(err).To(Equal(notify.Error))
			})
		})
	})
})
