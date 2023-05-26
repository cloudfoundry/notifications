package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpaceLoader", func() {
	Describe("Load", func() {
		var (
			loader services.SpaceLoader
			cc     *mocks.CloudController
		)

		BeforeEach(func() {
			cc = mocks.NewCloudController()
			cc.LoadSpaceCall.Returns.Space = cf.CloudControllerSpace{
				GUID:             "space-001",
				Name:             "space-name",
				OrganizationGUID: "org-001",
			}

			loader = services.NewSpaceLoader(cc)
		})

		It("returns the space", func() {
			space, err := loader.Load("space-001", "some-token")
			Expect(err).NotTo(HaveOccurred())
			Expect(space).To(Equal(cf.CloudControllerSpace{
				GUID:             "space-001",
				Name:             "space-name",
				OrganizationGUID: "org-001",
			}))

			Expect(cc.LoadSpaceCall.Receives.SpaceGUID).To(Equal("space-001"))
			Expect(cc.LoadSpaceCall.Receives.Token).To(Equal("some-token"))
		})

		Context("when the space cannot be found", func() {
			It("returns an error object", func() {
				cc.LoadSpaceCall.Returns.Error = cf.NewFailure(404, "not found")

				_, err := loader.Load("missing-space", "some-token")
				Expect(err).To(MatchError(services.CCNotFoundError{Err: cf.NewFailure(404, "not found")}))
			})
		})

		Context("when Load returns any other type of error", func() {
			It("returns a CCDownError when the error is cf.Failure", func() {
				cc.LoadSpaceCall.Returns.Error = cf.NewFailure(401, "BOOM!")

				_, err := loader.Load("space-001", "some-token")
				Expect(err).To(MatchError(services.CCDownError{Err: cf.NewFailure(401, "BOOM!")}))
			})

			It("returns the same error for all other cases", func() {
				cc.LoadSpaceCall.Returns.Error = errors.New("BOOM!")

				_, err := loader.Load("space-001", "some-token")
				Expect(err).To(Equal(errors.New("BOOM!")))
			})
		})
	})
})
