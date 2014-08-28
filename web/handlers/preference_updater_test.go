package handlers_test

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/web/handlers"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("PreferenceUpdater", func() {
    Describe("Execute", func() {

        var doorOpen models.Unsubscribe
        var barking models.Unsubscribe
        var repo *FakeUnsubscribesRepo
        var fakeDBConn *FakeDBConn
        var updater handlers.PreferenceUpdater

        BeforeEach(func() {
            fakeDBConn = &FakeDBConn{}
            repo = NewFakeUnsubscribesRepo()
            updater = handlers.NewPreferenceUpdater(repo)

            doorOpen = models.Unsubscribe{
                UserID:   "the-user",
                ClientID: "raptors",
                KindID:   "door-open",
            }

            barking = models.Unsubscribe{
                UserID:   "the-user",
                ClientID: "dogs",
                KindID:   "barking",
            }

        })

        It("Adds New Unsubscribes to the unsubscribes Repo", func() {

            updater.Execute(fakeDBConn, []models.Preference{
                {
                    ClientID: "raptors",
                    KindID:   "door-open",
                    Email:    false,
                },

                {
                    ClientID: "dogs",
                    KindID:   "barking",
                    Email:    false,
                },
            }, "the-user")

            Expect(len(repo.Unsubscribes)).To(Equal(2))
            Expect(repo.Unsubscribes).To(ContainElement(doorOpen))
            Expect(repo.Unsubscribes).To(ContainElement(barking))
        })

        It("does not insert duplicate unsubscribes", func() {
            _, err := repo.Create(fakeDBConn, models.Unsubscribe{
                UserID:   "my-user",
                ClientID: "raptors",
                KindID:   "door-open",
            })
            if err != nil {
                panic(err)
            }
            Expect(len(repo.Unsubscribes)).To(Equal(1))

            err = updater.Execute(fakeDBConn, []models.Preference{
                {
                    ClientID: "raptors",
                    KindID:   "door-open",
                    Email:    false,
                },
            }, "my-user")

            Expect(err).To(BeNil())
            Expect(len(repo.Unsubscribes)).To(Equal(1))
        })

        It("Does not add resubscriptions to the unsubscribes Repo", func() {
            updater.Execute(fakeDBConn, []models.Preference{
                {
                    ClientID: "dogs",
                    KindID:   "barking",
                    Email:    true,
                },
            }, "the-user")

            Expect(len(repo.Unsubscribes)).To(Equal(0))
        })
    })
})
