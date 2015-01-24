package rainmaker_test

import (
	"github.com/pivotal-golang/rainmaker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UsersService", func() {
	var token string
	var service *rainmaker.UsersService

	BeforeEach(func() {
		token = "token"
		service = rainmaker.NewUsersService(rainmaker.Config{
			Host: fakeCloudController.URL(),
		})
	})

	Describe("Create/Get", func() {
		It("creates a new user record in cloud controller and allows it to be fetched", func() {
			userGUID := "new-user-guid"

			user, err := service.Create(userGUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(user.GUID).To(Equal(userGUID))

			fetchedUser, err := service.Get(userGUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedUser.GUID).To(Equal(userGUID))
		})

		Context("when the request errors", func() {
			BeforeEach(func() {
				service = rainmaker.NewUsersService(rainmaker.Config{})
			})

			It("returns the error", func() {
				_, err := service.Create("user-guid", token)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
