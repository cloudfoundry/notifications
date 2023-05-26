package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AllUserGUIDs", func() {
	var allUsers services.AllUsers
	var uaaClient *mocks.ZonedUAAClient
	var users []uaa.User

	BeforeEach(func() {
		uaaClient = mocks.NewZonedUAAClient()
		allUsers = services.NewAllUsers(uaaClient)
	})

	Context("when the request succeeds", func() {
		BeforeEach(func() {
			users = []uaa.User{
				{
					Emails: []string{"user-123@example.com"},
					ID:     "user-123",
				},
				{
					Emails: []string{"user-456@example.com"},
					ID:     "user-456",
				},
				{
					Emails: []string{"user-999@example.com"},
					ID:     "user-999",
				},
			}

			uaaClient.AllUsersCall.Returns.Users = users
		})

		It("returns the UAAUsers, UserGUIDs, and an error", func() {
			guids, err := allUsers.AllUserGUIDs("token")
			Expect(err).NotTo(HaveOccurred())
			Expect(guids).To(ConsistOf("user-456", "user-999", "user-123"))

			Expect(uaaClient.AllUsersCall.Receives.Token).To(Equal("token"))
		})
	})

	Context("when the request to UAA fails", func() {
		It("bubbles up the error", func() {
			uaaClient.AllUsersCall.Returns.Error = errors.New("BOOM!")

			_, err := allUsers.AllUserGUIDs("token")
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(errors.New("BOOM!")))
		})
	})
})
