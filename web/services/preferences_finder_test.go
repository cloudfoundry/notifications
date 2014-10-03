package services_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/test_helpers/fakes"
    "github.com/cloudfoundry-incubator/notifications/web/services"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("PreferencesFinder", func() {
    var finder *services.PreferencesFinder
    var fakePreferencesRepo *fakes.FakePreferencesRepo
    var preferences []models.Preference

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
        fakeGlobalUnsubscribesRepo.Set(&fakes.FakeDBConn{}, "correct-user", true)
        fakePreferencesRepo = fakes.NewFakePreferencesRepo(preferences)
        finder = services.NewPreferencesFinder(fakePreferencesRepo, fakeGlobalUnsubscribesRepo)
    })

    Describe("Find", func() {
        It("returns the set of notifications that are not critical", func() {
            expectedResult := services.NewPreferencesBuilder()
            expectedResult.Add(preferences[0])
            expectedResult.Add(preferences[1])
            expectedResult.GlobalUnsubscribe = true

            resultPreferences, err := finder.Find("correct-user")
            if err != nil {
                panic(err)
            }

            Expect(resultPreferences).To(Equal(expectedResult))
        })

        Context("when the preferences repo returns an error", func() {
            It("should propagate the error", func() {
                fakePreferencesRepo.FindError = errors.New("BOOM!")
                _, err := finder.Find("correct-user")

                Expect(err).To(Equal(fakePreferencesRepo.FindError))
            })
        })
    })
})
