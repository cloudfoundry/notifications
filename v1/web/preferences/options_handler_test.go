package preferences_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/v1/web/preferences"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OptionsHandler", func() {
	var (
		handler preferences.OptionsHandler
		writer  *httptest.ResponseRecorder
		request *http.Request
		context stack.Context
	)

	BeforeEach(func() {
		var err error
		writer = httptest.NewRecorder()
		request, err = http.NewRequest("OPTIONS", "/user_preferences", nil)
		Expect(err).NotTo(HaveOccurred())
		context = stack.NewContext()
		handler = preferences.NewOptionsHandler()
	})

	Describe("ServeHTTP", func() {
		It("returns a 204 status code", func() {
			handler.ServeHTTP(writer, request, context)

			Expect(writer.Code).To(Equal(http.StatusNoContent))
		})
	})
})
