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
		var mailer *fakes.Mailer
		var conn *fakes.DBConn
		var options postal.Options
		var clientID string
		var emailID string

		BeforeEach(func() {
			mailer = fakes.NewMailer()
			emailStrategy = strategies.NewEmailStrategy(mailer)

			clientID = "raptors-123"
			emailID = ""

			options = postal.Options{
				Text: "email text",
				To:   "dr@strangelove.com",
			}

			conn = fakes.NewDBConn()
		})

		It("Calls Deliver on it's mailer with proper arguments", func() {
			Expect(options.Endorsement).To(BeEmpty())

			emailStrategy.Dispatch(clientID, emailID, options, conn)
			options.Endorsement = strategies.EmailEndorsement

			users := []strategies.User{{Email: options.To}}
			Expect(mailer.DeliverArguments).To(Equal(map[string]interface{}{
				"connection": conn,
				"users":      users,
				"options":    options,
				"space":      cf.CloudControllerSpace{},
				"org":        cf.CloudControllerOrganization{},
				"client":     clientID,
				"scope":      "",
			}))
		})
	})
})
