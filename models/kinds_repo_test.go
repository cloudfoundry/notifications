package models_test

import (
    "time"

    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/coopernurse/gorp"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("KindsRepo", func() {
    var repo models.KindsRepo
    var conn *gorp.DbMap

    BeforeEach(func() {
        TruncateTables()
        repo = models.NewKindsRepo()
        conn = models.Database().Connection
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

            kind, err = repo.Find(conn, "my-kind")
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

    Describe("Update", func() {
        It("updates the record in the database", func() {
            kind := models.Kind{
                ID: "my-kind",
            }

            kind, err := repo.Create(conn, kind)
            if err != nil {
                panic(err)
            }

            kind.Description = "My Kind"
            kind.Critical = true
            kind.ClientID = "my-client"

            kind, err = repo.Update(conn, kind)
            if err != nil {
                panic(err)
            }

            kind, err = repo.Find(conn, "my-kind")
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

                kind, err = repo.Find(conn, "my-kind")
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
                    ID: "my-kind",
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

                kind, err = repo.Find(conn, "my-kind")
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

            _, err = repo.Find(conn, "my-kind")
            Expect(err).To(BeAssignableToTypeOf(models.ErrRecordNotFound{}))

            _, err = repo.Find(conn, "ignored-kind")
            if err != nil {
                panic(err)
            }

            _, err = repo.Find(conn, "other-kind")
            if err != nil {
                panic(err)
            }
        })
    })
})
