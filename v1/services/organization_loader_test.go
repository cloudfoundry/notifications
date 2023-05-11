package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrganizationLoader", func() {
	Describe("Load", func() {
		var (
			loader services.OrganizationLoader
			cc     *mocks.CloudController
		)

		BeforeEach(func() {
			cc = mocks.NewCloudController()

			cc.LoadOrganizationCall.Returns.Organization = cf.CloudControllerOrganization{
				GUID: "org-001",
				Name: "org-name",
			}

			loader = services.NewOrganizationLoader(cc)
		})

		It("returns the org", func() {
			org, err := loader.Load("org-001", "some-token")
			Expect(err).NotTo(HaveOccurred())
			Expect(org).To(Equal(cf.CloudControllerOrganization{
				GUID: "org-001",
				Name: "org-name",
			}))

			Expect(cc.LoadOrganizationCall.Receives.OrgGUID).To(Equal("org-001"))
			Expect(cc.LoadOrganizationCall.Receives.Token).To(Equal("some-token"))
		})

		Context("when the org cannot be found", func() {
			It("returns an error object", func() {
				cc.LoadOrganizationCall.Returns.Error = cf.NewFailure(404, "BOOM!")

				_, err := loader.Load("missing-org", "some-token")
				Expect(err).To(MatchError(services.CCNotFoundError{Err: cf.NewFailure(404, "BOOM!")}))
			})
		})

		Context("when Load returns any other type of error", func() {
			It("returns a CCDownError when the error is cf.Failure", func() {
				cc.LoadOrganizationCall.Returns.Error = cf.NewFailure(401, "BOOM!")

				_, err := loader.Load("org-001", "some-token")
				Expect(err).To(Equal(services.CCDownError{Err: cf.NewFailure(401, "BOOM!")}))
			})

			It("returns the same error for all other cases", func() {
				cc.LoadOrganizationCall.Returns.Error = errors.New("BOOM!")

				_, err := loader.Load("org-001", "some-token")
				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})
	})
})
