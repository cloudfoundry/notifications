package web_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificatonsRouter", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = web.NewNotificationsRouter(fakes.NewRegistrar(), fakes.NewErrorWriter(), middleware.RequestLogging{}, middleware.Authenticator{Scopes: []string{"notifications.write"}}, middleware.DatabaseAllocator{}, fakes.NewNotificationsFinder(), middleware.Authenticator{Scopes: []string{"notifications.manage"}})
	})

	Describe("/notifications", func() {
		It("routes PUT /notifications", func() {
			s := router.Get("PUT /notifications").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.RegisterClientWithNotifications{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
		})

		It("routes GET /notifications", func() {
			s := router.Get("GET /notifications").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetAllNotifications{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
		})
	})

	Describe("/registration", func() {
		It("routes PUT /registration", func() {
			s := router.Get("PUT /registration").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.RegisterNotifications{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
		})
	})
})
