package common_test

import (
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/uaa"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserLoader", func() {
	var (
		loader    common.UserLoader
		token     string
		uaaClient *mocks.ZonedUAAClient
	)

	Describe("Load", func() {
		BeforeEach(func() {
			tokenHeader := map[string]interface{}{
				"alg": "RS256",
			}
			tokenClaims := map[string]interface{}{
				"client_id": "mister-client",
				"exp":       int64(3404281214),
				"scope":     []string{"notifications.write"},
			}
			token = helpers.BuildToken(tokenHeader, tokenClaims)

			uaaClient = mocks.NewZonedUAAClient()
			uaaClient.UsersEmailsByIDsCall.Returns.Users = []uaa.User{
				{
					Emails: []string{"user-123@example.com"},
					ID:     "user-123",
				},
			}

			loader = common.NewUserLoader(uaaClient)
		})

		Context("UAA returns a collection of users", func() {
			It("returns a map of users from GUID to uaa.User using a list of user GUIDs", func() {
				users, err := loader.Load([]string{"user-123", "user-789"}, token)

				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(HaveLen(2))

				user123 := users["user-123"]
				Expect(user123.Emails[0]).To(Equal("user-123@example.com"))
				Expect(user123.ID).To(Equal("user-123"))

				user789, ok := users["user-789"]
				Expect(ok).To(BeTrue())
				Expect(user789).To(Equal(uaa.User{}))

				Expect(uaaClient.UsersEmailsByIDsCall.Receives.Token).To(Equal(token))
				Expect(uaaClient.UsersEmailsByIDsCall.Receives.IDs).To(Equal([]string{"user-123", "user-789"}))
			})
		})

		Describe("UAA Error Responses", func() {
			Context("when UAA cannot be reached", func() {
				It("returns a UAADownError", func() {
					uaaClient.UsersEmailsByIDsCall.Returns.Error = uaa.NewFailure(404, []byte("Requested route ('uaa.10.244.0.34.xip.io') does not exist"))

					_, err := loader.Load([]string{"user-123"}, token)
					Expect(err).To(BeAssignableToTypeOf(common.UAADownError{}))
				})
			})

			Context("when UAA returns an unknown UAA 404 error", func() {
				It("returns a UAAGenericError", func() {
					uaaClient.UsersEmailsByIDsCall.Returns.Error = uaa.NewFailure(404, []byte("Weird message we haven't seen"))

					_, err := loader.Load([]string{"user-123"}, token)

					Expect(err).To(BeAssignableToTypeOf(common.UAAGenericError{}))
				})
			})

			Context("when UAA returns an failure code that is not 404", func() {
				It("returns a UAADownError", func() {
					uaaClient.UsersEmailsByIDsCall.Returns.Error = uaa.NewFailure(500, []byte("Doesn't matter"))

					_, err := loader.Load([]string{"user-123"}, token)

					Expect(err).To(BeAssignableToTypeOf(common.UAADownError{}))
				})
			})
		})
	})
})
