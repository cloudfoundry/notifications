package messages_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v1/messages"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
	var router *mux.Router

	BeforeEach(func() {
		router = messages.NewRouter(messages.RouterConfig{
			RequestLogging:                               middleware.RequestLogging{},
			DatabaseAllocator:                            middleware.DatabaseAllocator{},
			NotificationsWriteOrEmailsWriteAuthenticator: middleware.Authenticator{Scopes: []string{"notifications.write", "emails.write"}},

			ErrorWriter:   fakes.NewErrorWriter(),
			MessageFinder: fakes.NewMessageFinder(),
		})
	})

	It("routes GET /messages/{message_id}", func() {
		s := router.Get("GET /messages/{message_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(messages.GetHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(ConsistOf([]string{"notifications.write", "emails.write"}))
	})
})
