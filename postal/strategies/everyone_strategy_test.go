package strategies_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Everyone Strategy", func() {
	var (
		strategy      strategies.EveryoneStrategy
		options       postal.Options
		tokenLoader   *fakes.TokenLoader
		allUsers      *fakes.AllUsers
		mailer        *fakes.Mailer
		clientID      string
		conn          *fakes.DBConn
		vcapRequestID string
	)

	BeforeEach(func() {
		clientID = "my-client"
		vcapRequestID = "some-request-id"
		conn = fakes.NewDBConn()

		tokenHeader := map[string]interface{}{
			"alg": "FAST",
		}

		tokenClaims := map[string]interface{}{
			"client_id": "mister-client",
			"exp":       int64(3404281214),
			"scope":     []string{"notifications.write"},
		}
		tokenLoader = fakes.NewTokenLoader()
		tokenLoader.Token = fakes.BuildToken(tokenHeader, tokenClaims)

		mailer = fakes.NewMailer()
		allUsers = fakes.NewAllUsers()
		allUsers.GUIDs = []string{"user-380", "user-319"}

		strategy = strategies.NewEveryoneStrategy(tokenLoader, allUsers, mailer)
	})

	Describe("Dispatch", func() {
		BeforeEach(func() {
			options = postal.Options{
				KindID:            "welcome_user",
				KindDescription:   "Your Official Welcome",
				SourceDescription: "Welcome system",
				Text:              "Welcome to the system, now get off my lawn.",
				HTML:              postal.HTML{BodyContent: "<p>Welcome to the system, now get off my lawn.</p>"},
			}
		})

		It("call mailer.Deliver with the correct arguments for an organization", func() {
			Expect(options.Endorsement).To(BeEmpty())
			_, err := strategy.Dispatch(clientID, "", vcapRequestID, options, conn)
			if err != nil {
				panic(err)
			}

			options.Endorsement = strategies.EveryoneEndorsement
			var users []strategies.User
			for _, guid := range allUsers.GUIDs {
				users = append(users, strategies.User{GUID: guid})
			}

			Expect(mailer.DeliverCall.Args.Connection).To(Equal(conn))
			Expect(mailer.DeliverCall.Args.Users).To(Equal(users))
			Expect(mailer.DeliverCall.Args.Options).To(Equal(options))
			Expect(mailer.DeliverCall.Args.Space).To(Equal(cf.CloudControllerSpace{}))
			Expect(mailer.DeliverCall.Args.Org).To(Equal(cf.CloudControllerOrganization{}))
			Expect(mailer.DeliverCall.Args.Client).To(Equal(clientID))
			Expect(mailer.DeliverCall.Args.Scope).To(Equal(""))
			Expect(mailer.DeliverCall.Args.VCAPRequestID).To(Equal(vcapRequestID))
		})
	})

	Context("failure cases", func() {
		Context("when token loader fails to return a token", func() {
			It("returns an error", func() {
				tokenLoader.LoadError = errors.New("BOOM!")
				_, err := strategy.Dispatch(clientID, "", vcapRequestID, options, conn)

				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})

		Context("when allUsers fails to load users", func() {
			It("returns the error", func() {
				allUsers.LoadError = errors.New("BOOM!")
				_, err := strategy.Dispatch(clientID, "", vcapRequestID, options, conn)

				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})
	})
})
