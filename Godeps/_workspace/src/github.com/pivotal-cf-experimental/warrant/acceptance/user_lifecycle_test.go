package acceptance

import (
	"time"

	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Lifecycle", func() {
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
		client.Users.Delete(user.ID, UAAToken)

		_, err := client.Users.Get(user.ID, UAAToken)
		Expect(err).To(BeAssignableToTypeOf(warrant.NotFoundError{}))
	})

	It("creates, retrieves, and deletes a user", func() {
		By("creating a new user", func() {
			var err error
			user, err = client.Users.Create(UAADefaultUsername, "warrant-user@example.com", UAAToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(user.UserName).To(Equal(UAADefaultUsername))
			Expect(user.Emails).To(ConsistOf([]string{"warrant-user@example.com"}))
			Expect(user.CreatedAt).To(BeTemporally("~", time.Now().UTC(), 10*time.Minute))
			Expect(user.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), 10*time.Minute))
			Expect(user.Version).To(Equal(0))
			Expect(user.Active).To(BeTrue())
			Expect(user.Verified).To(BeFalse())
			Expect(user.Origin).To(Equal("uaa"))
			//Expect(user.Groups).To(ConsistOf([]warrant.Group{})) TODO: finish up groups implementation
		})

		By("finding the user", func() {
			fetchedUser, err := client.Users.Get(user.ID, UAAToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedUser).To(Equal(user))
		})

		By("updating the user", func() {
			updatedUser, err := client.Users.Update(user, UAAToken)
			Expect(err).NotTo(HaveOccurred())

			fetchedUser, err := client.Users.Get(user.ID, UAAToken)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedUser).To(Equal(updatedUser))
		})
	})

	It("does not allow a user to be created without an email address", func() {
		_, err := client.Users.Create("warrant-invalid-email-user", "", UAAToken)
		Expect(err).To(BeAssignableToTypeOf(warrant.BadRequestError{}))
	})

	It("does not allow non-existant users to be updated", func() {
		user, err := client.Users.Create(UAADefaultUsername, "warrant-user@example.com", UAAToken)
		Expect(err).NotTo(HaveOccurred())

		originalUserID := user.ID
		user.ID = "non-existant-user-guid"
		_, err = client.Users.Update(user, UAAToken)
		Expect(err).To(BeAssignableToTypeOf(warrant.NotFoundError{}))

		err = client.Users.Delete(originalUserID, UAAToken)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when a user already exists", func() {
		BeforeEach(func() {
			var err error
			user, err = client.Users.Create(UAADefaultUsername, "warrant-user@example.com", UAAToken)
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not allow duplicate users to be created", func() {
			_, err := client.Users.Create(UAADefaultUsername, "warrant-user@example.com", UAAToken)
			Expect(err).To(BeAssignableToTypeOf(warrant.DuplicateResourceError{}))
		})
	})
})
