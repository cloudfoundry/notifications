package unsubscribers_test

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/v2/collections"
	"github.com/cloudfoundry-incubator/notifications/v2/web/unsubscribers"
	"github.com/cloudfoundry-incubator/notifications/web"
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
		unsubscribers.Routes{
			RequestLogging:          logging,
			Authenticator:           auth,
			DatabaseAllocator:       dbAllocator,
			UnsubscribersCollection: collections.UnsubscribersCollection{},
		}.Register(muxer)
	})

	It("routes PUT /senders/{sender_id}/campaign_types/{campaign_type_id}/unsubscribers/{user_guid}", func() {
		request, err := http.NewRequest("PUT", "/senders/some-sender-id/campaign_types/some-campaign-type-id/unsubscribers/some-user-guid", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(unsubscribers.UpdateHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})
})
