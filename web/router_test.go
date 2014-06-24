package web_test

import (
    "github.com/pivotal-cf/cf-notifications/web"
    "github.com/pivotal-cf/cf-notifications/web/handlers"
    "github.com/pivotal-cf/cf-notifications/web/middleware"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
    var router web.Router

    BeforeEach(func() {
        router = web.NewRouter()
    })

    It("routes GET /info", func() {
        stack := router.Routes().Get("GET /info").GetHandler().(web.Stack)
        Expect(stack.Handler).To(BeAssignableToTypeOf(handlers.GetInfo{}))
        Expect(stack.Middleware[0]).To(BeAssignableToTypeOf(middleware.Logging{}))
    })
})
