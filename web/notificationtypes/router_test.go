package notificationtypes_test

import (
	"github.com/cloudfoundry-incubator/notifications/collections"
	"github.com/cloudfoundry-incubator/notifications/web/notificationtypes"
	"github.com/gorilla/mux"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationTypesRouter", func() {
	var (
		router *mux.Router
	)

	BeforeEach(func() {
		//logging = web.NewRequestLogging(lager.NewLogger("log-prefix"))
		//auth = web.NewAuthenticator("some-public-key", "notifications.write")
		//dbAllocator = web.NewDatabaseAllocator(&sql.DB{}, false)
		router = notificationtypes.NewRouter(notificationtypes.RouterConfig{
			NotificationTypesCollection: collections.NotificationTypesCollection{},
		})
	})

	It("routes POST /senders/{sender_id}/notification_types", func() {
		s := router.Get("POST /senders/{sender_id}/notification_types").GetHandler().(stack.Stack)
		Expect(s.Handler).To(BeAssignableToTypeOf(notificationtypes.CreateHandler{}))
		Expect(s.Middleware).To(HaveLen(0))
	})
})
