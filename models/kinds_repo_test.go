package models_test

import (
    "path"
    "time"

    "github.com/cloudfoundry-incubator/notifications/config"
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
        env := config.NewEnvironment()
        migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
        db := models.NewDatabase(env.DatabaseURL, migrationsPath)
        conn = db.Connection().(*models.Connection)
    })

    Describe("Create", func() {
        It("stores the kind record into the database", func() {
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
    })

    Describe("Update", func() {
        It("updates the record in the database", func() {
            kind := models.Kind{
                ID:       "my-kind",
                ClientID: "my-client",
            }

            kind, err := repo.Create(conn, kind)
            if err != nil {
                panic(err)
            }

            kind.Description = "My Kind"
            kind.Critical = true

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
            Expect(kind.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
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
            Expect(err).To(BeAssignableToTypeOf(models.ErrRecordNotFound{}))

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
})
