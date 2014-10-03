package models_test

import (
    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("GlobalUnsubscribesRepo", func() {
    var repo models.GlobalUnsubscribesRepo
    var conn *models.Connection

    Describe("Set/Get", func() {
        BeforeEach(func() {
            TruncateTables()
            repo = models.NewGlobalUnsubscribesRepo()
            conn = models.Database().Connection()
        })

        It("sets the global unsubscribe field for a user, allowing it to be retrieved later", func() {
            err := repo.Set(conn, "my-user", true)
            if err != nil {
                panic(err)
            }

            unsubscribed, err := repo.Get(conn, "my-user")
            if err != nil {
                panic(err)
            }

            Expect(unsubscribed).To(BeTrue())

            err = repo.Set(conn, "my-user", false)
            if err != nil {
                panic(err)
            }

            unsubscribed, err = repo.Get(conn, "my-user")
            if err != nil {
                panic(err)
            }

            Expect(unsubscribed).To(BeFalse())
        })
    })
})
