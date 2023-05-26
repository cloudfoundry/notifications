package cf_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/pivotal-cf-experimental/rainmaker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpaceFinder", func() {
	var (
		tokenGetter *mocks.WarrantClientService
		spaceGetter *mocks.RainmakerSpacesService
		finder      cf.SpaceFinder
	)

	BeforeEach(func() {
		tokenGetter = mocks.NewWarrantClientService()
		tokenGetter.GetTokenCall.Returns.Token = "some-token"
		spaceGetter = mocks.NewRainmakerSpacesService()
		spaceGetter.GetCall.Returns.Space = rainmaker.Space{
			GUID: "some-guid",
		}
		finder = cf.NewSpaceFinder("some-id", "some-secret", tokenGetter, spaceGetter)
	})

	It("finds a space given a guid", func() {
		exists, err := finder.Exists("some-guid")
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())

		Expect(spaceGetter.GetCall.Receives.GUID).To(Equal("some-guid"))
		Expect(spaceGetter.GetCall.Receives.Token).To(Equal("some-token"))

		Expect(tokenGetter.GetTokenCall.Receives.ID).To(Equal("some-id"))
		Expect(tokenGetter.GetTokenCall.Receives.Secret).To(Equal("some-secret"))
	})

	Context("when a space cannot be retrieved", func() {
		It("returns false", func() {
			spaceGetter.GetCall.Returns.Error = rainmaker.NotFoundError{}
			exists, err := finder.Exists("some-guid")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})
	})

	Context("when an error occurs", func() {
		Context("when a token cannot be retrieved", func() {
			It("returns an error", func() {
				tokenGetter.GetTokenCall.Returns.Error = errors.New("some error getting a token")
				_, err := finder.Exists("some-guid")
				Expect(err).To(MatchError(errors.New("some error getting a token")))
			})
		})

		Context("when rainmaker errors", func() {
			It("returns an error", func() {
				spaceGetter.GetCall.Returns.Error = errors.New("some error getting a space")
				_, err := finder.Exists("some-guid")
				Expect(err).To(MatchError(errors.New("some error getting a space")))
			})
		})
	})
})
