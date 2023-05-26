package preferences_test

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/v1/web/preferences"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/ryanmoran/stack"

	. "github.com/cloudfoundry-incubator/notifications/testing/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	var muxer web.Muxer

	BeforeEach(func() {
		muxer = web.NewMuxer()
		preferences.Routes{
			ErrorWriter:       mocks.NewErrorWriter(),
			PreferencesFinder: mocks.NewPreferencesFinder(),
			PreferenceUpdater: mocks.NewPreferenceUpdater(),

			CORS:                                     middleware.CORS{},
			RequestCounter:                           middleware.RequestCounter{},
			RequestLogging:                           middleware.RequestLogging{},
			DatabaseAllocator:                        middleware.DatabaseAllocator{},
			NotificationPreferencesReadAuthenticator: middleware.Authenticator{Scopes: []string{"notification_preferences.read"}},
			NotificationPreferencesAdminAuthenticator: middleware.Authenticator{Scopes: []string{"notification_preferences.admin"}},
			NotificationPreferencesWriteAuthenticator: middleware.Authenticator{Scopes: []string{"notification_preferences.write"}},
		}.Register(muxer)
	})

	Describe("/user_preferences", func() {
		It("routes GET /user_preferences", func() {
			request, err := http.NewRequest("GET", "/user_preferences", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.GetPreferencesHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.read"}))
		})

		It("routes PATCH /user_preferences", func() {
			request, err := http.NewRequest("PATCH", "/user_preferences", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.UpdatePreferencesHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.write"}))
		})

		It("routes OPTIONS /user_preferences", func() {
			request, err := http.NewRequest("OPTIONS", "/user_preferences", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.OptionsHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{})
		})
	})

	Describe("/user_preferences/{user_id}", func() {
		It("routes GET /user_preferences/{user_id}", func() {
			request, err := http.NewRequest("GET", "/user_preferences/some-user-id", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.GetUserPreferencesHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.admin"}))
		})

		It("routes PATCH /user_preferences/{user_id}", func() {
			request, err := http.NewRequest("PATCH", "/user_preferences/some-user-id", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.UpdateUserPreferencesHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.admin"}))
		})

		It("routes OPTIONS /user_preferences/{user_id}", func() {
			request, err := http.NewRequest("OPTIONS", "/user_preferences/some-user-id", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.OptionsHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{})
		})
	})
})
