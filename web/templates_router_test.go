package web_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplatesRouter", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = web.NewTemplatesRouter(fakes.NewTemplateFinder(), fakes.NewErrorWriter(), middleware.RequestLogging{}, middleware.Authenticator{Scopes: []string{"notification_templates.read"}}, middleware.Authenticator{Scopes: []string{"notification_templates.write"}}, middleware.DatabaseAllocator{}, fakes.NewTemplateUpdater(), fakes.NewTemplateCreator(), fakes.NewTemplateDeleter(), fakes.NewTemplateAssociationLister(), middleware.Authenticator{Scopes: []string{"notifications.manage"}}, fakes.NewTemplateLister())
	})

	Describe("/templates", func() {
		It("routes GET /templates", func() {
			s := router.Get("GET /templates").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.ListTemplates{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.read"}))
		})

		It("routes POST /templates", func() {
			s := router.Get("POST /templates").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.CreateTemplate{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
		})
	})

	Describe("/templates/{template_id}", func() {
		It("routes GET /templates/{template_id}", func() {
			s := router.Get("GET /templates/{template_id}").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetTemplates{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.read"}))
		})

		It("routes PUT /templates/{template_id}", func() {
			s := router.Get("PUT /templates/{template_id}").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdateTemplates{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
		})

		It("routes DELETE /templates/{template_id}", func() {
			s := router.Get("DELETE /templates/{template_id}").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.DeleteTemplates{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
		})

		It("routes GET /templates/{template_id}/associations", func() {
			s := router.Get("GET /templates/{template_id}/associations").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.ListTemplateAssociations{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
		})
	})

	Describe("/default_template", func() {
		It("routes GET /default_template", func() {
			s := router.Get("GET /default_template").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetDefaultTemplate{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.read"}))
		})

		It("routes PUT /default_template", func() {
			s := router.Get("PUT /default_template").GetHandler().(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(handlers.UpdateDefaultTemplate{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
		})
	})
})
