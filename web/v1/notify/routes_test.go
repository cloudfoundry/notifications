package notify_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v1/notify"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = mux.NewRouter()
		notify.Routes{
			Notify:               fakes.NewNotify(),
			ErrorWriter:          fakes.NewErrorWriter(),
			UserStrategy:         fakes.NewStrategy(),
			SpaceStrategy:        fakes.NewStrategy(),
			OrganizationStrategy: fakes.NewStrategy(),
			EveryoneStrategy:     fakes.NewStrategy(),
			UAAScopeStrategy:     fakes.NewStrategy(),
			EmailStrategy:        fakes.NewStrategy(),

			RequestLogging:                  middleware.RequestLogging{},
			DatabaseAllocator:               middleware.DatabaseAllocator{},
			NotificationsWriteAuthenticator: middleware.Authenticator{Scopes: []string{"notifications.write"}},
			EmailsWriteAuthenticator:        middleware.Authenticator{Scopes: []string{"emails.write"}},
		}.Register(router)
	})

	It("routes POST /users/{user_id}", func() {
		s := router.Get("POST /users/{user_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.UserHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /spaces/{space_id}", func() {
		s := router.Get("POST /spaces/{space_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.SpaceHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /organizations/{org_id}", func() {
		s := router.Get("POST /organizations/{org_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.OrganizationHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /everyone", func() {
		s := router.Get("POST /everyone").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.EveryoneHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /uaa_scopes/{scope}", func() {
		s := router.Get("POST /uaa_scopes/{scope}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.UAAScopeHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /emails", func() {
		s := router.Get("POST /emails").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.EmailHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"emails.write"}))
	})
})
