package web_test

import (
    "github.com/cloudfoundry-incubator/notifications/web"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    "github.com/cloudfoundry-incubator/notifications/web/middleware"
    "github.com/ryanmoran/stack"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
    var router web.Router

    BeforeEach(func() {
        router = web.NewRouter()
    })

    It("routes GET /info", func() {
        s := router.Routes().Get("GET /info").GetHandler().(stack.Stack)
        Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetInfo{}))
        Expect(s.Middleware[0]).To(BeAssignableToTypeOf(stack.Logging{}))
    })

    It("routes POST /users/{guid}", func() {
        s := router.Routes().Get("POST /users/{guid}").GetHandler().(stack.Stack)
        Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyUser{}))
        Expect(s.Middleware[0]).To(BeAssignableToTypeOf(stack.Logging{}))
        Expect(s.Middleware[1]).To(BeAssignableToTypeOf(middleware.Authenticator{}))
    })

    It("routes POST /spaces/{guid}", func() {
        s := router.Routes().Get("POST /spaces/{guid}").GetHandler().(stack.Stack)
        Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifySpace{}))
        Expect(s.Middleware[0]).To(BeAssignableToTypeOf(stack.Logging{}))
        Expect(s.Middleware[1]).To(BeAssignableToTypeOf(middleware.Authenticator{}))
    })

    It("routes PUT /registration", func() {
        s := router.Routes().Get("PUT /registration").GetHandler().(stack.Stack)
        Expect(s.Handler).To(BeAssignableToTypeOf(handlers.Registration{}))
        Expect(s.Middleware[0]).To(BeAssignableToTypeOf(stack.Logging{}))
        Expect(s.Middleware[1]).To(BeAssignableToTypeOf(middleware.Authenticator{}))
    })
})
