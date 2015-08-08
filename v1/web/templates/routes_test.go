package templates_test

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
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
		templates.Routes{
			ErrorWriter:               fakes.NewErrorWriter(),
			TemplateFinder:            fakes.NewTemplateFinder(),
			TemplateUpdater:           fakes.NewTemplateUpdater(),
			TemplateCreator:           fakes.NewTemplateCreator(),
			TemplateDeleter:           fakes.NewTemplateDeleter(),
			TemplateLister:            fakes.NewTemplateLister(),
			TemplateAssociationLister: fakes.NewTemplateAssociationLister(),

			RequestLogging:                          middleware.RequestLogging{},
			DatabaseAllocator:                       middleware.DatabaseAllocator{},
			NotificationsManageAuthenticator:        middleware.Authenticator{Scopes: []string{"notifications.manage"}},
			NotificationTemplatesReadAuthenticator:  middleware.Authenticator{Scopes: []string{"notification_templates.read"}},
			NotificationTemplatesWriteAuthenticator: middleware.Authenticator{Scopes: []string{"notification_templates.write"}},
		}.Register(muxer)
	})

	Describe("/templates", func() {
		It("routes GET /templates", func() {
			request, err := http.NewRequest("GET", "/templates", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(templates.ListHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.read"}))
		})

		It("routes POST /templates", func() {
			request, err := http.NewRequest("POST", "/templates", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(templates.CreateHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
		})
	})

	Describe("/templates/{template_id}", func() {
		It("routes GET /templates/{template_id}", func() {
			request, err := http.NewRequest("GET", "/templates/{template_id}", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(templates.GetHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.read"}))
		})

		It("routes PUT /templates/{template_id}", func() {
			request, err := http.NewRequest("PUT", "/templates/{template_id}", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(templates.UpdateHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
		})

		It("routes DELETE /templates/{template_id}", func() {
			request, err := http.NewRequest("DELETE", "/templates/{template_id}", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(templates.DeleteHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
		})

		It("routes GET /templates/{template_id}/associations", func() {
			request, err := http.NewRequest("GET", "/templates/{template_id}/associations", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(templates.ListAssociationsHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
		})
	})

	Describe("/default_template", func() {
		It("routes GET /default_template", func() {
			request, err := http.NewRequest("GET", "/default_template", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(templates.GetDefaultHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.read"}))
		})

		It("routes PUT /default_template", func() {
			request, err := http.NewRequest("PUT", "/default_template", nil)
			Expect(err).NotTo(HaveOccurred())

			s := muxer.Match(request).(stack.Stack)
			Expect(s.Handler).To(BeAssignableToTypeOf(templates.UpdateDefaultHandler{}))
			ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

			authenticator := s.Middleware[2].(middleware.Authenticator)
			Expect(authenticator.Scopes).To(Equal([]string{"notification_templates.write"}))
		})
	})
})
