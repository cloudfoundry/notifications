package services_test

import (
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotificationsPreferences", func() {
	var builder services.PreferencesBuilder
	var TRUE = true
	var FALSE = false

	Describe("Add", func() {
		BeforeEach(func() {
			builder = services.NewPreferencesBuilder()
		})

		It("adds new preferences", func() {
			builder.Add(models.Preference{
				ClientID:          "clientID",
				KindID:            "kindID",
				Email:             true,
				KindDescription:   "kind description",
				SourceDescription: "client description",
			})

			node := builder.Clients["clientID"]["kindID"]
			Expect(node).To(Equal(services.Kind{
				Email:             &TRUE,
				KindDescription:   "kind description",
				SourceDescription: "client description",
			}))
		})

		It("adds new preferences to an old client", func() {
			builder.Add(models.Preference{
				ClientID:          "clientID",
				KindID:            "kindID",
				Email:             true,
				KindDescription:   "kind description",
				SourceDescription: "clientID description",
			})
			builder.Add(models.Preference{
				ClientID:          "clientID",
				KindID:            "new_kind",
				Email:             true,
				KindDescription:   "new kind description",
				SourceDescription: "clientID description",
			})

			node := builder.Clients["clientID"]["kindID"]
			Expect(node).To(Equal(services.Kind{
				Email:             &TRUE,
				KindDescription:   "kind description",
				SourceDescription: "clientID description",
			}))

			node = builder.Clients["clientID"]["new_kind"]
			Expect(node).To(Equal(services.Kind{
				Email:             &TRUE,
				KindDescription:   "new kind description",
				SourceDescription: "clientID description",
			}))
		})

		It("changes the value of an email", func() {
			builder.Add(models.Preference{
				ClientID: "clientID",
				KindID:   "kindID",
				Email:    true,
			})

			Expect(builder.Clients["clientID"]["kindID"].Email).To(Equal(&TRUE))

			builder.Add(models.Preference{
				ClientID: "clientID",
				KindID:   "kindID",
				Email:    false,
			})

			Expect(builder.Clients["clientID"]["kindID"].Email).To(Equal(&FALSE))
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

			Expect(builder.Clients["client1"]["kind1"].Email).To(Equal(&TRUE))
			Expect(builder.Clients["client1"]["kind2"].Email).To(Equal(&TRUE))
			Expect(builder.Clients["client2"]["kind1"].Email).To(Equal(&TRUE))
			Expect(builder.Clients["client2"]["kind2"].Email).To(Equal(&TRUE))
		})

		It("uses the fallback values for descriptions and counts, when there are none", func() {
			builder.Add(models.Preference{
				ClientID:          "raptors",
				KindID:            "hungry",
				Email:             true,
				KindDescription:   "",
				SourceDescription: "",
			})

			node := builder.Clients["raptors"]["hungry"]
			Expect(node).To(Equal(services.Kind{
				Email:             &TRUE,
				KindDescription:   "hungry",
				SourceDescription: "raptors",
			}))
		})
	})

	Describe("ToPreferences", func() {
		BeforeEach(func() {
			builder = services.NewPreferencesBuilder()
		})

		It("returns a slice of preferences from the populated map", func() {
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

				delete(badBuilder.Clients["electric-fence"], "zap")

				_, err := badBuilder.ToPreferences()

				Expect(err).ToNot(BeNil())

			})

			It("returns an error when the email data map is empty", func() {
				badBuilder.Add(models.Preference{
					ClientID: "TRex",
					KindID:   "glass-of-water",
					Email:    false,
				})

				kind := badBuilder.Clients["TRex"]["glass-of-water"]
				kind.Email = nil
				badBuilder.Clients["TRex"]["glass-of-water"] = kind

				_, err := badBuilder.ToPreferences()

				Expect(err).ToNot(BeNil())
			})
		})
	})
})
