package horde_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/horde"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("spaces audience", func() {
	var (
		userFinder  *mocks.FindsUserIDs
		orgFinder   *mocks.OrganizationLoader
		spaceFinder *mocks.SpaceLoader
		tokenLoader *mocks.TokenLoader
		spaces      horde.Spaces
	)

	BeforeEach(func() {
		userFinder = mocks.NewFindsUserIDs()
		userFinder.UserIDsBelongingToSpaceCall.Returns.UserIDs = []string{"some-random-guid"}

		orgFinder = mocks.NewOrganizationLoader()
		orgFinder.LoadCall.Returns.Organization = cf.CloudControllerOrganization{
			GUID: "some-silly-org-guid",
			Name: "SOME-SILLY",
		}

		spaceFinder = mocks.NewSpaceLoader()
		spaceFinder.LoadCall.Returns.Space = cf.CloudControllerSpace{
			OrganizationGUID: "some-silly-org-guid",
			GUID:             "some-silly-space",
			Name:             "SILLY-SPACE",
		}

		tokenLoader = mocks.NewTokenLoader()
		tokenLoader.LoadCall.Returns.Token = "token"

		spaces = horde.NewSpaces(userFinder, orgFinder, spaceFinder, tokenLoader, "https://uaa.example.com")
	})

	Describe("GenerateAudiences", func() {
		It("looks up userGUIDs and wraps them in User objects", func() {
			audiences, err := spaces.GenerateAudiences([]string{"some-silly-space"})
			Expect(err).NotTo(HaveOccurred())
			Expect(audiences).To(HaveLen(1))

			audience := audiences[0]
			Expect(audience.Users).To(Equal([]horde.User{{GUID: "some-random-guid"}}))
			Expect(audience.Endorsement).To(Equal(`You received this message because you belong to the "SILLY-SPACE" space in the "SOME-SILLY" organization.`))

			Expect(tokenLoader.LoadCall.Receives.UAAHost).To(Equal("https://uaa.example.com"))

			Expect(userFinder.UserIDsBelongingToSpaceCall.Receives.SpaceGUID).To(Equal("some-silly-space"))
			Expect(userFinder.UserIDsBelongingToSpaceCall.Receives.Token).To(Equal("token"))

			Expect(spaceFinder.LoadCall.Receives.SpaceGUID).To(Equal("some-silly-space"))
			Expect(spaceFinder.LoadCall.Receives.Token).To(Equal("token"))

			Expect(orgFinder.LoadCall.Receives.OrganizationGUID).To(Equal("some-silly-org-guid"))
			Expect(orgFinder.LoadCall.Receives.Token).To(Equal("token"))
		})

		Context("when a error occurs", func() {
			Context("when the token loader encounters an error", func() {
				It("returns the error", func() {
					tokenLoader.LoadCall.Returns.Error = errors.New("some token error")
					_, err := spaces.GenerateAudiences([]string{"some-silly-space"})
					Expect(err).To(MatchError(errors.New("some token error")))
				})
			})

			Context("when the organizaton loader encounters an error", func() {
				It("returns the error", func() {
					orgFinder.LoadCall.Returns.Error = errors.New("some org finding error")
					_, err := spaces.GenerateAudiences([]string{"some-silly-space"})
					Expect(err).To(MatchError(errors.New("some org finding error")))
				})
			})

			Context("when the space loader encounters an error", func() {
				It("returns the error", func() {
					spaceFinder.LoadCall.Returns.Error = errors.New("some space finding error")
					_, err := spaces.GenerateAudiences([]string{"some-silly-space"})
					Expect(err).To(MatchError(errors.New("some space finding error")))
				})
			})

			Context("when the user loader encounters an error", func() {
				It("returns the error", func() {
					userFinder.UserIDsBelongingToSpaceCall.Returns.Error = errors.New("some user finding error")
					_, err := spaces.GenerateAudiences([]string{"some-silly-space"})
					Expect(err).To(MatchError(errors.New("some user finding error")))
				})
			})
		})
	})
})
