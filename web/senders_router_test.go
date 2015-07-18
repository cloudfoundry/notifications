package web_test

import (
	"database/sql"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/senders"
	"github.com/gorilla/mux"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SendersRouter", func() {
	var (
		logging     web.RequestLogging
		auth        web.Authenticator
		dbAllocator web.DatabaseAllocator
		router      *mux.Router
	)

	BeforeEach(func() {
		logging = web.NewRequestLogging(lager.NewLogger("log-prefix"))
		auth = web.NewAuthenticator("some-public-key", "notifications.write")
		dbAllocator = web.NewDatabaseAllocator(&sql.DB{}, false)
		router = web.NewSendersRouter(web.SendersRouterConfig{
			RequestLogging:    logging,
			Authenticator:     auth,
			DatabaseAllocator: dbAllocator,
			SendersCollection: collections.SendersCollection{},
		})
	})

	It("routes POST /senders", func() {
		s := router.Get("POST /senders").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(senders.CreateHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(web.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(web.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(web.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})

	It("routes GET /senders/{sender_id}", func() {
		s := router.Get("GET /senders/{sender_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(senders.GetHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(web.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(web.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(web.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})
})
