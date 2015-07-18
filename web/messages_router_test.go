package web_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessagesRouter", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = web.NewMessagesRouter(web.MessagesRouterConfig{
			MessageFinder:                                fakes.NewMessageFinder(),
			ErrorWriter:                                  fakes.NewErrorWriter(),
			RequestLogging:                               web.RequestLogging{},
			NotificationsWriteOrEmailsWriteAuthenticator: web.Authenticator{Scopes: []string{"notifications.write", "emails.write"}},
			DatabaseAllocator:                            web.DatabaseAllocator{},
		})
	})

	It("routes GET /messages/{message_id}", func() {
		s := router.Get("GET /messages/{message_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(handlers.GetMessages{}))
		ExpectToContainMiddlewareStack(s.Middleware, web.RequestLogging{}, web.RequestCounter{}, web.Authenticator{}, web.DatabaseAllocator{})

		authenticator := s.Middleware[2].(web.Authenticator)
		Expect(authenticator.Scopes).To(ConsistOf([]string{"notifications.write", "emails.write"}))
	})
})
