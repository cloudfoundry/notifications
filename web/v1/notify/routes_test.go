package notify_test

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v1/notify"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	var muxer web.Muxer

	BeforeEach(func() {
		muxer = web.NewMuxer()
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
		}.Register(muxer)
	})

	It("routes POST /users/{user_id}", func() {
		request, err := http.NewRequest("POST", "/users/{user_id}", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.UserHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /spaces/{space_id}", func() {
		request, err := http.NewRequest("POST", "/spaces/{space_id}", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.SpaceHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /organizations/{org_id}", func() {
		request, err := http.NewRequest("POST", "/organizations/{org_id}", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.OrganizationHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /everyone", func() {
		request, err := http.NewRequest("POST", "/everyone", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.EveryoneHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /uaa_scopes/{scope}", func() {
		request, err := http.NewRequest("POST", "/uaa_scopes/{scope}", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.UAAScopeHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /emails", func() {
		request, err := http.NewRequest("POST", "/emails", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notify.EmailHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"emails.write"}))
	})
})
