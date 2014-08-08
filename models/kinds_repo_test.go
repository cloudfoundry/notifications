package models_test

import (
    "time"

    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("KindsRepo", func() {
    BeforeEach(func() {
        TruncateTables()
    })

    Describe("Create", func() {
        It("stores the kind record into the database", func() {
            kind := models.Kind{
                ID:          "my-kind",
                Description: "My Kind",
                Critical:    false,
                ClientID:    "my-client",
            }

            repo := models.NewKindsRepo()

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
})
