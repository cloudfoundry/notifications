package models_test

import (
    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
    It("returns a connection to the database", func() {
        db := models.Database()
        err := db.Connection.Db.Ping()
        Expect(err).To(BeNil())

        _, err = db.Connection.Db.Exec("SHOW TABLES")
        Expect(err).To(BeNil())
    })

    It("returns a single connection only", func() {
        db1 := models.Database()
        db2 := models.Database()

        for i := 0; i < 20; i++ {
            go models.Database()
        }

        Expect(db1).To(Equal(db2))
    })
})
