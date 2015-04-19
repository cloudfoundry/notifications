package web_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func ExpectToContainMiddlewareStack(actualMiddleware []stack.Middleware, expectedMiddleware ...stack.Middleware) {
	for i, ware := range expectedMiddleware {
		Expect(actualMiddleware[i]).To(BeAssignableToTypeOf(ware))
	}

}

var _ = Describe("Router", func() {
	var router web.Router

	BeforeEach(func() {
		router = web.NewRouter(fakes.NewMother())
	})

	It("routes GET /info", func() {
		s := router.Routes().Get("GET /info").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetInfo{}))
		Expect(s.Middleware).To(HaveLen(2))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{})
	})

	It("routes POST /users/{user_id}", func() {
		s := router.Routes().Get("POST /users/{user_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyUser{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /spaces/{space_id}", func() {
		s := router.Routes().Get("POST /spaces/{space_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifySpace{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /organizations/{org_id}", func() {
		s := router.Routes().Get("POST /organizations/{org_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyOrganization{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /everyone", func() {
		s := router.Routes().Get("POST /everyone").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyEveryone{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /uaa_scopes/{scope}", func() {
		s := router.Routes().Get("POST /uaa_scopes/{scope}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyUAAScope{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes POST /emails", func() {
		s := router.Routes().Get("POST /emails").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.NotifyEmail{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"emails.write"}))
	})

	It("routes PUT /registration", func() {
		s := router.Routes().Get("PUT /registration").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.RegisterNotifications{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes PUT /notifications", func() {
		s := router.Routes().Get("PUT /notifications").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.RegisterClientWithNotifications{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.write"}))
	})

	It("routes PUT /clients/{client_id}/notifications/{notification_id}", func() {
		s := router.Routes().Get("PUT /clients/{client_id}/notifications/{notification_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdateNotifications{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})

	It("routes GET /notifications", func() {
		s := router.Routes().Get("GET /notifications").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetAllNotifications{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})

	It("routes GET /user_preferences", func() {
		s := router.Routes().Get("GET /user_preferences").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetPreferences{}))
		Expect(s.Middleware).To(HaveLen(4))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{})

		authenticator := s.Middleware[3].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.read"}))
	})

	It("routes GET /user_preferences/{user_id}", func() {
		s := router.Routes().Get("GET /user_preferences/{user_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetPreferencesForUser{}))
		Expect(s.Middleware).To(HaveLen(4))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{})

		authenticator := s.Middleware[3].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.admin"}))
	})

	It("routes PATCH /user_preferences", func() {
		s := router.Routes().Get("PATCH /user_preferences").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdatePreferences{}))
		Expect(s.Middleware).To(HaveLen(4))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{})

		authenticator := s.Middleware[3].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.write"}))
	})

	It("routes OPTIONS /user_preferences", func() {
		s := router.Routes().Get("OPTIONS /user_preferences").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.OptionsPreferences{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.CORS{})
	})

	It("routes PATCH /user_preferences/{user_id}", func() {
		s := router.Routes().Get("PATCH /user_preferences/{user_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdateSpecificUserPreferences{}))
		Expect(s.Middleware).To(HaveLen(4))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.CORS{}, middleware.Authenticator{})

		authenticator := s.Middleware[3].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_preferences.admin"}))
	})

	It("routes OPTIONS /user_preferences/{user_id}", func() {
		s := router.Routes().Get("OPTIONS /user_preferences/{user_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.OptionsPreferences{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.CORS{})
	})

	It("routes POST /templates", func() {
		s := router.Routes().Get("POST /templates").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.CreateTemplate{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
	})

	It("routes GET /templates/{template_id}", func() {
		s := router.Routes().Get("GET /templates/{template_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetTemplates{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.read"}))
	})

	It("routes PUT /templates/{template_id}", func() {
		s := router.Routes().Get("PUT /templates/{template_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdateTemplates{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
	})

	It("routes DELETE /templates/{template_id}", func() {
		s := router.Routes().Get("DELETE /templates/{template_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.DeleteTemplates{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
	})

	It("routes GET /templates", func() {
		s := router.Routes().Get("GET /templates").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.ListTemplates{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.read"}))
	})

	It("routes PUT /clients/{client_id}/template", func() {
		s := router.Routes().Get("PUT /clients/{client_id}/template").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.AssignClientTemplate{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})

	It("routes PUT /clients/{client_id}/notifications/{notification_id}/template", func() {
		s := router.Routes().Get("PUT /clients/{client_id}/notifications/{notification_id}/template").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.AssignNotificationTemplate{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})

	It("routes GET /templates/{template_id}/associations", func() {
		s := router.Routes().Get("GET /templates/{template_id}/associations").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.ListTemplateAssociations{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})

	It("routes GET /default_template", func() {
		s := router.Routes().Get("GET /default_template").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetDefaultTemplate{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.read"}))
	})

	It("routes PUT /default_template", func() {
		s := router.Routes().Get("PUT /default_template").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdateDefaultTemplate{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
	})

	It("routes GET /messages/{message_id}", func() {
		s := router.Routes().Get("GET /messages/{message_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetMessages{}))
		Expect(s.Middleware).To(HaveLen(3))
		ExpectToContainMiddlewareStack(s.Middleware, stack.Logging{}, middleware.RequestCounter{}, middleware.Authenticator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(ConsistOf([]string{"notifications.write", "emails.write"}))
	})
})
