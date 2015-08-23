package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FindsUserGUIDs", func() {
	var finder services.FindsUserGUIDs
	var cc *mocks.CloudController
	var uaa *mocks.ZonedUAAClient

	BeforeEach(func() {
		cc = mocks.NewCloudController()
		uaa = mocks.NewZonedUAAClient()
		finder = services.NewFindsUserGUIDs(cc, uaa)
	})

	Context("when looking for GUIDs that have a scope", func() {
		BeforeEach(func() {
			uaa.UsersGUIDsByScopeResponse["this.scope"] = []string{"user-402", "user-525"}
		})

		It("returns the userGUIDs that have that scope", func() {
			guids, err := finder.UserGUIDsBelongingToScope("token", "this.scope")

			Expect(guids).To(Equal([]string{"user-402", "user-525"}))
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when uaa has an error", func() {
			It("returns the error", func() {
				uaa.UsersGUIDsByScopeError = errors.New("foobar")
				_, err := finder.UserGUIDsBelongingToScope("token", "this.scope")

				Expect(err).To(MatchError(errors.New("foobar")))
			})
		})
	})

	Context("when looking for GUIDs belonging to a space", func() {
		BeforeEach(func() {
			cc.UsersBySpaceGuid["space-001"] = []cf.CloudControllerUser{
				{GUID: "user-123"},
				{GUID: "user-789"},
			}
		})

		It("returns the user GUIDs for the space", func() {
			guids, err := finder.UserGUIDsBelongingToSpace("space-001", "token")

			Expect(guids).To(Equal([]string{"user-123", "user-789"}))
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when CloudController causes an error", func() {
			BeforeEach(func() {
				cc.GetUsersBySpaceGuidError = errors.New("BOOM!")
			})

			It("returns the error", func() {
				_, err := finder.UserGUIDsBelongingToSpace("space-001", "token")

				Expect(err).To(Equal(cc.GetUsersBySpaceGuidError))
			})
		})
	})

	Context("when looking for GUIDs belonging to an organization", func() {
		BeforeEach(func() {
			cc.UsersByOrganizationGuid["org-001"] = []cf.CloudControllerUser{
				{GUID: "user-456"},
				{GUID: "user-001"},
			}
		})

		It("returns the user GUIDs for the organization", func() {
			guids, err := finder.UserGUIDsBelongingToOrganization("org-001", "", "token")

			Expect(guids).To(Equal([]string{"user-456", "user-001"}))
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when CloudController causes an error", func() {
			BeforeEach(func() {
				cc.GetUsersByOrganizationGuidError = errors.New("BOOM!")
			})

			It("returns the error", func() {
				_, err := finder.UserGUIDsBelongingToOrganization("org-001", "", "token")

				Expect(err).To(Equal(cc.GetUsersByOrganizationGuidError))
			})
		})

		Context("when the role is OrgManager", func() {
			BeforeEach(func() {
				cc.ManagersByOrganization["org-001"] = []cf.CloudControllerUser{
					{GUID: "user-678"},
					{GUID: "user-xxx"},
				}
			})

			It("returns the organization managers for the organization", func() {
				guids, err := finder.UserGUIDsBelongingToOrganization("org-001", "OrgManager", "token")

				Expect(guids).To(Equal([]string{"user-678", "user-xxx"}))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the role is OrgAuditor", func() {
			BeforeEach(func() {
				cc.AuditorsByOrganization["org-001"] = []cf.CloudControllerUser{
					{GUID: "user-abc"},
					{GUID: "user-zzz"},
				}
			})

			It("returns the organization auditors for the organization", func() {
				guids, err := finder.UserGUIDsBelongingToOrganization("org-001", "OrgAuditor", "token")

				Expect(guids).To(Equal([]string{"user-abc", "user-zzz"}))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the role is BillingManager", func() {
			BeforeEach(func() {
				cc.BillingManagersByOrganization["org-001"] = []cf.CloudControllerUser{
					{GUID: "user-jkl"},
					{GUID: "user-aaa"},
				}
			})

			It("returns the billing managers for the organization", func() {
				guids, err := finder.UserGUIDsBelongingToOrganization("org-001", "BillingManager", "token")

				Expect(guids).To(Equal([]string{"user-jkl", "user-aaa"}))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
