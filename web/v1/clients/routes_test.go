package clients_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v1/clients"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = mux.NewRouter()
		clients.Routes{
			RequestLogging:                   middleware.RequestLogging{},
			DatabaseAllocator:                middleware.DatabaseAllocator{},
			NotificationsManageAuthenticator: middleware.Authenticator{Scopes: []string{"notifications.manage"}},

			ErrorWriter:      fakes.NewErrorWriter(),
			TemplateAssigner: fakes.NewTemplateAssigner(),
		}.Register(router)
	})

	It("routes PUT /clients/{client_id}/template", func() {
		s := router.Get("PUT /clients/{client_id}/template").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(clients.AssignTemplateHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})
})
