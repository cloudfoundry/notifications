package horde_test

import (
	"bytes"
	"errors"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/horde"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("organizations audience", func() {
	var (
		userFinder    *mocks.FindsUserIDs
		orgFinder     *mocks.OrganizationLoader
		tokenLoader   *mocks.TokenLoader
		organizations horde.Organizations
		logger        lager.Logger
		logStream     *bytes.Buffer
	)

	BeforeEach(func() {
		userFinder = mocks.NewFindsUserIDs()
		userFinder.UserIDsBelongingToOrganizationCall.Returns.UserIDs = []string{"some-random-guid"}

		orgFinder = mocks.NewOrganizationLoader()
		orgFinder.LoadCall.Returns.Organizations = []cf.CloudControllerOrganization{
			{
				GUID: "some-silly-org-guid",
				Name: "SOME-SILLY",
			},
		}

		tokenLoader = mocks.NewTokenLoader()
		tokenLoader.LoadCall.Returns.Token = "token"

		organizations = horde.NewOrganizations(userFinder, orgFinder, tokenLoader, "https://uaa.example.com")

		logStream = bytes.NewBuffer([]byte{})
		logger = lager.NewLogger("notifications-whatever")
		logger.RegisterSink(lager.NewWriterSink(logStream, lager.DEBUG))
	})

	Describe("GenerateAudiences", func() {
		It("looks up userGUIDs and wraps them in User objects", func() {
			audiences, err := organizations.GenerateAudiences([]string{"some-silly-org-guid"}, logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(audiences).To(HaveLen(1))

			audience := audiences[0]
			Expect(audience.Users).To(Equal([]horde.User{{GUID: "some-random-guid"}}))
			Expect(audience.Endorsement).To(Equal("You received this message because you belong to the SOME-SILLY organization."))

			Expect(tokenLoader.LoadCall.Receives.UAAHost).To(Equal("https://uaa.example.com"))

			Expect(userFinder.UserIDsBelongingToOrganizationCall.Receives.OrgGUID).To(Equal("some-silly-org-guid"))
			Expect(userFinder.UserIDsBelongingToOrganizationCall.Receives.Role).To(Equal(""))
			Expect(userFinder.UserIDsBelongingToOrganizationCall.Receives.Token).To(Equal("token"))

			Expect(orgFinder.LoadCall.Receives.OrganizationGUID).To(Equal("some-silly-org-guid"))
			Expect(orgFinder.LoadCall.Receives.Token).To(Equal("token"))
		})

		Context("when we count 100 OrgGUIDs", func() {
			It("logs the count to the logger", func() {
				allOrgs := make([]string, 101)

				_, err := organizations.GenerateAudiences(allOrgs, logger)
				Expect(err).NotTo(HaveOccurred())

				message, err := logStream.ReadString('\n')
				Expect(err).NotTo(HaveOccurred())
				Expect(message).To(ContainSubstring(`{"processed":0}`))

				message, err = logStream.ReadString('\n')
				Expect(err).NotTo(HaveOccurred())
				Expect(message).To(ContainSubstring(`{"processed":100}`))
			})
		})

		Context("when a error occurs", func() {
			Context("when the token loader encounters an error", func() {
				It("returns the error", func() {
					tokenLoader.LoadCall.Returns.Error = errors.New("some token error")
					_, err := organizations.GenerateAudiences([]string{"some-silly-org-guid"}, logger)
					Expect(err).To(MatchError(errors.New("some token error")))
				})
			})

			Context("when the organizaton loader encounters an error", func() {
				Context("when the error is a NotFoundError", func() {
					BeforeEach(func() {
						orgFinder.LoadCall.Returns.Organizations = []cf.CloudControllerOrganization{
							{
								GUID: "some-silly-org-guid",
								Name: "SOME-SILLY",
							},
							{},
						}

						orgFinder.LoadCall.Returns.Errors = []error{
							nil,
							cf.NotFoundError{
								Message: "some org error",
							},
						}
					})

					It("returns the correct audience", func() {
						audiences, err := organizations.GenerateAudiences([]string{"some-silly-org-guid", "some-other-org-guid"}, logger)
						Expect(err).NotTo(HaveOccurred())
						Expect(audiences).To(ContainElement(horde.Audience{
							Users: []horde.User{
								{
									Email: "",
									GUID:  "some-random-guid",
								},
							},
							Endorsement: "You received this message because you belong to the SOME-SILLY organization.",
						}))
					})
				})

				Context("when any other error occurs", func() {
					It("returns the error", func() {
						orgFinder.LoadCall.Returns.Errors = []error{
							cf.Failure{
								Message: "some org finding error",
							},
						}

						_, err := organizations.GenerateAudiences([]string{"some-silly-org-guid"}, logger)
						Expect(err).To(MatchError(cf.Failure{Message: "some org finding error"}))
					})
				})
			})

			Context("when the user loader encounters an error", func() {
				It("returns the error", func() {
					userFinder.UserIDsBelongingToOrganizationCall.Returns.Error = errors.New("some user finding error")
					_, err := organizations.GenerateAudiences([]string{"some-silly-org-guid"}, logger)
					Expect(err).To(MatchError(errors.New("some user finding error")))
				})
			})
		})
	})
})
