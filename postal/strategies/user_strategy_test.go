package strategies_test

import (
	"reflect"

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
		conn          *fakes.Connection
		vcapRequestID string
	)

	BeforeEach(func() {
		clientID = "mister-client"
		vcapRequestID = "some-request-id"
		conn = fakes.NewConnection()

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

			Expect(reflect.ValueOf(mailer.DeliverCall.Args.Connection).Pointer()).To(Equal(reflect.ValueOf(conn).Pointer()))
			Expect(mailer.DeliverCall.Args.Users).To(Equal(users))
			Expect(mailer.DeliverCall.Args.Options).To(Equal(options))
			Expect(mailer.DeliverCall.Args.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(mailer.DeliverCall.Args.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(mailer.DeliverCall.Args.Client).To(Equal(clientID))
			Expect(mailer.DeliverCall.Args.Scope).To(Equal(""))
			Expect(mailer.DeliverCall.Args.VCAPRequestID).To(Equal(vcapRequestID))
		})
	})
})
