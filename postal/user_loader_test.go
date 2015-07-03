package postal_test

import (
	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserLoader", func() {
	var loader postal.UserLoader
	var token string
	var uaaClient *fakes.ZonedUAAClient

	Describe("Load", func() {
		BeforeEach(func() {
			tokenHeader := map[string]interface{}{
				"alg": "FAST",
			}
			tokenClaims := map[string]interface{}{
				"client_id": "mister-client",
				"exp":       int64(3404281214),
				"scope":     []string{"notifications.write"},
			}
			token = fakes.BuildToken(tokenHeader, tokenClaims)

			uaaClient = fakes.NewZonedUAAClient()
			uaaClient.UsersByID = map[string]uaa.User{
				"user-123": {
					Emails: []string{"user-123@example.com"},
					ID:     "user-123",
				},
				"user-456": {
					Emails: []string{"user-456@example.com"},
					ID:     "user-456",
				},
				"user-999": {
					Emails: []string{"user-999@example.com"},
					ID:     "user-999",
				},
			}

			loader = postal.NewUserLoader(uaaClient)
		})

		Context("UAA returns a collection of users", func() {
			It("returns a map of users from GUID to uaa.User using a list of user GUIDs", func() {
				users, err := loader.Load([]string{"user-123", "user-789"}, token)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(users)).To(Equal(2))

				user123 := users["user-123"]
				Expect(user123.Emails[0]).To(Equal("user-123@example.com"))
				Expect(user123.ID).To(Equal("user-123"))

				user789, ok := users["user-789"]
				Expect(ok).To(BeTrue())
				Expect(user789).To(Equal(uaa.User{}))
			})
		})

		Describe("UAA Error Responses", func() {
			Context("when UAA cannot be reached", func() {
				It("returns a UAADownError", func() {
					uaaClient.ErrorForUserByID = uaa.NewFailure(404, []byte("Requested route ('uaa.10.244.0.34.xip.io') does not exist"))

					_, err := loader.Load([]string{"user-123"}, token)

					Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
				})
			})

			Context("when UAA returns an unknown UAA 404 error", func() {
				It("returns a UAAGenericError", func() {
					uaaClient.ErrorForUserByID = uaa.NewFailure(404, []byte("Weird message we haven't seen"))

					_, err := loader.Load([]string{"user-123"}, token)

					Expect(err).To(BeAssignableToTypeOf(postal.UAAGenericError("")))
				})
			})

			Context("when UAA returns an failure code that is not 404", func() {
				It("returns a UAADownError", func() {
					uaaClient.ErrorForUserByID = uaa.NewFailure(500, []byte("Doesn't matter"))

					_, err := loader.Load([]string{"user-123"}, token)

					Expect(err).To(BeAssignableToTypeOf(postal.UAADownError("")))
				})
			})
		})
	})
})
