package services_test

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/test_helpers/fakes"
    "github.com/cloudfoundry-incubator/notifications/web/services"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("PreferenceUpdater", func() {
    Describe("Execute", func() {

        var doorOpen models.Unsubscribe
        var barking models.Unsubscribe
        var repo *fakes.FakeUnsubscribesRepo
        var kinds *fakes.FakeKindsRepo
        var fakeDBConn *fakes.FakeDBConn
        var updater services.PreferenceUpdater

        Context("When unsubscribing from existing kinds of existing clients", func() {
            BeforeEach(func() {
                fakeDBConn = &fakes.FakeDBConn{}
                repo = fakes.NewFakeUnsubscribesRepo()
                kinds = fakes.NewFakeKindsRepo()
                updater = services.NewPreferenceUpdater(repo, kinds)

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

                kinds.Create(fakeDBConn, models.Kind{
                    ID:       "door-open",
                    ClientID: "raptors",
                })

                kinds.Create(fakeDBConn, models.Kind{
                    ID:       "barking",
                    ClientID: "dogs",
                })

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

            It("removes unsubscribes when they are resubscribed", func() {
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
                        Email:    true,
                    },
                }, "my-user")

                Expect(err).To(BeNil())
                Expect(len(repo.Unsubscribes)).To(Equal(0))
            })
        })

        Context("when unsubscribing from missing client", func() {

            var hungry models.Preference
            var boo models.Preference

            BeforeEach(func() {
                fakeDBConn = &fakes.FakeDBConn{}
                repo = fakes.NewFakeUnsubscribesRepo()
                kinds = fakes.NewFakeKindsRepo()
                updater = services.NewPreferenceUpdater(repo, kinds)

                hungry = models.Preference{
                    ClientID: "raptors",
                    KindID:   "hungry",
                }

                boo = models.Preference{
                    ClientID: "ghosts",
                    KindID:   "boo",
                }

                kinds.Create(fakeDBConn, models.Kind{
                    ID:       "hungry",
                    ClientID: "raptors",
                })

                kinds.Create(fakeDBConn, models.Kind{
                    ID:       "boo",
                    ClientID: "missing-client",
                })

            })

            It("should return a MissingKindOrClientError", func() {
                err := updater.Execute(fakeDBConn, []models.Preference{hungry, boo}, "the-user")

                Expect(err).To(Equal(services.MissingKindOrClientError("The kind 'boo' cannot be found for client 'ghosts'")))
            })
        })

        Context("when unsubscribing from a missing kind", func() {

            var hungry models.Preference
            var dead models.Preference

            BeforeEach(func() {
                fakeDBConn = &fakes.FakeDBConn{}
                repo = fakes.NewFakeUnsubscribesRepo()
                kinds = fakes.NewFakeKindsRepo()
                updater = services.NewPreferenceUpdater(repo, kinds)

                hungry = models.Preference{
                    ClientID: "raptors",
                    KindID:   "hungry",
                }

                dead = models.Preference{
                    ClientID: "raptors",
                    KindID:   "dead",
                }

                kinds.Create(fakeDBConn, models.Kind{
                    ID:       "hungry",
                    ClientID: "raptors",
                })

            })

            It("should return a MissingKindOrClientError", func() {

                err := updater.Execute(fakeDBConn, []models.Preference{hungry, dead}, "the-user")

                Expect(err).To(Equal(services.MissingKindOrClientError("The kind 'dead' cannot be found for client 'raptors'")))
            })
        })

        Context("when unsubscribing from a critical kind", func() {

            var hungry models.Preference
            var barking models.Preference

            BeforeEach(func() {
                fakeDBConn = &fakes.FakeDBConn{}
                repo = fakes.NewFakeUnsubscribesRepo()
                kinds = fakes.NewFakeKindsRepo()
                updater = services.NewPreferenceUpdater(repo, kinds)

                hungry = models.Preference{
                    ClientID: "raptors",
                    KindID:   "hungry",
                }

                barking = models.Preference{
                    ClientID: "dogs",
                    KindID:   "barking",
                }

                kinds.Create(fakeDBConn, models.Kind{
                    ClientID: "raptors",
                    ID:       "hungry",
                    Critical: true,
                })

                kinds.Create(fakeDBConn, models.Kind{
                    ClientID: "dogs",
                    ID:       "barking",
                })

            })

            It("should return a CriticalKindError", func() {
                err := updater.Execute(fakeDBConn, []models.Preference{barking, hungry}, "the-user")

                Expect(err).To(Equal(services.CriticalKindError("The kind 'hungry' for the 'raptors' client is critical and cannot be unsubscribed from")))
            })
        })

    })

})
