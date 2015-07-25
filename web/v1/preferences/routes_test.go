package preferences_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v1/preferences"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = mux.NewRouter()
		preferences.Routes{
			ErrorWriter:       fakes.NewErrorWriter(),
			PreferencesFinder: fakes.NewPreferencesFinder(services.PreferencesBuilder{}),
			PreferenceUpdater: fakes.NewPreferenceUpdater(),

			CORS:                                      middleware.CORS{},
			RequestLogging:                            middleware.RequestLogging{},
			DatabaseAllocator:                         middleware.DatabaseAllocator{},
			NotificationPreferencesReadAuthenticator:  middleware.Authenticator{Scopes: []string{"notification_preferences.read"}},
			NotificationPreferencesAdminAuthenticator: middleware.Authenticator{Scopes: []string{"notification_preferences.admin"}},
			NotificationPreferencesWriteAuthenticator: middleware.Authenticator{Scopes: []string{"notification_preferences.write"}},
		}.Register(router)
	})

	Describe("/user_preferences", func() {
		It("routes GET /user_preferences", func() {
			s := router.Get("GET /user_preferences").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.GetPreferencesHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.read"}))
		})

		It("routes PATCH /user_preferences", func() {
			s := router.Get("PATCH /user_preferences").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.UpdatePreferencesHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.write"}))
		})

		It("routes OPTIONS /user_preferences", func() {
			s := router.Get("OPTIONS /user_preferences").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.OptionsHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{})
		})
	})

	Describe("/user_preferences/{user_id}", func() {
		It("routes GET /user_preferences/{user_id}", func() {
			s := router.Get("GET /user_preferences/{user_id}").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.GetUserPreferencesHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.admin"}))
		})

		It("routes PATCH /user_preferences/{user_id}", func() {
			s := router.Get("PATCH /user_preferences/{user_id}").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.UpdateUserPreferencesHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.admin"}))
		})

		It("routes OPTIONS /user_preferences/{user_id}", func() {
			s := router.Get("OPTIONS /user_preferences/{user_id}").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(preferences.OptionsHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{})
		})
	})
})
