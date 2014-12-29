package models_test

import (
	"path"
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UnsubscribesRepo", func() {
	var repo models.UnsubscribesRepo
	var conn *models.Connection

	BeforeEach(func() {
		TruncateTables()
		repo = models.NewUnsubscribesRepo()

		env := application.NewEnvironment()
		migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
		db := models.NewDatabase(models.Config{
			DatabaseURL:    env.DatabaseURL,
			MigrationsPath: migrationsPath,
		})
		conn = db.Connection().(*models.Connection)
	})

	Describe("Create/Find", func() {
		It("stores the unsubscribe record into the database", func() {
			unsubscribe := models.Unsubscribe{
				ClientID: "raptors",
				KindID:   "hungry-kind",
				UserID:   "correct-user",
			}

			unsubscribe, err := repo.Create(conn, unsubscribe)
			if err != nil {
				panic(err)
			}

			unsubscribe, err = repo.Find(conn, "raptors", "hungry-kind", "correct-user")
			if err != nil {
				panic(err)
			}

			Expect(unsubscribe.ClientID).To(Equal("raptors"))
			Expect(unsubscribe.KindID).To(Equal("hungry-kind"))
			Expect(unsubscribe.UserID).To(Equal("correct-user"))
			Expect(unsubscribe.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
		})

		It("returns an error duplicate record if the record is already in the database", func() {
			unsubscribe := models.Unsubscribe{
				ClientID: "raptors",
				KindID:   "hungry-kind",
				UserID:   "correct-user",
			}

			_, err := repo.Create(conn, unsubscribe)
			if err != nil {
				panic(err)
			}

			_, err = repo.Create(conn, unsubscribe)
			Expect(err).To(BeAssignableToTypeOf(models.DuplicateRecordError{}))
		})

		It("returns a record not found error when the record does not exist", func() {
			_, err := repo.Find(conn, "bad-client", "bad-kind", "bad-user")
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
		})
	})

	Describe("Upsert", func() {
		It("inserts new records into the database", func() {
			unsubscribe := models.Unsubscribe{
				ClientID: "raptors",
				KindID:   "hungry-kind",
				UserID:   "correct-user",
			}

			unsubscribe, err := repo.Upsert(conn, unsubscribe)
			if err != nil {
				panic(err)
			}

			results, err := conn.Select(models.Unsubscribe{}, "SELECT * FROM `unsubscribes`")
			if err != nil {
				panic(err)
			}

			Expect(len(results)).To(Equal(1))
			unsub := *(results[0].(*models.Unsubscribe))

			Expect(unsub).To(Equal(unsubscribe))
		})

		It("updates existing records in the database", func() {
			unsubscribe := models.Unsubscribe{
				ClientID: "raptors",
				KindID:   "hungry-kind",
				UserID:   "correct-user",
			}

			unsubscribe, err := repo.Create(conn, unsubscribe)
			if err != nil {
				panic(err)
			}

			unsubscribe, err = repo.Upsert(conn, unsubscribe)
			if err != nil {
				panic(err)
			}

			results, err := conn.Select(models.Unsubscribe{}, "SELECT * FROM `unsubscribes`")
			if err != nil {
				panic(err)
			}

			Expect(len(results)).To(Equal(1))
			unsub := *(results[0].(*models.Unsubscribe))

			Expect(unsub).To(Equal(unsubscribe))
		})
	})

	Describe("FindAllByUserID", func() {
		It("finds all unsubscribes for a user", func() {
			unsub1, err := repo.Create(conn, models.Unsubscribe{
				UserID:   "correct-user",
				ClientID: "raptors",
				KindID:   "hungry",
			})
			if err != nil {
				panic(err)
			}

			unsub2, err := repo.Create(conn, models.Unsubscribe{
				UserID:   "correct-user",
				ClientID: "raptors",
				KindID:   "sleepy",
			})
			if err != nil {
				panic(err)
			}

			_, err = repo.Create(conn, models.Unsubscribe{
				UserID:   "other-user",
				ClientID: "dogs",
				KindID:   "barking",
			})
			if err != nil {
				panic(err)
			}

			unsubscribes, err := repo.FindAllByUserID(conn, "correct-user")
			if err != nil {
				panic(err)
			}

			Expect(len(unsubscribes)).To(Equal(2))
			Expect(unsubscribes).To(ContainElement(unsub1))
			Expect(unsubscribes).To(ContainElement(unsub2))
		})
	})

	Describe("Destroy", func() {
		It("removes the record from the database", func() {
			unsub1, err := repo.Create(conn, models.Unsubscribe{
				UserID:   "correct-user",
				ClientID: "raptors",
				KindID:   "hungry",
			})
			if err != nil {
				panic(err)
			}

			unsub2, err := repo.Create(conn, models.Unsubscribe{
				UserID:   "correct-user",
				ClientID: "raptors",
				KindID:   "sleepy",
			})
			if err != nil {
				panic(err)
			}

			unsubscribes, err := repo.FindAllByUserID(conn, "correct-user")
			if err != nil {
				panic(err)
			}

			Expect(len(unsubscribes)).To(Equal(2))
			Expect(unsubscribes).To(ContainElement(unsub1))
			Expect(unsubscribes).To(ContainElement(unsub2))

			_, err = repo.Destroy(conn, unsub1)
			if err != nil {
				panic(err)
			}

			unsubscribes, err = repo.FindAllByUserID(conn, "correct-user")
			if err != nil {
				panic(err)
			}

			Expect(len(unsubscribes)).To(Equal(1))
			Expect(unsubscribes).ToNot(ContainElement(unsub1))
			Expect(unsubscribes).To(ContainElement(unsub2))
		})
	})
})
