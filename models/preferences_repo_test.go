package models_test

import (
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PreferencesRepo", func() {
	var repo models.PreferencesRepo
	var kinds models.KindsRepo
	var clients models.ClientsRepo
	var receipts models.ReceiptsRepo
	var conn *models.Connection
	var unsubscribeRepo models.UnsubscribesRepo

	BeforeEach(func() {
		TruncateTables()

		env := application.NewEnvironment()
		db := models.NewDatabase(sqlDB, models.Config{
			MigrationsPath: env.ModelMigrationsDir,
		})

		conn = db.Connection().(*models.Connection)

		kinds = models.NewKindsRepo()
		clients = models.NewClientsRepo()
		receipts = models.NewReceiptsRepo()
		unsubscribeRepo = models.NewUnsubscribesRepo()
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
				raptorClient := models.Client{
					ID:          "raptors",
					Description: "raptors description",
				}

				clients.Create(conn, raptorClient)

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

				kinds.Create(conn, nonCriticalKind)
				kinds.Create(conn, secondNonCriticalKind)
				kinds.Create(conn, nonCriticalKindThatUserHasNotReceived)
				kinds.Create(conn, criticalKind)
				kinds.Create(conn, otherUserKind)

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

				receipts.Create(conn, nonCriticalReceipt)
				receipts.Create(conn, secondNonCriticalReceipt)
				receipts.Create(conn, criticalReceipt)
				receipts.Create(conn, otherUserReceipt)
			})

			It("Returns a slice of non-critical notifications for this user", func() {
				err := unsubscribeRepo.Set(conn, "correct-user", "raptors", "sleepy", true)
				Expect(err).NotTo(HaveOccurred())

				results, err := repo.FindNonCriticalPreferences(conn, "correct-user")
				if err != nil {
					panic(err)
				}

				Expect(results).To(HaveLen(3))

				Expect(results).To(ContainElement(models.Preference{
					ClientID:          "raptors",
					KindID:            "sleepy",
					Email:             false,
					KindDescription:   "sleepy description",
					SourceDescription: "raptors description",
					Count:             402,
				}))

				Expect(results).To(ContainElement(models.Preference{
					ClientID:          "raptors",
					KindID:            "dead",
					Email:             true,
					KindDescription:   "dead description",
					SourceDescription: "raptors description",
					Count:             525,
				}))

				Expect(results).To(ContainElement(models.Preference{
					ClientID:          "raptors",
					KindID:            "orange",
					Email:             true,
					KindDescription:   "orange description",
					SourceDescription: "raptors description",
					Count:             0,
				}))

			})
		})
	})
})
