package services_test

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/services"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotificationsPreferences", func() {
    var builder services.PreferencesBuilder

    Describe("Add", func() {
        BeforeEach(func() {
            builder = services.NewPreferencesBuilder()
        })

        It("adds new preferences", func() {
            builder.Add(models.Preference{
                ClientID:          "client",
                KindID:            "kind",
                Email:             true,
                KindDescription:   "kind description",
                SourceDescription: "client description",
            })

            node := builder["client"]["kind"]
            Expect(node).To(Equal(map[string]interface{}{
                "email":              true,
                "kind_description":   "kind description",
                "source_description": "client description",
            }))
        })

        It("adds new preferences to an old client", func() {
            builder.Add(models.Preference{
                ClientID:          "client",
                KindID:            "kind",
                Email:             true,
                KindDescription:   "kind description",
                SourceDescription: "client description",
            })
            builder.Add(models.Preference{
                ClientID:          "client",
                KindID:            "new_kind",
                Email:             true,
                KindDescription:   "new kind description",
                SourceDescription: "client description",
            })

            node := builder["client"]["kind"]
            Expect(node).To(Equal(map[string]interface{}{
                "email":              true,
                "kind_description":   "kind description",
                "source_description": "client description",
            }))

            node = builder["client"]["new_kind"]
            Expect(node).To(Equal(map[string]interface{}{
                "email":              true,
                "kind_description":   "new kind description",
                "source_description": "client description",
            }))
        })

        It("changes the value of an email", func() {
            builder.Add(models.Preference{
                ClientID: "client",
                KindID:   "kind",
                Email:    true,
            })

            Expect(builder["client"]["kind"]["email"]).To(Equal(true))

            builder.Add(models.Preference{
                ClientID: "client",
                KindID:   "kind",
                Email:    false,
            })

            Expect(builder["client"]["kind"]["email"]).To(Equal(false))
        })

        It("can have multiple clients", func() {
            builder.Add(models.Preference{
                ClientID: "client1",
                KindID:   "kind1",
                Email:    true,
            })
            builder.Add(models.Preference{
                ClientID: "client1",
                KindID:   "kind2",
                Email:    true,
            })
            builder.Add(models.Preference{
                ClientID: "client2",
                KindID:   "kind1",
                Email:    true,
            })
            builder.Add(models.Preference{
                ClientID: "client2",
                KindID:   "kind2",
                Email:    true,
            })

            Expect(builder["client1"]["kind1"]["email"]).To(Equal(true))
            Expect(builder["client1"]["kind2"]["email"]).To(Equal(true))
            Expect(builder["client2"]["kind1"]["email"]).To(Equal(true))
            Expect(builder["client2"]["kind2"]["email"]).To(Equal(true))
        })

        It("uses the fallback values for descriptions, when there are none", func() {
            builder.Add(models.Preference{
                ClientID:          "client",
                KindID:            "kind",
                Email:             true,
                KindDescription:   "",
                SourceDescription: "",
            })

            node := builder["client"]["kind"]
            Expect(node).To(Equal(map[string]interface{}{
                "email":              true,
                "kind_description":   "kind",
                "source_description": "client",
            }))
        })
    })

    Describe("ToPreferences", func() {
        BeforeEach(func() {
            builder = services.NewPreferencesBuilder()
        })

        It("returns a slice of buildererences from the populated map", func() {
            builder.Add(models.Preference{
                ClientID: "raptors",
                KindID:   "door-open",
                Email:    true,
            })
            builder.Add(models.Preference{
                ClientID: "raptors",
                KindID:   "feeding-time",
                Email:    false,
            })
            builder.Add(models.Preference{
                ClientID: "dogs",
                KindID:   "barking",
                Email:    true,
            })

            preferences, err := builder.ToPreferences()
            if err != nil {
                panic(err)
            }

            Expect(len(preferences)).To(Equal(3))
            Expect(preferences).To(ContainElement(models.Preference{
                ClientID: "raptors",
                KindID:   "door-open",
                Email:    true,
            }))
            Expect(preferences).To(ContainElement(models.Preference{
                ClientID: "raptors",
                KindID:   "feeding-time",
                Email:    false,
            }))
            Expect(preferences).To(ContainElement(models.Preference{
                ClientID: "dogs",
                KindID:   "barking",
                Email:    true,
            }))
        })

        Context("invalid preferences", func() {

            var badBuilder services.PreferencesBuilder

            BeforeEach(func() {
                badBuilder = services.NewPreferencesBuilder()
            })

            It("returns an error when there are no kinds within a client", func() {

                badBuilder.Add(models.Preference{
                    ClientID: "electric-fence",
                    KindID:   "zap",
                    Email:    false,
                })

                delete(badBuilder["electric-fence"], "zap")

                _, err := badBuilder.ToPreferences()

                Expect(err).ToNot(BeNil())

            })

            It("returns an error when the email data map is empty", func() {
                badBuilder.Add(models.Preference{
                    ClientID: "TRex",
                    KindID:   "glass-of-water",
                    Email:    false,
                })

                delete(badBuilder["TRex"]["glass-of-water"], "email")

                _, err := badBuilder.ToPreferences()

                Expect(err).ToNot(BeNil())
            })

            It("returns an error when the email data map for emails cannot be coerced to a bool", func() {
                badBuilder.Add(models.Preference{
                    ClientID: "raptors",
                    KindID:   "feeding-time",
                    Email:    true,
                })

                badBuilder["raptors"]["feeding-time"]["email"] = "RUNNNNNNNNNNNNNNNN!"

                _, err := badBuilder.ToPreferences()

                Expect(err).ToNot(BeNil())
            })
        })
    })
})
