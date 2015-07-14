package acceptance

import (
	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Passwords", func() {
	var (
		client warrant.Warrant
		user   warrant.User
	)

	BeforeEach(func() {
		client = warrant.New(warrant.Config{
			Host:          UAAHost,
			SkipVerifySSL: true,
			TraceWriter:   TraceWriter,
		})

	})

	AfterEach(func() {
		err := client.Users.Delete(user.ID, UAAToken)
		Expect(err).NotTo(HaveOccurred())
	})

	It("allows a user password to be set/updated", func() {
		var (
			userToken string
		)

		By("creating a new user", func() {
			var err error
			user, err = client.Users.Create(UAADefaultUsername, "warrant-user@example.com", UAAToken)
			Expect(err).NotTo(HaveOccurred())
		})

		By("setting the user password using a valid client", func() {
			err := client.Users.SetPassword(user.ID, "password", UAAToken)
			Expect(err).NotTo(HaveOccurred())
		})

		By("retrieving the user token using the new password", func() {
			var err error
			userToken, err = client.Users.GetToken(user.UserName, "password")
			Expect(err).NotTo(HaveOccurred())

			decodedToken, err := client.Tokens.Decode(userToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(decodedToken.UserID).To(Equal(user.ID))
		})

		By("changing a user's own password", func() {
			err := client.Users.ChangePassword(user.ID, "password", "new-password", userToken)
			Expect(err).NotTo(HaveOccurred())
		})

		By("retrieving the user token using the new password", func() {
			var err error
			userToken, err = client.Users.GetToken(user.UserName, "new-password")
			Expect(err).NotTo(HaveOccurred())

			decodedToken, err := client.Tokens.Decode(userToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(decodedToken.UserID).To(Equal(user.ID))
		})
	})
})
