package notifications_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v1/notifications"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = notifications.NewRouter(notifications.RouterConfig{
			RequestLogging:                   middleware.RequestLogging{},
			DatabaseAllocator:                middleware.DatabaseAllocator{},
			NotificationsWriteAuthenticator:  middleware.Authenticator{Scopes: []string{"notifications.write"}},
			NotificationsManageAuthenticator: middleware.Authenticator{Scopes: []string{"notifications.manage"}},

			Registrar:            fakes.NewRegistrar(),
			ErrorWriter:          fakes.NewErrorWriter(),
			NotificationsFinder:  fakes.NewNotificationsFinder(),
			NotificationsUpdater: &fakes.NotificationUpdater{},
		})
	})

	Describe("/notifications", func() {
		It("routes PUT /notifications", func() {
			s := router.Get("PUT /notifications").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(notifications.PutHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
		})

		It("routes GET /notifications", func() {
			s := router.Get("GET /notifications").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(notifications.ListHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
		})

		It("routes PUT /clients/{client_id}/notifications/{notification_id}", func() {
			s := router.Get("PUT /clients/{client_id}/notifications/{notification_id}").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(notifications.UpdateHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
		})

		It("routes PUT /clients/{client_id}/notifications/{notification_id}/template", func() {
			s := router.Get("PUT /clients/{client_id}/notifications/{notification_id}/template").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(notifications.AssignTemplateHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
		})
	})

	Describe("/registration", func() {
		It("routes PUT /registration", func() {
			s := router.Get("PUT /registration").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(notifications.RegistrationHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
		})
	})
})
