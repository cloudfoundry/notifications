package handlers_test

import (
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("NotificationsPreferences", func() {

    var pref handlers.NotificationPreferences

    Describe("Add", func() {

        BeforeEach(func() {
            pref = handlers.NewNotificationPreferences()
        })
        It("Adding a new client, kind, and email", func() {
            pref.Add("client", "kind", true)

            Expect(pref["client"]["kind"]["email"]).To(Equal(true))

        })
        It("Adding a new kind to an old client", func() {
            pref.Add("client", "kind", true)
            pref.Add("client", "new_kind", true)

            Expect(pref["client"]["kind"]["email"]).To(Equal(true))
            Expect(pref["client"]["new_kind"]["email"]).To(Equal(true))
        })

        It("Changing the value of an email", func() {
            pref.Add("client", "kind", true)

            Expect(pref["client"]["kind"]["email"]).To(Equal(true))

            pref.Add("client", "kind", false)

            Expect(pref["client"]["kind"]["email"]).To(Equal(false))
        })

        It("Can have multiple clients", func() {
            Expect(pref["client"]["new_kind"]["email"]).To(Equal(false))
            pref.Add("client1", "kind1", true)
            pref.Add("client1", "kind2", true)
            pref.Add("client2", "kind1", true)
            pref.Add("client2", "kind2", true)

            Expect(pref["client1"]["kind1"]["email"]).To(Equal(true))
            Expect(pref["client1"]["kind2"]["email"]).To(Equal(true))
            Expect(pref["client2"]["kind1"]["email"]).To(Equal(true))
            Expect(pref["client2"]["kind2"]["email"]).To(Equal(true))
        })
    })
})
