package info_test

import (
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v1/info"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = mux.NewRouter()
		info.Routes{
			Version:        1,
			RequestLogging: middleware.RequestLogging{},
		}.Register(router)
	})

	It("routes GET /info", func() {
		s := router.Get("GET /info").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(info.GetHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{})
	})
})
