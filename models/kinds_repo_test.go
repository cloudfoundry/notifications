package models_test

import (
    "time"

    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("KindsRepo", func() {
    var repo models.KindsRepo

    BeforeEach(func() {
        TruncateTables()
        repo = models.NewKindsRepo()
    })

    Describe("Create", func() {
        It("stores the kind record into the database", func() {
            kind := models.Kind{
                ID:          "my-kind",
                Description: "My Kind",
                Critical:    false,
                ClientID:    "my-client",
            }

            kind, err := repo.Create(kind)
            if err != nil {
                panic(err)
            }

            kind, err = repo.Find("my-kind")
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

            kind, err := repo.Create(kind)
            if err != nil {
                panic(err)
            }

            kind.Description = "My Kind"
            kind.Critical = true
            kind.ClientID = "my-client"

            kind, err = repo.Update(kind)
            if err != nil {
                panic(err)
            }

            kind, err = repo.Find("my-kind")
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

                kind, err := repo.Upsert(kind)
                if err != nil {
                    panic(err)
                }

                kind, err = repo.Find("my-kind")
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

                kind, err := repo.Create(kind)
                if err != nil {
                    panic(err)
                }

                kind = models.Kind{
                    ID:          "my-kind",
                    Description: "My Kind",
                    Critical:    true,
                    ClientID:    "my-client",
                }

                kind, err = repo.Upsert(kind)
                if err != nil {
                    panic(err)
                }

                kind, err = repo.Find("my-kind")
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
})
