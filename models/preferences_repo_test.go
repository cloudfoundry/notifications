package models_test

import (
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/coopernurse/gorp"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("PreferencesRepo", func() {
    var repo models.PreferencesRepo
    var kinds models.KindsRepo
    var receipts models.ReceiptsRepo
    var conn *gorp.DbMap

    BeforeEach(func() {
        TruncateTables()
        conn = models.Database().Connection
        kinds = models.NewKindsRepo()
        receipts = models.NewReceiptsRepo()
        repo = models.NewPreferencesRepo()
    })

    Context("when there are no matching results in the database", func() {
        It("returns an an empty slice", func() {
            results, err := repo.FindNonCriticalPreferences(conn, "irrelevant-user")
            if err != nil {
                panic(err)
            }
            Expect(len(results)).To(Equal(0))
        })
    })

    Context("when there are matching results in the database", func() {
        Describe("FindNonCriticalPreferences", func() {
            BeforeEach(func() {
                nonCriticalKind := models.Kind{
                    ID:       "sleepy",
                    ClientID: "raptors",
                    Critical: false,
                }

                secondNonCriticalKind := models.Kind{
                    ID:       "dead",
                    ClientID: "raptors",
                    Critical: false,
                }

                nonCriticalKindThatUserHasNotReceived := models.Kind{
                    ID:       "orange",
                    ClientID: "raptors",
                    Critical: false,
                }

                criticalKind := models.Kind{
                    ID:       "hungry",
                    ClientID: "raptors",
                    Critical: true,
                }

                otherUserKind := models.Kind{
                    ID:       "fast",
                    ClientID: "raptors",
                    Critical: true,
                }

                kinds.Create(conn, nonCriticalKind)
                kinds.Create(conn, secondNonCriticalKind)
                kinds.Create(conn, nonCriticalKindThatUserHasNotReceived)
                kinds.Create(conn, criticalKind)
                kinds.Create(conn, otherUserKind)

                nonCriticalReceipt := models.Receipt{
                    ClientID: "raptors",
                    KindID:   "sleepy",
                    UserGUID: "correct-user",
                }

                secondNonCriticalReceipt := models.Receipt{
                    ClientID: "raptors",
                    KindID:   "dead",
                    UserGUID: "correct-user",
                }

                criticalReceipt := models.Receipt{
                    ClientID: "raptors",
                    KindID:   "hungry",
                    UserGUID: "correct-user",
                }

                otherUserReceipt := models.Receipt{
                    ClientID: "raptors",
                    KindID:   "fast",
                    UserGUID: "other-user",
                }

                receipts.Create(conn, nonCriticalReceipt)
                receipts.Create(conn, secondNonCriticalReceipt)
                receipts.Create(conn, criticalReceipt)
                receipts.Create(conn, otherUserReceipt)
            })

            It("Returns a slice of non-critical notifications for this user", func() {
                results, err := repo.FindNonCriticalPreferences(conn, "correct-user")
                if err != nil {
                    panic(err)
                }

                Expect(len(results)).To(Equal(3))

                Expect(results).To(ContainElement(models.Preference{
                    ClientID: "raptors",
                    KindID:   "sleepy",
                    Email:    true,
                }))

                Expect(results).To(ContainElement(models.Preference{
                    ClientID: "raptors",
                    KindID:   "dead",
                    Email:    true,
                }))

                Expect(results).To(ContainElement(models.Preference{
                    ClientID: "raptors",
                    KindID:   "orange",
                    Email:    true,
                }))

                Expect(results).To(ContainElement(models.Preference{
                    ClientID: "raptors",
                    KindID:   "orange",
                    Email:    "true",
                }))
            })

        })
    })
})
