package services_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/services"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("PreferencesFinder", func() {
    var finder *services.PreferencesFinder
    var fakePreferencesRepo *FakePreferencesRepo
    var preferences []models.Preference

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

        fakePreferencesRepo = NewFakePreferencesRepo(preferences)
        finder = services.NewPreferencesFinder(fakePreferencesRepo)
    })

    Describe("Find", func() {
        It("returns the set of notifications that are not critical", func() {
            result := services.NewPreferencesBuilder()
            result.Add(preferences[0])
            result.Add(preferences[1])

            preferences, err := finder.Find("correct-user")
            if err != nil {
                panic(err)
            }

            Expect(preferences).To(Equal(result))
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
