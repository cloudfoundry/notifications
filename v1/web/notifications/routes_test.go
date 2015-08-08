package notifications_test

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/web/notifications"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	var muxer web.Muxer

	BeforeEach(func() {
		muxer = web.NewMuxer()
		notifications.Routes{
			RequestLogging:                   middleware.RequestLogging{},
			DatabaseAllocator:                middleware.DatabaseAllocator{},
			NotificationsWriteAuthenticator:  middleware.Authenticator{Scopes: []string{"notifications.write"}},
			NotificationsManageAuthenticator: middleware.Authenticator{Scopes: []string{"notifications.manage"}},

			Registrar:            fakes.NewRegistrar(),
			ErrorWriter:          fakes.NewErrorWriter(),
			NotificationsFinder:  fakes.NewNotificationsFinder(),
			NotificationsUpdater: &fakes.NotificationUpdater{},
		}.Register(muxer)
	})

	Describe("/notifications", func() {
		It("routes PUT /notifications", func() {
			request, err := http.NewRequest("PUT", "/notifications", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(notifications.PutHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
		})

		It("routes GET /notifications", func() {
			request, err := http.NewRequest("GET", "/notifications", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(notifications.ListHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
		})

		It("routes PUT /clients/{client_id}/notifications/{notification_id}", func() {
			request, err := http.NewRequest("PUT", "/clients/{client_id}/notifications/{notification_id}", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(notifications.UpdateHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
		})

		It("routes PUT /clients/{client_id}/notifications/{notification_id}/template", func() {
			request, err := http.NewRequest("PUT", "/clients/{client_id}/notifications/{notification_id}/template", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(notifications.AssignTemplateHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
		})
	})

	Describe("/registration", func() {
		It("routes PUT /registration", func() {
			request, err := http.NewRequest("PUT", "/registration", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(notifications.RegistrationHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
		})
	})
})
