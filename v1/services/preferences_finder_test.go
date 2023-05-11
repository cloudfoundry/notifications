package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PreferencesFinder", func() {
	var (
		finder          *services.PreferencesFinder
		preferencesRepo *mocks.PreferencesRepo
		preferences     []models.Preference
		database        *mocks.Database
		conn            *mocks.Connection
	)

	BeforeEach(func() {
		preferences = []models.Preference{
			{
				ClientID:          "raptors",
				SourceDescription: "raptors description",
				KindID:            "non-critical-kind",
				KindDescription:   "non critical kind description",
				Email:             true,
			},
			{
				ClientID:          "raptors",
				SourceDescription: "raptors description",
				KindID:            "other-kind",
				KindDescription:   "other kind description",
				Email:             false,
			},
		}

		fakeGlobalUnsubscribesRepo := mocks.NewGlobalUnsubscribesRepo()
		fakeGlobalUnsubscribesRepo.GetCall.Returns.Unsubscribed = true

		preferencesRepo = mocks.NewPreferencesRepo()
		preferencesRepo.FindNonCriticalPreferencesCall.Returns.Preferences = preferences

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		finder = services.NewPreferencesFinder(preferencesRepo, fakeGlobalUnsubscribesRepo)
	})

	Describe("Find", func() {
		It("returns the set of notifications that are not critical", func() {
			expectedResult := services.NewPreferencesBuilder()
			expectedResult.Add(preferences[0])
			expectedResult.Add(preferences[1])
			expectedResult.GlobalUnsubscribe = true

			resultPreferences, err := finder.Find(database, "correct-user")
			Expect(err).NotTo(HaveOccurred())
			Expect(resultPreferences).To(Equal(expectedResult))

			Expect(preferencesRepo.FindNonCriticalPreferencesCall.Receives.Connection).To(Equal(conn))
			Expect(preferencesRepo.FindNonCriticalPreferencesCall.Receives.UserGUID).To(Equal("correct-user"))
		})

		Context("when the preferences repo returns an error", func() {
			It("should propagate the error", func() {
				preferencesRepo.FindNonCriticalPreferencesCall.Returns.Error = errors.New("BOOM!")

				_, err := finder.Find(database, "correct-user")
				Expect(err).To(Equal(preferencesRepo.FindNonCriticalPreferencesCall.Returns.Error))
			})
		})
	})
})
