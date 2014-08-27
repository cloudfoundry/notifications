package handlers_test

import (
    "errors"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Preferences", func() {
    var preference *handlers.Preference
    var fakePreferencesRepo *FakePreferencesRepo

    BeforeEach(func() {
        preferences := []models.Preference{models.Preference{
            ClientID: "raptors",
            KindID:   "non-critical-kind",
            Email:    true,
        }}

        preferences = append(preferences, models.Preference{ClientID: "raptors", KindID: "other-kind", Email: false})

        fakePreferencesRepo = NewFakePreferencesRepo(preferences)
        preference = handlers.NewPreference(fakePreferencesRepo)
    })

    Describe("Execute", func() {
        It("returns the set of notifications that are not critical", func() {

            result := handlers.NewNotificationPreferences()
            result.Add("raptors", "non-critical-kind", true)
            result.Add("raptors", "other-kind", true)

            preferences, err := preference.Execute("correct-user")
            if err != nil {
                panic(err)
            }

            Expect(preferences).To(Equal(result))
        })

        Context("when the preferences repo returns an error", func() {
            It("should propagate the error", func() {
                fakePreferencesRepo.FindError = errors.New("BOOM!")

                _, err := preference.Execute("correct-user")

                Expect(err).To(Equal(fakePreferencesRepo.FindError))
            })
        })
    })
})
