package clients_test

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/web/clients"
	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
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
		clients.Routes{
			RequestCounter:                   middleware.RequestCounter{},
			RequestLogging:                   middleware.RequestLogging{},
			DatabaseAllocator:                middleware.DatabaseAllocator{},
			NotificationsManageAuthenticator: middleware.Authenticator{Scopes: []string{"notifications.manage"}},

			ErrorWriter:      mocks.NewErrorWriter(),
			TemplateAssigner: mocks.NewTemplateAssigner(),
		}.Register(muxer)
	})

	It("routes PUT /clients/{client_id}/template", func() {
		request, err := http.NewRequest("PUT", "/clients/some-client-id/template", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(clients.AssignTemplateHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(Equal([]string{"notifications.manage"}))
	})
})
