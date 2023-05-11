package cf_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/pivotal-cf-experimental/rainmaker"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrgFinder", func() {
	var (
		tokenGetter *mocks.WarrantClientService
		orgGetter   *mocks.RainmakerOrganizationsService
		finder      cf.OrgFinder
	)

	BeforeEach(func() {
		tokenGetter = mocks.NewWarrantClientService()
		tokenGetter.GetTokenCall.Returns.Token = "some-token"
		orgGetter = mocks.NewRainmakerOrganizationsService()
		orgGetter.GetCall.Returns.Organization = rainmaker.Organization{
			GUID: "some-guid",
		}
		finder = cf.NewOrgFinder("some-id", "some-secret", tokenGetter, orgGetter)
	})

	It("finds an org given a guid", func() {
		exists, err := finder.Exists("some-guid")
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())

		Expect(orgGetter.GetCall.Receives.GUID).To(Equal("some-guid"))
		Expect(orgGetter.GetCall.Receives.Token).To(Equal("some-token"))

		Expect(tokenGetter.GetTokenCall.Receives.ID).To(Equal("some-id"))
		Expect(tokenGetter.GetTokenCall.Receives.Secret).To(Equal("some-secret"))
	})

	Context("when a org cannot be retrieved", func() {
		It("returns false", func() {
			orgGetter.GetCall.Returns.Error = rainmaker.NotFoundError{}

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
				orgGetter.GetCall.Returns.Error = errors.New("some error getting the org")

				_, err := finder.Exists("some-guid")
				Expect(err).To(MatchError(errors.New("some error getting the org")))
			})
		})
	})
})
