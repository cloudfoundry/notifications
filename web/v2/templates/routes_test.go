package templates_test

import (
	"database/sql"
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/web"
	"github.com/cloudfoundry-incubator/notifications/web/middleware"
	"github.com/cloudfoundry-incubator/notifications/web/v2/templates"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager"
	"github.com/ryanmoran/stack"
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
		templates.Routes{
			RequestLogging:      logging,
			Authenticator:       auth,
			DatabaseAllocator:   dbAllocator,
			TemplatesCollection: "",
		}.Register(muxer)
	})

	It("routes POST /templates", func() {
		request, err := http.NewRequest("POST", "/templates", nil)
		Expect(err).NotTo(HaveOccurred())

		s := muxer.Match(request).(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(templates.CreateHandler{}))
		Expect(s.Middleware).To(HaveLen(3))

		requestLogging := s.Middleware[0].(middleware.RequestLogging)
		Expect(requestLogging).To(Equal(logging))

		authenticator := s.Middleware[1].(middleware.Authenticator)
		Expect(authenticator).To(Equal(auth))

		databaseAllocator := s.Middleware[2].(middleware.DatabaseAllocator)
		Expect(databaseAllocator).To(Equal(dbAllocator))
	})
})
