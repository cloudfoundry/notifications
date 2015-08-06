package campaigntypes_test

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v2/campaigntypes"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routes", func() {
	var (
		logging     middleware.RequestLogging
		dbAllocator middleware.DatabaseAllocator
		auth        middleware.Authenticator
		muxer       web.Muxer
	)

	BeforeEach(func() {
		logging = middleware.NewRequestLogging(lager.NewLogger("log-prefix"))
		auth = middleware.NewAuthenticator("some-public-key", "notifications.write")
		dbAllocator = middleware.NewDatabaseAllocator(&sql.DB{}, false)
		muxer = web.NewMuxer()
		campaigntypes.Routes{
			RequestLogging:          logging,
			Authenticator:           auth,
			DatabaseAllocator:       dbAllocator,
			CampaignTypesCollection: collections.CampaignTypesCollection{},
		}.Register(muxer)
	})

	It("routes POST /senders/{sender_id}/campaign_types", func() {
		request, err := http.NewRequest("POST", "/senders/some-sender-id/campaign_types", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(campaigntypes.CreateHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})

	It("routes GET /senders/{sender_id}/campaign_types", func() {
		request, err := http.NewRequest("GET", "/senders/some-sender-id/campaign_types", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(campaigntypes.ListHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})

	It("routes GET /senders/{sender_id}/campaign_types/{campaign_type_id}", func() {
		request, err := http.NewRequest("GET", "/senders/some-sender-id/campaign_types/some-campaign-type-id", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(campaigntypes.ShowHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})

	It("routes PUT /senders/{sender_id}/campaign_types/{campaign_type}", func() {
		request, err := http.NewRequest("PUT", "/senders/some-sender-id/campaign_types/some-campaign-type-id", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(campaigntypes.UpdateHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})

	It("routes DELETE /senders/{sender_id}/campaign_types/{campaign_type}", func() {
		request, err := http.NewRequest("DELETE", "/senders/some-sender-id/campaign_types/some-campaign-type-id", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(campaigntypes.DeleteHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})
})
