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

var _ = Describe("ClientsRouter", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = web.NewClientsRouter(fakes.NewTemplateAssigner(), fakes.NewErrorWriter(), middleware.RequestLogging{}, middleware.Authenticator{Scopes: []string{"notifications.manage"}}, middleware.DatabaseAllocator{}, &fakes.NotificationUpdater{})
	})

	It("routes PUT /clients/{client_id}/notifications/{notification_id}", func() {
		s := router.Get("PUT /clients/{client_id}/notifications/{notification_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdateNotifications{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})

	It("routes PUT /clients/{client_id}/template", func() {
		s := router.Get("PUT /clients/{client_id}/template").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.AssignClientTemplate{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})

	It("routes PUT /clients/{client_id}/notifications/{notification_id}/template", func() {
		s := router.Get("PUT /clients/{client_id}/notifications/{notification_id}/template").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.AssignNotificationTemplate{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})
})
