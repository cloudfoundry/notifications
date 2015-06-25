package web_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/services"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserPreferencesRouter", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = web.NewUserPreferencesRouter(middleware.RequestLogging{}, middleware.CORS{}, fakes.NewPreferencesFinder(services.PreferencesBuilder{}), fakes.NewErrorWriter(), middleware.Authenticator{Scopes: []string{"notification_preferences.read"}}, middleware.DatabaseAllocator{}, middleware.Authenticator{Scopes: []string{"notification_preferences.admin"}}, fakes.NewPreferenceUpdater(), middleware.Authenticator{Scopes: []string{"notification_preferences.write"}})
	})

	Describe("/user_preferences", func() {
		It("routes GET /user_preferences", func() {
			s := router.Get("GET /user_preferences").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetPreferences{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.read"}))
		})

		It("routes PATCH /user_preferences", func() {
			s := router.Get("PATCH /user_preferences").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdatePreferences{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.write"}))
		})

		It("routes OPTIONS /user_preferences", func() {
			s := router.Get("OPTIONS /user_preferences").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.OptionsPreferences{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{})
		})
	})

	Describe("/user_preferences/{user_id}", func() {
		It("routes GET /user_preferences/{user_id}", func() {
			s := router.Get("GET /user_preferences/{user_id}").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetPreferencesForUser{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.admin"}))
		})

		It("routes PATCH /user_preferences/{user_id}", func() {
			s := router.Get("PATCH /user_preferences/{user_id}").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdateSpecificUserPreferences{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[3].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.admin"}))
		})

		It("routes OPTIONS /user_preferences/{user_id}", func() {
			s := router.Get("OPTIONS /user_preferences/{user_id}").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.OptionsPreferences{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.CORS{})
		})
	})
})
