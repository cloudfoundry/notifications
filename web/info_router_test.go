package web_test

import (
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InfoRouter", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = web.NewInfoRouter(1, web.RequestLogging{})
	})

	It("routes GET /info", func() {
		s := router.Get("GET /info").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetInfo{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{})
	})
})
