package strategies_test

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmailStrategy", func() {
	var emailStrategy strategies.EmailStrategy

	Describe("Dispatch", func() {
		var (
			mailer        *fakes.Mailer
			conn          *fakes.DBConn
			options       postal.Options
			clientID      string
			emailID       string
			vcapRequestID string
		)

		BeforeEach(func() {
			mailer = fakes.NewMailer()
			emailStrategy = strategies.NewEmailStrategy(mailer)

			clientID = "raptors-123"
			emailID = ""
			vcapRequestID = "some-request-id"

			options = postal.Options{
				Text: "email text",
				To:   "dr@strangelove.com",
			}

			conn = fakes.NewDBConn()
		})

		It("calls Deliver on it's mailer with proper arguments", func() {
			Expect(options.Endorsement).To(BeEmpty())

			emailStrategy.Dispatch(clientID, emailID, vcapRequestID, options, conn)
			options.Endorsement = strategies.EmailEndorsement

			users := []strategies.User{{Email: options.To}}
			Expect(mailer.DeliverArguments).To(Equal(map[string]interface{}{
				"connection":      conn,
				"users":           users,
				"options":         options,
				"space":           cf.CloudControllerSpace{},
				"org":             cf.CloudControllerOrganization{},
				"client":          clientID,
				"scope":           "",
				"vcap-request-id": vcapRequestID,
			}))
		})
	})
})
