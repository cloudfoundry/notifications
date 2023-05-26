package uaa_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account Exists", func() {
	var (
		warrantUserService   *mocks.WarrantUserService
		warrantClientService *mocks.WarrantClientService
		userFinder           uaa.UserFinder
	)

	BeforeEach(func() {
		warrantUserService = mocks.NewWarrantUserService()
		warrantClientService = mocks.NewWarrantClientService()
		warrantClientService.GetTokenCall.Returns.Token = "client-token"

		userFinder = uaa.NewUserFinder("client-id", "client-secret", warrantUserService, warrantClientService)
	})

	It("determines whether a user exists or not", func() {
		exists, err := userFinder.Exists("some-guid")
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeTrue())

		Expect(warrantClientService.GetTokenCall.Receives.ID).To(Equal("client-id"))
		Expect(warrantClientService.GetTokenCall.Receives.Secret).To(Equal("client-secret"))

		Expect(warrantUserService.GetCall.Receives.Token).To(Equal("client-token"))
	})

	Context("when the user does not exist", func() {
		It("returns false", func() {
			warrantUserService.GetCall.Returns.Error = warrant.NotFoundError{}

			exists, err := userFinder.Exists("some-guid")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})
	})

	Context("when getting the user causes an error", func() {
		It("returns an error", func() {
			warrantUserService.GetCall.Returns.Error = errors.New("UAA has gone away")

			_, err := userFinder.Exists("some-guid")
			Expect(err).To(MatchError(errors.New("UAA has gone away")))
		})
	})

	Context("when getting the token causes an error", func() {
		It("returns an error", func() {
			warrantClientService.GetTokenCall.Returns.Error = errors.New("no token")

			_, err := userFinder.Exists("some-guid")
			Expect(err).To(MatchError(errors.New("no token")))
		})
	})
})
