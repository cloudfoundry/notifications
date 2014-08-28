package models_test

import (
    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
    var db *models.DB

    BeforeEach(func() {
        db = models.Database()
    })

    It("returns a connection to the database", func() {
        err := db.Connection.Db.Ping()
        Expect(err).To(BeNil())

        _, err = db.Connection.Db.Query("SHOW TABLES")
        Expect(err).To(BeNil())
    })

    It("returns a single connection only", func() {
        db2 := models.Database()

        Expect(db).To(Equal(db2))
    })

    It("has the correct tables", func() {
        err := db.Connection.Db.Ping()
        Expect(err).To(BeNil())

        rows, err := db.Connection.Db.Query("SHOW TABLES")
        Expect(err).To(BeNil())

        tables := []string{}
        for rows.Next() {
            var table string
            err = rows.Scan(&table)
            if err != nil {
                panic(err)
            }
            tables = append(tables, table)
        }
        err = rows.Err()
        if err != nil {
            panic(err)
        }

        rows.Close()

        Expect(tables).To(ContainElement("clients"))
        Expect(tables).To(ContainElement("kinds"))
        Expect(tables).To(ContainElement("receipts"))
        Expect(tables).To(ContainElement("unsubscribes"))
    })
})
