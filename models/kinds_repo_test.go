package models_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KindsRepo", func() {
	var repo models.KindsRepo
	var conn *models.Connection

	BeforeEach(func() {
		TruncateTables()
		repo = models.NewKindsRepo()
		env := application.NewEnvironment()
		db := models.NewDatabase(sqlDB, models.Config{
			MigrationsPath: env.ModelMigrationsDir,
		})

		conn = db.Connection().(*models.Connection)
	})

	Describe("Create", func() {
		It("stores the kind record into the database", func() {
			kind := models.Kind{
				ID:          "my-kind",
				Description: "My Kind",
				Critical:    false,
				ClientID:    "my-client",
				TemplateID:  "my-template",
			}

			kind, err := repo.Create(conn, kind)
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
			Expect(kind.TemplateID).To(Equal("my-template"))
			Expect(kind.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
			Expect(kind.UpdatedAt).To(Equal(kind.CreatedAt))
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

			kind1, err := repo.Create(conn, kind1)
			Expect(err).To(BeNil())

			kind2, err = repo.Create(conn, kind2)
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

			kind, err := repo.Create(conn, kind)
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

	Describe("FindByClient", func() {
		It("finds all the records matching the given client ID", func() {
			kind1 := models.Kind{
				ID:          "kind1",
				Description: "kind1-description",
				Critical:    true,
				ClientID:    "client1",
			}

			kind2 := models.Kind{
				ID:          "kind2",
				Description: "kind2-description",
				Critical:    false,
				ClientID:    "client2",
			}

			kind3 := models.Kind{
				ID:          "kind3",
				Description: "kind3-description",
				Critical:    false,
				ClientID:    "client1",
			}

			_, err := repo.Create(conn, kind1)
			if err != nil {
				panic(err)
			}

			_, err = repo.Create(conn, kind2)
			if err != nil {
				panic(err)
			}

			_, err = repo.Create(conn, kind3)
			if err != nil {
				panic(err)
			}

			kindsForClient1, err := repo.FindByClient(conn, "client1")
			if err != nil {
				panic(err)
			}

			Expect(kindsForClient1[0].ID).To(Equal(kind1.ID))
			Expect(kindsForClient1[0].Description).To(Equal(kind1.Description))
			Expect(kindsForClient1[0].Critical).To(Equal(kind1.Critical))
			Expect(kindsForClient1[0].ClientID).To(Equal(kind1.ClientID))
			Expect(kindsForClient1[0].CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))

			Expect(kindsForClient1[1].ID).To(Equal(kind3.ID))
			Expect(kindsForClient1[1].Description).To(Equal(kind3.Description))
			Expect(kindsForClient1[1].Critical).To(Equal(kind3.Critical))
			Expect(kindsForClient1[1].ClientID).To(Equal(kind3.ClientID))
			Expect(kindsForClient1[1].CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
		})

		Context("when there are no kinds for the given client", func() {
			It("returns an empty slice of kinds", func() {
				kindsForClient, err := repo.FindByClient(conn, "i-have-no-kinds")

				Expect(err).ToNot(HaveOccurred())
				Expect(kindsForClient).To(BeEmpty())
			})
		})
	})

	Describe("Update", func() {
		Context("when the template id is meant to be set", func() {
			It("updates the record in the database", func() {
				kind := models.Kind{
					ID:         "my-kind",
					ClientID:   "my-client",
					TemplateID: "my-template",
				}

				kind, err := repo.Create(conn, kind)
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
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})
		})

		Context("when the template id is not meant to be set", func() {
			It("updates the record in the database, using the existing template ID", func() {
				kind := models.Kind{
					ID:         "my-kind",
					ClientID:   "my-client",
					TemplateID: "my-template",
				}

				kind, err := repo.Create(conn, kind)
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
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
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
		})

		Context("when the record exists", func() {
			It("updates the record in the database", func() {
				kind := models.Kind{
					ID:       "my-kind",
					ClientID: "my-client",
				}

				kind, err := repo.Create(conn, kind)
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

			kind, err := repo.Create(conn, kind)
			if err != nil {
				panic(err)
			}

			kindToDelete, err = repo.Create(conn, kindToDelete)
			if err != nil {
				panic(err)
			}

			ignoredKind, err = repo.Create(conn, ignoredKind)
			if err != nil {
				panic(err)
			}

			count, err := repo.Trim(conn, "the-client-id", []string{"other-kind"})

			if err != nil {
				panic(err)
			}

			Expect(count).To(Equal(1))

			_, err = repo.Find(conn, "my-kind", "the-client-id")
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))

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
			kind1, err := repo.Create(conn, models.Kind{
				ID:       "my-kind",
				ClientID: "the-client-id",
			})
			if err != nil {
				panic(err)
			}

			kind2, err := repo.Create(conn, models.Kind{
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
			kind, err := repo.Create(conn, models.Kind{
				ID:         "some-id",
				TemplateID: "some-template",
			})

			_, err = repo.Create(conn, models.Kind{
				ID: "another-id",
			})

			kinds, err := repo.FindAllByTemplateID(conn, "some-template")
			Expect(err).NotTo(HaveOccurred())
			Expect(kinds).To(HaveLen(1))
			Expect(kinds).To(ContainElement(kind))
		})
	})
})
