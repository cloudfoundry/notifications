package strategies_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/strategies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Everyone Strategy", func() {
	var (
		strategy            strategies.EveryoneStrategy
		options             postal.Options
		tokenLoader         *fakes.TokenLoader
		allUsers            *fakes.AllUsers
		mailer              *fakes.Mailer
		clientID            string
		conn                *fakes.Connection
		vcapRequestID       string
		requestReceivedTime time.Time
	)

	BeforeEach(func() {
		clientID = "my-client"
		vcapRequestID = "some-request-id"
		requestReceivedTime, _ = time.Parse(time.RFC3339Nano, "2015-06-08T14:38:03.180764129-07:00")
		conn = fakes.NewConnection()

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
		allUsers.AllUserGUIDsCall.Returns = []string{"user-380", "user-319"}

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
			_, err := strategy.Dispatch(clientID, "", vcapRequestID, requestReceivedTime, options, conn)
			if err != nil {
				panic(err)
			}

			options.Endorsement = strategies.EveryoneEndorsement
			var users []strategies.User
			for _, guid := range allUsers.AllUserGUIDsCall.Returns {
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
			Expect(mailer.DeliverCall.Args.RequestReceived).To(Equal(requestReceivedTime))
		})
	})

	Context("failure cases", func() {
		Context("when token loader fails to return a token", func() {
			It("returns an error", func() {
				tokenLoader.LoadError = errors.New("BOOM!")
				_, err := strategy.Dispatch(clientID, "", vcapRequestID, requestReceivedTime, options, conn)

				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})

		Context("when allUsers fails to load users", func() {
			It("returns the error", func() {
				allUsers.AllUserGUIDsCall.Error = errors.New("BOOM!")
				_, err := strategy.Dispatch(clientID, "", vcapRequestID, requestReceivedTime, options, conn)

				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})
	})
})
