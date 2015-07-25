package senders_test

import (
	"database/sql"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v2/senders"
	"github.com/gorilla/mux"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	var (
		logging     middleware.RequestLogging
		auth        middleware.Authenticator
		dbAllocator middleware.DatabaseAllocator
		router      *mux.Router
	)

	BeforeEach(func() {
		logging = middleware.NewRequestLogging(lager.NewLogger("log-prefix"))
		auth = middleware.NewAuthenticator("some-public-key", "notifications.write")
		dbAllocator = middleware.NewDatabaseAllocator(&sql.DB{}, false)

		router = mux.NewRouter()
		senders.Routes{
			RequestLogging:    logging,
			Authenticator:     auth,
			DatabaseAllocator: dbAllocator,
			SendersCollection: collections.SendersCollection{},
		}.Register(router)
	})

	It("routes POST /senders", func() {
		s := router.Get("POST /senders").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(senders.CreateHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})

	It("routes GET /senders/{sender_id}", func() {
		s := router.Get("GET /senders/{sender_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(senders.GetHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})
})
