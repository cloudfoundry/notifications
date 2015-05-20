package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PreferencesFinder", func() {
	var (
		finder          *services.PreferencesFinder
		preferencesRepo *fakes.PreferencesRepo
		preferences     []models.Preference
		database        *fakes.Database
	)

	BeforeEach(func() {
		preferences = []models.Preference{
			{
				ClientID:          "raptors",
				SourceDescription: "raptors description",
				KindID:            "non-critical-kind",
				KindDescription:   "non critical kind description",
				Email:             true,
				Count:             3,
			},
			{
				ClientID:          "raptors",
				SourceDescription: "raptors description",
				KindID:            "other-kind",
				KindDescription:   "other kind description",
				Email:             false,
				Count:             10,
			},
		}

		fakeGlobalUnsubscribesRepo := fakes.NewGlobalUnsubscribesRepo()
		fakeGlobalUnsubscribesRepo.Set(fakes.NewConnection(), "correct-user", true)
		preferencesRepo = fakes.NewPreferencesRepo(preferences)
		database = fakes.NewDatabase()

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

			Expect(database.ConnectionWasCalled).To(BeTrue())
		})

		Context("when the preferences repo returns an error", func() {
			It("should propagate the error", func() {
				preferencesRepo.FindError = errors.New("BOOM!")

				_, err := finder.Find(database, "correct-user")
				Expect(err).To(Equal(preferencesRepo.FindError))
			})
		})
	})
})
