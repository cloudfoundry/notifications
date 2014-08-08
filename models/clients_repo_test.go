package models_test

import (
    "time"

    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("ClientsRepo", func() {
    BeforeEach(func() {
        TruncateTables()
    })

    Describe("Create", func() {
        It("stores the client record into the database", func() {
            client := models.Client{
                ID:          "my-client",
                Description: "My Client",
            }

            repo := models.NewClientsRepo()
            client, err := repo.Create(client)
            if err != nil {
                panic(err)
            }

            client, err = repo.Find("my-client")
            if err != nil {
                panic(err)
            }

            Expect(client.ID).To(Equal("my-client"))
            Expect(client.Description).To(Equal("My Client"))
            Expect(client.CreatedAt).To(BeTemporally("~", time.Now(), 2*time.Second))
        })
    })
})
