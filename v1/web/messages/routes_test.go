package messages_test

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/web/messages"
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
		messages.Routes{
			RequestLogging:                               middleware.RequestLogging{},
			DatabaseAllocator:                            middleware.DatabaseAllocator{},
			NotificationsWriteOrEmailsWriteAuthenticator: middleware.Authenticator{Scopes: []string{"notifications.write", "emails.write"}},

			ErrorWriter:   fakes.NewErrorWriter(),
			MessageFinder: fakes.NewMessageFinder(),
		}.Register(muxer)
	})

	It("routes GET /messages/{message_id}", func() {
		request, err := http.NewRequest("GET", "/messages/some-message-id", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(messages.GetHandler{}))
		ExpectToContainMiddlewareStack(s.Middleware, middleware.RequestLogging{}, middleware.RequestCounter{}, middleware.Authenticator{}, middleware.DatabaseAllocator{})

		authenticator := s.Middleware[2].(middleware.Authenticator)
		Expect(authenticator.Scopes).To(ConsistOf([]string{"notifications.write", "emails.write"}))
	})
})
