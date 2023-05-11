package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PreferencesRepo", func() {
	var (
		repo            models.PreferencesRepo
		kinds           models.KindsRepo
		clients         models.ClientsRepo
		conn            *db.Connection
		unsubscribeRepo models.UnsubscribesRepo
	)

	BeforeEach(func() {
		database := db.NewDatabase(sqlDB, db.Config{})
		helpers.TruncateTables(database)

		conn = database.Connection().(*db.Connection)

		kinds = models.NewKindsRepo()
		clients = models.NewClientsRepo()
		unsubscribeRepo = models.NewUnsubscribesRepo()
		repo = models.NewPreferencesRepo()
	})

	Context("when there are no matching results in the database", func() {
		It("returns an an empty slice", func() {
			results, err := repo.FindNonCriticalPreferences(conn, "irrelevant-user")
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(0))
		})
	})

	Context("when there are matching results in the database", func() {
		Describe("FindNonCriticalPreferences", func() {
			BeforeEach(func() {
				raptorClient := models.Client{
					ID:          "raptors",
					Description: "raptors description",
				}

				_, err := clients.Upsert(conn, raptorClient)
				Expect(err).NotTo(HaveOccurred())

				nonCriticalKind := models.Kind{
					ID:          "sleepy",
					Description: "sleepy description",
					ClientID:    "raptors",
					Critical:    false,
				}

				secondNonCriticalKind := models.Kind{
					ID:          "dead",
					Description: "dead description",
					ClientID:    "raptors",
					Critical:    false,
				}

				nonCriticalKindThatUserHasNotReceived := models.Kind{
					ID:          "orange",
					Description: "orange description",
					ClientID:    "raptors",
					Critical:    false,
				}

				criticalKind := models.Kind{
					ID:          "hungry",
					Description: "hungry description",
					ClientID:    "raptors",
					Critical:    true,
				}

				otherUserKind := models.Kind{
					ID:          "fast",
					Description: "fast description",
					ClientID:    "raptors",
					Critical:    true,
				}

				kinds.Upsert(conn, nonCriticalKind)
				kinds.Upsert(conn, secondNonCriticalKind)
				kinds.Upsert(conn, nonCriticalKindThatUserHasNotReceived)
				kinds.Upsert(conn, criticalKind)
				kinds.Upsert(conn, otherUserKind)

				nonCriticalReceipt := models.Receipt{
					ClientID: "raptors",
					KindID:   "sleepy",
					UserGUID: "correct-user",
					Count:    402,
				}

				secondNonCriticalReceipt := models.Receipt{
					ClientID: "raptors",
					KindID:   "dead",
					UserGUID: "correct-user",
					Count:    525,
				}

				criticalReceipt := models.Receipt{
					ClientID: "raptors",
					KindID:   "hungry",
					UserGUID: "correct-user",
					Count:    89,
				}

				otherUserReceipt := models.Receipt{
					ClientID: "raptors",
					KindID:   "fast",
					UserGUID: "other-user",
					Count:    83,
				}

				createReceipt(conn, nonCriticalReceipt)
				createReceipt(conn, secondNonCriticalReceipt)
				createReceipt(conn, criticalReceipt)
				createReceipt(conn, otherUserReceipt)
			})

			It("returns a slice of non-critical notifications for this user", func() {
				err := unsubscribeRepo.Set(conn, "correct-user", "raptors", "sleepy", true)
				Expect(err).NotTo(HaveOccurred())

				results, err := repo.FindNonCriticalPreferences(conn, "correct-user")
				Expect(err).NotTo(HaveOccurred())

				Expect(results).To(HaveLen(3))

				Expect(results).To(ContainElement(models.Preference{
					ClientID:          "raptors",
					KindID:            "sleepy",
					Email:             false,
					KindDescription:   "sleepy description",
					SourceDescription: "raptors description",
				}))

				Expect(results).To(ContainElement(models.Preference{
					ClientID:          "raptors",
					KindID:            "dead",
					Email:             true,
					KindDescription:   "dead description",
					SourceDescription: "raptors description",
				}))

				Expect(results).To(ContainElement(models.Preference{
					ClientID:          "raptors",
					KindID:            "orange",
					Email:             true,
					KindDescription:   "orange description",
					SourceDescription: "raptors description",
				}))
			})
		})
	})
})
