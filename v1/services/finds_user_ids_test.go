package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FindsUserIDs", func() {
	var (
		finder services.FindsUserIDs
		cc     *mocks.CloudController
		uaa    *mocks.ZonedUAAClient
	)

	BeforeEach(func() {
		cc = mocks.NewCloudController()
		uaa = mocks.NewZonedUAAClient()
		finder = services.NewFindsUserIDs(cc, uaa)
	})

	Context("UserIDsBelongingToScope", func() {
		BeforeEach(func() {
			uaa.UsersGUIDsByScopeCall.Returns.UserGUIDs = []string{"user-402", "user-525"}
		})

		It("returns the userIDs that have that scope", func() {
			guids, err := finder.UserIDsBelongingToScope("token", "this.scope")

			Expect(guids).To(Equal([]string{"user-402", "user-525"}))
			Expect(err).NotTo(HaveOccurred())

			Expect(uaa.UsersGUIDsByScopeCall.Receives.Token).To(Equal("token"))
			Expect(uaa.UsersGUIDsByScopeCall.Receives.Scope).To(Equal("this.scope"))
		})

		Context("when uaa has an error", func() {
			It("returns the error", func() {
				uaa.UsersGUIDsByScopeCall.Returns.Error = errors.New("foobar")

				_, err := finder.UserIDsBelongingToScope("token", "this.scope")
				Expect(err).To(MatchError(errors.New("foobar")))
			})
		})
	})

	Context("UserIDsBelongingToSpace", func() {
		BeforeEach(func() {
			cc.GetUsersBySpaceGuidCall.Returns.Users = []cf.CloudControllerUser{
				{GUID: "user-123"},
				{GUID: "user-789"},
			}
		})

		It("returns the user IDs for the space", func() {
			guids, err := finder.UserIDsBelongingToSpace("space-001", "token")
			Expect(err).NotTo(HaveOccurred())
			Expect(guids).To(Equal([]string{"user-123", "user-789"}))

			Expect(cc.GetUsersBySpaceGuidCall.Receives.SpaceGUID).To(Equal("space-001"))
			Expect(cc.GetUsersBySpaceGuidCall.Receives.Token).To(Equal("token"))
		})

		Context("when CloudController causes an error", func() {
			It("returns the error", func() {
				cc.GetUsersBySpaceGuidCall.Returns.Error = errors.New("BOOM!")

				_, err := finder.UserIDsBelongingToSpace("space-001", "token")
				Expect(err).To(MatchError(errors.New("BOOM!")))
			})
		})
	})

	Context("UserIDsBelongingToOrganization", func() {
		BeforeEach(func() {
			cc.GetUsersByOrgGuidCall.Returns.Users = []cf.CloudControllerUser{
				{GUID: "user-456"},
				{GUID: "user-001"},
			}
		})

		Context("when there is no role", func() {
			It("returns the user IDs for the organization", func() {
				guids, err := finder.UserIDsBelongingToOrganization("org-001", "", "token")
				Expect(err).NotTo(HaveOccurred())
				Expect(guids).To(Equal([]string{"user-456", "user-001"}))

				Expect(cc.GetUsersByOrgGuidCall.Receives.OrgGUID).To(Equal("org-001"))
				Expect(cc.GetUsersByOrgGuidCall.Receives.Token).To(Equal("token"))
			})

			Context("when CloudController causes an error", func() {
				It("returns the error", func() {
					cc.GetUsersByOrgGuidCall.Returns.Error = errors.New("BOOM!")
					_, err := finder.UserIDsBelongingToOrganization("org-001", "", "token")
					Expect(err).To(MatchError(errors.New("BOOM!")))
				})
			})
		})

		Context("when the role is OrgManager", func() {
			BeforeEach(func() {
				cc.GetManagersByOrgGuidCall.Returns.Users = []cf.CloudControllerUser{
					{GUID: "user-678"},
					{GUID: "user-xxx"},
				}
			})

			It("returns the organization managers for the organization", func() {
				guids, err := finder.UserIDsBelongingToOrganization("org-001", "OrgManager", "token")
				Expect(err).NotTo(HaveOccurred())
				Expect(guids).To(Equal([]string{"user-678", "user-xxx"}))

				Expect(cc.GetManagersByOrgGuidCall.Receives.OrgGUID).To(Equal("org-001"))
				Expect(cc.GetManagersByOrgGuidCall.Receives.Token).To(Equal("token"))
			})

			Context("when CloudController causes an error", func() {
				It("returns the error", func() {
					cc.GetManagersByOrgGuidCall.Returns.Error = errors.New("BOOM!")

					_, err := finder.UserIDsBelongingToOrganization("org-001", "OrgManager", "token")
					Expect(err).To(MatchError(errors.New("BOOM!")))
				})
			})
		})

		Context("when the role is OrgAuditor", func() {
			BeforeEach(func() {
				cc.GetAuditorsByOrgGuidCall.Returns.Users = []cf.CloudControllerUser{
					{GUID: "user-abc"},
					{GUID: "user-zzz"},
				}
			})

			It("returns the organization auditors for the organization", func() {
				guids, err := finder.UserIDsBelongingToOrganization("org-001", "OrgAuditor", "token")
				Expect(err).NotTo(HaveOccurred())
				Expect(guids).To(Equal([]string{"user-abc", "user-zzz"}))

				Expect(cc.GetAuditorsByOrgGuidCall.Receives.OrgGUID).To(Equal("org-001"))
				Expect(cc.GetAuditorsByOrgGuidCall.Receives.Token).To(Equal("token"))
			})

			Context("when CloudController causes an error", func() {
				It("returns the error", func() {
					cc.GetAuditorsByOrgGuidCall.Returns.Error = errors.New("BOOM!")

					_, err := finder.UserIDsBelongingToOrganization("org-001", "OrgAuditor", "token")
					Expect(err).To(MatchError(errors.New("BOOM!")))
				})
			})
		})

		Context("when the role is BillingManager", func() {
			BeforeEach(func() {
				cc.GetBillingManagersByOrgGuidCall.Returns.Users = []cf.CloudControllerUser{
					{GUID: "user-jkl"},
					{GUID: "user-aaa"},
				}
			})

			It("returns the billing managers for the organization", func() {
				guids, err := finder.UserIDsBelongingToOrganization("org-001", "BillingManager", "token")
				Expect(err).NotTo(HaveOccurred())
				Expect(guids).To(Equal([]string{"user-jkl", "user-aaa"}))

				Expect(cc.GetBillingManagersByOrgGuidCall.Receives.OrgGUID).To(Equal("org-001"))
				Expect(cc.GetBillingManagersByOrgGuidCall.Receives.Token).To(Equal("token"))
			})

			Context("when CloudController causes an error", func() {
				It("returns the error", func() {
					cc.GetBillingManagersByOrgGuidCall.Returns.Error = errors.New("BOOM!")

					_, err := finder.UserIDsBelongingToOrganization("org-001", "BillingManager", "token")
					Expect(err).To(MatchError(errors.New("BOOM!")))
				})
			})
		})
	})
})
