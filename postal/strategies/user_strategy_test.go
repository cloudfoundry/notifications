package strategies_test

import (
	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserStrategy", func() {
	var (
		strategy      strategies.UserStrategy
		options       postal.Options
		mailer        *fakes.Mailer
		clientID      string
		conn          *fakes.DBConn
		vcapRequestID string
	)

	BeforeEach(func() {
		clientID = "mister-client"
		vcapRequestID = "some-request-id"

		mailer = fakes.NewMailer()
		strategy = strategies.NewUserStrategy(mailer)
	})

	Describe("Dispatch", func() {
		BeforeEach(func() {
			options = postal.Options{
				KindID:            "forgot_password",
				KindDescription:   "Password reminder",
				SourceDescription: "Login system",
				Text:              "Please reset your password by clicking on this link...",
				HTML: postal.HTML{
					BodyContent: "<p>Please reset your password by clicking on this link...</p>",
				},
			}
		})

		It("calls mailer.Deliver with the correct arguments for a user", func() {
			Expect(options.Endorsement).To(BeEmpty())

			_, err := strategy.Dispatch(clientID, "user-123", vcapRequestID, options, conn)
			if err != nil {
				panic(err)
			}
			users := []strategies.User{{GUID: "user-123"}}
			options.Endorsement = strategies.UserEndorsement

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
