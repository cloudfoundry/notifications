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

var _ = Describe("NotifyRouter", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = web.NewNotifyRouter(fakes.NewNotify(), fakes.NewErrorWriter(), fakes.NewStrategy(), web.RequestLogging{}, web.Authenticator{Scopes: []string{"notifications.write"}}, web.DatabaseAllocator{}, fakes.NewStrategy(), fakes.NewStrategy(), fakes.NewStrategy(), fakes.NewStrategy(), fakes.NewStrategy(), web.Authenticator{Scopes: []string{"emails.write"}})
	})

	It("routes POST /users/{user_id}", func() {
		s := router.Get("POST /users/{user_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyUser{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{}, web.Authenticator{}, web.DatabaseAllocator{})

		authenticator := s.Middleware[2].(web.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /spaces/{space_id}", func() {
		s := router.Get("POST /spaces/{space_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifySpace{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{}, web.Authenticator{}, web.DatabaseAllocator{})

		authenticator := s.Middleware[2].(web.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /organizations/{org_id}", func() {
		s := router.Get("POST /organizations/{org_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyOrganization{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{}, web.Authenticator{}, web.DatabaseAllocator{})

		authenticator := s.Middleware[2].(web.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /everyone", func() {
		s := router.Get("POST /everyone").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyEveryone{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{}, web.Authenticator{}, web.DatabaseAllocator{})

		authenticator := s.Middleware[2].(web.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /uaa_scopes/{scope}", func() {
		s := router.Get("POST /uaa_scopes/{scope}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyUAAScope{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{}, web.Authenticator{}, web.DatabaseAllocator{})

		authenticator := s.Middleware[2].(web.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /emails", func() {
		s := router.Get("POST /emails").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyEmail{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{}, web.Authenticator{}, web.DatabaseAllocator{})

		authenticator := s.Middleware[2].(web.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"emails.write"}))
	})
})
