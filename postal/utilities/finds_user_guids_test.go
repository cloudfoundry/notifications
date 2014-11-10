package utilities_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/postal/utilities"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("FindsUserGUIDs", func() {
    var finder utilities.FindsUserGUIDs
    var cc *fakes.CloudController

    BeforeEach(func() {
        cc = fakes.NewCloudController()
        finder = utilities.NewFindsUserGUIDs(cc)
    })

    Context("when looking for GUIDs belonging to a space", func() {
        BeforeEach(func() {
            cc.UsersBySpaceGuid["space-001"] = []cf.CloudControllerUser{
                cf.CloudControllerUser{GUID: "user-123"},
                cf.CloudControllerUser{GUID: "user-789"},
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
                cf.CloudControllerUser{GUID: "user-456"},
                cf.CloudControllerUser{GUID: "user-001"},
            }
        })

        It("returns the user GUIDs for the organization", func() {
            guids, err := finder.UserGUIDsBelongingToOrganization("org-001", "token")

            Expect(guids).To(Equal([]string{"user-456", "user-001"}))
            Expect(err).NotTo(HaveOccurred())
        })

        Context("when CloudController causes an error", func() {
            BeforeEach(func() {
                cc.GetUsersByOrganizationGuidError = errors.New("BOOM!")
            })

            It("returns the error", func() {
                _, err := finder.UserGUIDsBelongingToOrganization("org-001", "token")

                Expect(err).To(Equal(cc.GetUsersByOrganizationGuidError))
            })
        })
    })
})
