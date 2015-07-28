package campaigntypes_test

import (
	"database/sql"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v2/campaigntypes"
	"github.com/gorilla/mux"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CampaignTypesRouter", func() {
	var (
		logging     middleware.RequestLogging
		dbAllocator middleware.DatabaseAllocator
		auth        middleware.Authenticator
		router      *mux.Router
	)

	BeforeEach(func() {
		logging = middleware.NewRequestLogging(lager.NewLogger("log-prefix"))
		auth = middleware.NewAuthenticator("some-public-key", "notifications.write")
		dbAllocator = middleware.NewDatabaseAllocator(&sql.DB{}, false)
		router = campaigntypes.NewRouter(campaigntypes.RouterConfig{
			RequestLogging:              logging,
			Authenticator:               auth,
			DatabaseAllocator:           dbAllocator,
			CampaignTypesCollection: collections.CampaignTypesCollection{},
		})
	})

	It("routes POST /senders/{sender_id}/campaign_types", func() {
		s := router.Get("POST /senders/{sender_id}/campaign_types").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(campaigntypes.CreateHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})

	It("routes GET /senders/{sender_id}/campaign_types/{campaign_type_id}", func() {
		s := router.Get("GET /senders/{sender_id}/campaign_types/{campaign_type_id}").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(campaigntypes.ShowHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})

	It("routes GET /senders/{sender_id}/campaign_types", func() {
		s := router.Get("GET /senders/{sender_id}/campaign_types").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(campaigntypes.ListHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})
})
