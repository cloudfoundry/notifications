package models_test

import (
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("KindsRepo", func() {
	var (
		repo models.KindsRepo
		conn *db.Connection
	)

	BeforeEach(func() {

		repo = models.NewKindsRepo()
		database := db.NewDatabase(sqlDB, db.Config{})
		helpers.TruncateTables(database)
		conn = database.Connection().(*db.Connection)
	})

	Describe("Update", func() {
		Context("when the template id is meant to be set", func() {
			It("updates the record in the database", func() {
				kind := models.Kind{
					ID:         "my-kind",
					ClientID:   "my-client",
					TemplateID: "my-template",
				}

				kind, err := repo.Upsert(conn, kind)
				if err != nil {
					panic(err)
				}

				primary := kind.Primary
				createdAt := kind.CreatedAt

				kind.Description = "My Kind"
				kind.Critical = true
				kind.Primary = 42069
				kind.TemplateID = "new-template"
				kind.CreatedAt = time.Now().Add(-3 * time.Minute)

				kind, err = repo.Update(conn, kind)
				if err != nil {
					panic(err)
				}

				kind, err = repo.Find(conn, "my-kind", "my-client")
				if err != nil {
					panic(err)
				}

				Expect(kind.ID).To(Equal("my-kind"))
				Expect(kind.Description).To(Equal("My Kind"))
				Expect(kind.Critical).To(BeTrue())
				Expect(kind.ClientID).To(Equal("my-client"))
				Expect(kind.TemplateID).To(Equal("new-template"))
				Expect(kind.UpdatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
				Expect(kind.CreatedAt).To(Equal(createdAt))
				Expect(kind.Primary).To(Equal(primary))
			})

			It("returns a record not found error when the record does not exist", func() {
				kind := models.Kind{
					ID:         "my-kind",
					ClientID:   "my-client",
					TemplateID: "my-template",
				}
				_, err := repo.Update(conn, kind)
				Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("Notification with ID \"my-kind\" belonging to client \"my-client\" could not be found")}))
			})
		})

		Context("when the template id is not meant to be set", func() {
			It("updates the record in the database, using the existing template ID", func() {
				kind := models.Kind{
					ID:         "my-kind",
					ClientID:   "my-client",
					TemplateID: "my-template",
				}

				kind, err := repo.Upsert(conn, kind)
				if err != nil {
					panic(err)
				}

				primary := kind.Primary
				createdAt := kind.CreatedAt

				kind.Description = "My Kind"
				kind.Critical = true
				kind.Primary = 42069
				kind.TemplateID = models.DoNotSetTemplateID
				kind.CreatedAt = time.Now().Add(-3 * time.Minute)

				kind, err = repo.Update(conn, kind)
				if err != nil {
					panic(err)
				}

				kind, err = repo.Find(conn, "my-kind", "my-client")
				if err != nil {
					panic(err)
				}

				Expect(kind.ID).To(Equal("my-kind"))
				Expect(kind.Description).To(Equal("My Kind"))
				Expect(kind.Critical).To(BeTrue())
				Expect(kind.ClientID).To(Equal("my-client"))
				Expect(kind.TemplateID).To(Equal("my-template"))
				Expect(kind.UpdatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
				Expect(kind.CreatedAt).To(Equal(createdAt))
				Expect(kind.Primary).To(Equal(primary))
			})

			It("returns a record not found error when the record does not exist", func() {
				kind := models.Kind{
					ID:       "my-kind",
					ClientID: "my-client",
				}
				_, err := repo.Update(conn, kind)
				Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("Notification with ID \"my-kind\" belonging to client \"my-client\" could not be found")}))
			})
		})
	})

	Describe("Upsert", func() {
		Context("when the record is new", func() {
			It("inserts the record in the database", func() {
				kind := models.Kind{
					ID:          "my-kind",
					Description: "My Kind",
					Critical:    false,
					ClientID:    "my-client",
				}

				kind, err := repo.Upsert(conn, kind)
				if err != nil {
					panic(err)
				}

				kind, err = repo.Find(conn, "my-kind", "my-client")
				if err != nil {
					panic(err)
				}

				Expect(kind.ID).To(Equal("my-kind"))
				Expect(kind.Description).To(Equal("My Kind"))
				Expect(kind.Critical).To(BeFalse())
				Expect(kind.ClientID).To(Equal("my-client"))
				Expect(kind.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			})

			It("allows duplicate kindIDs that are unique by clientID", func() {
				kind1 := models.Kind{
					ID:       "forgotten-password",
					ClientID: "a-client",
				}

				kind2 := models.Kind{
					ID:       "forgotten-password",
					ClientID: "another-client",
				}

				kind1, err := repo.Upsert(conn, kind1)
				Expect(err).To(BeNil())

				kind2, err = repo.Upsert(conn, kind2)
				Expect(err).To(BeNil())

				firstKind, err := repo.Find(conn, "forgotten-password", "a-client")
				if err != nil {
					panic(err)
				}

				secondKind, err := repo.Find(conn, "forgotten-password", "another-client")
				if err != nil {
					panic(err)
				}

				Expect(firstKind).To(Equal(kind1))
				Expect(secondKind).To(Equal(kind2))
			})

			It("sets the template ID to 'default' when the field is empty", func() {
				kind := models.Kind{
					ID:          "my-kind",
					Description: "My Kind",
					Critical:    false,
					ClientID:    "my-client",
				}

				kind, err := repo.Upsert(conn, kind)
				if err != nil {
					panic(err)
				}

				kind, err = repo.Find(conn, "my-kind", "my-client")
				if err != nil {
					panic(err)
				}

				Expect(kind.ID).To(Equal("my-kind"))
				Expect(kind.Description).To(Equal("My Kind"))
				Expect(kind.Critical).To(BeFalse())
				Expect(kind.ClientID).To(Equal("my-client"))
				Expect(kind.TemplateID).To(Equal(models.DefaultTemplateID))
				Expect(kind.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
				Expect(kind.UpdatedAt).To(Equal(kind.CreatedAt))
			})
		})

		Context("when the record exists", func() {
			It("updates the record in the database", func() {
				kind := models.Kind{
					ID:       "my-kind",
					ClientID: "my-client",
				}

				kind, err := repo.Upsert(conn, kind)
				if err != nil {
					panic(err)
				}

				kind = models.Kind{
					ID:          "my-kind",
					Description: "My Kind",
					Critical:    true,
					ClientID:    "my-client",
				}

				kind, err = repo.Upsert(conn, kind)
				if err != nil {
					panic(err)
				}

				kind, err = repo.Find(conn, "my-kind", "my-client")
				if err != nil {
					panic(err)
				}

				Expect(kind.ID).To(Equal("my-kind"))
				Expect(kind.Description).To(Equal("My Kind"))
				Expect(kind.Critical).To(BeTrue())
				Expect(kind.ClientID).To(Equal("my-client"))
				Expect(kind.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			})
		})

		Context("when the record comes into existence before we create it", func() {
			It("updates the record in the database", func() {
				kind := models.Kind{
					ID:          "my-kind",
					Description: "My Kind",
					Critical:    true,
					ClientID:    "my-client",
				}

				conn := mocks.NewConnection()
				conn.InsertCall.Returns.Error = errors.New("Duplicate entry")

				_, err := repo.Upsert(conn, kind)
				Expect(err).NotTo(HaveOccurred())
				Expect(conn.UpdateCall.Receives.List).To(HaveLen(1))
				Expect(conn.UpdateCall.Receives.List[0].(*models.Kind).ID).To(Equal("my-kind"))
			})
		})
	})

	Describe("Trim", func() {
		It("deletes any kinds for the clientID that are not in the kindArray", func() {
			kind := models.Kind{
				ID:       "my-kind",
				ClientID: "the-client-id",
			}

			kindToDelete := models.Kind{
				ID:       "other-kind",
				ClientID: "the-client-id",
			}

			ignoredKind := models.Kind{
				ID:       "ignored-kind",
				ClientID: "other-client-id",
			}

			kind, err := repo.Upsert(conn, kind)
			if err != nil {
				panic(err)
			}

			kindToDelete, err = repo.Upsert(conn, kindToDelete)
			if err != nil {
				panic(err)
			}

			ignoredKind, err = repo.Upsert(conn, ignoredKind)
			if err != nil {
				panic(err)
			}

			count, err := repo.Trim(conn, "the-client-id", []string{"other-kind"})

			if err != nil {
				panic(err)
			}

			Expect(count).To(Equal(1))

			_, err = repo.Find(conn, "my-kind", "the-client-id")
			Expect(err).To(MatchError(models.NotFoundError{Err: errors.New("Notification with ID \"my-kind\" belonging to client \"the-client-id\" could not be found")}))

			_, err = repo.Find(conn, "ignored-kind", "other-client-id")
			if err != nil {
				panic(err)
			}

			_, err = repo.Find(conn, "other-kind", "the-client-id")
			if err != nil {
				panic(err)
			}
		})
	})

	Describe("FindAll", func() {
		It("returns all the records in the database", func() {
			kind1, err := repo.Upsert(conn, models.Kind{
				ID:       "my-kind",
				ClientID: "the-client-id",
			})
			if err != nil {
				panic(err)
			}

			kind2, err := repo.Upsert(conn, models.Kind{
				ID:       "another-kind",
				ClientID: "some-client-id",
			})
			if err != nil {
				panic(err)
			}

			kinds, err := repo.FindAll(conn)
			Expect(err).NotTo(HaveOccurred())

			Expect(kinds).To(HaveLen(2))
			Expect(kinds).To(ContainElement(kind1))
			Expect(kinds).To(ContainElement(kind2))
		})
	})

	Describe("FindAllByTemplateID", func() {
		It("returns all kinds with a given template ID", func() {
			kind, err := repo.Upsert(conn, models.Kind{
				ID:         "some-id",
				TemplateID: "some-template",
			})

			_, err = repo.Upsert(conn, models.Kind{
				ID: "another-id",
			})

			kinds, err := repo.FindAllByTemplateID(conn, "some-template")
			Expect(err).NotTo(HaveOccurred())
			Expect(kinds).To(HaveLen(1))
			Expect(kinds).To(ContainElement(kind))
		})
	})
})
