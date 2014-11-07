package utilities_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/fakes"
    "github.com/cloudfoundry-incubator/notifications/postal/utilities"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("OrganizationLoader", func() {
    Describe("Load", func() {
        var loader utilities.OrganizationLoader
        var token string
        var cc *fakes.CloudController

        BeforeEach(func() {
            cc = fakes.NewCloudController()
            cc.Orgs = map[string]cf.CloudControllerOrganization{
                "org-001": cf.CloudControllerOrganization{
                    GUID: "org-001",
                    Name: "org-name",
                }, "org-123": cf.CloudControllerOrganization{
                    GUID: "org-123",
                    Name: "org-piggies",
                },
            }
            loader = utilities.NewOrganizationLoader(cc)
        })

        It("returns the org", func() {
            org, err := loader.Load("org-001", token)
            if err != nil {
                panic(err)
            }

            Expect(org).To(Equal(cf.CloudControllerOrganization{
                GUID: "org-001",
                Name: "org-name",
            }))
        })

        Context("when the org cannot be found", func() {
            It("returns an error object", func() {
                _, err := loader.Load("org-doesnotexist", token)

                Expect(err).To(BeAssignableToTypeOf(utilities.CCNotFoundError("")))
                Expect(err.Error()).To(Equal(`CloudController Error: CloudController Failure (404): {"code":30003,"description":"The organization could not be found: org-doesnotexist","error_code":"CF-OrganizationNotFound"}`))
            })
        })

        Context("when Load returns any other type of error", func() {
            It("returns a CCDownError when the error is cf.Failure", func() {
                failure := cf.NewFailure(401, "BOOM!")
                cc.LoadOrganizationError = failure
                _, err := loader.Load("org-001", token)

                Expect(err).To(Equal(utilities.CCDownError(failure.Error())))
            })

            It("returns the same error for all other cases", func() {
                cc.LoadOrganizationError = errors.New("BOOM!")
                _, err := loader.Load("org-001", token)

                Expect(err).To(Equal(errors.New("BOOM!")))
            })
        })
    })
})
