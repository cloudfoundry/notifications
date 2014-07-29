package postal_test

import (
    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/postal"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("SpaceLoader", func() {
    Describe("Load", func() {
        var loader postal.SpaceLoader
        var token string
        var fakeCC *FakeCloudController

        BeforeEach(func() {
            fakeCC = NewFakeCloudController()
            fakeCC.Spaces = map[string]cf.CloudControllerSpace{
                "space-001": cf.CloudControllerSpace{
                    Guid:             "space-001",
                    Name:             "space-name",
                    OrganizationGuid: "org-001",
                },
            }
            fakeCC.Orgs = map[string]cf.CloudControllerOrganization{
                "org-001": cf.CloudControllerOrganization{
                    Guid: "org-001",
                    Name: "org-name",
                },
            }
            loader = postal.NewSpaceLoader(fakeCC)
        })

        Context("when GUID represents a space", func() {
            It("returns the name of the space and org", func() {
                space, org, err := loader.Load("space-001", token, postal.IsSpace)
                if err != nil {
                    panic(err)
                }

                Expect(space).To(Equal("space-name"))
                Expect(org).To(Equal("org-name"))
            })

            Context("when the space cannot be found", func() {
                It("returns an error object", func() {
                    _, _, err := loader.Load("space-doesnotexist", token, postal.IsSpace)

                    Expect(err.Error()).To(Equal("CloudController Failure (404): NOT FOUND"))
                })
            })
        })

        Context("when GUID represents a user", func() {
            It("returns empty values for space, org, and error", func() {
                space, org, err := loader.Load("user-001", token, postal.IsUser)
                if err != nil {
                    panic(err)
                }

                Expect(space).To(Equal(""))
                Expect(org).To(Equal(""))
                Expect(err).To(BeNil())
            })
        })
    })
})
