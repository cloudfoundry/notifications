package info_test

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/web/info"
	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/ryanmoran/stack"

	. "github.com/cloudfoundry-incubator/notifications/testing/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	var muxer web.Muxer

	BeforeEach(func() {
		muxer = web.NewMuxer()
		info.Routes{
			RequestCounter: middleware.RequestCounter{},
			RequestLogging: middleware.RequestLogging{},
		}.Register(muxer)
	})

	It("routes GET /info", func() {
		request, err := http.NewRequest("GET", "/info", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(info.GetHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{})
	})
})
