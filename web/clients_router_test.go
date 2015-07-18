package web_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClientsRouter", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = web.NewClientsRouter(web.ClientsRouterConfig{
			TemplateAssigner:                 fakes.NewTemplateAssigner(),
			ErrorWriter:                      fakes.NewErrorWriter(),
			RequestLogging:                   web.RequestLogging{},
			NotificationsManageAuthenticator: web.Authenticator{Scopes: []string{"notifications.manage"}},
			DatabaseAllocator:                web.DatabaseAllocator{},
			NotificationsUpdater:             &fakes.NotificationUpdater{},
		})
	})

	It("routes PUT /clients/{client_id}/notifications/{notification_id}", func() {
		s := router.Get("PUT /clients/{client_id}/notifications/{notification_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdateNotifications{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{}, web.Authenticator{}, web.DatabaseAllocator{})

		authenticator := s.Middleware[2].(web.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})

	It("routes PUT /clients/{client_id}/template", func() {
		s := router.Get("PUT /clients/{client_id}/template").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.AssignClientTemplate{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{}, web.Authenticator{}, web.DatabaseAllocator{})

		authenticator := s.Middleware[2].(web.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})

	It("routes PUT /clients/{client_id}/notifications/{notification_id}/template", func() {
		s := router.Get("PUT /clients/{client_id}/notifications/{notification_id}/template").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.AssignNotificationTemplate{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{}, web.Authenticator{}, web.DatabaseAllocator{})

		authenticator := s.Middleware[2].(web.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})
})
