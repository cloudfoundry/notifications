package models_test

import (
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
    var env config.Environment
    var db *models.DB
    var connection *models.Connection

    BeforeEach(func() {
        env = config.NewEnvironment()
        db = models.NewDatabase(env.DatabaseURL)
        connection = db.Connection().(*models.Connection)
    })

    It("returns a connection to the database", func() {
        err := connection.Db.Ping()
        Expect(err).To(BeNil())

        _, err = connection.Db.Query("SHOW TABLES")
        Expect(err).To(BeNil())
    })

    It("returns a single connection only", func() {
        db2 := models.NewDatabase(env.DatabaseURL)

        Expect(db).To(Equal(db2))
    })

    It("has the correct tables", func() {
        err := connection.Db.Ping()
        Expect(err).To(BeNil())

        rows, err := connection.Db.Query("SHOW TABLES")
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

    Describe("Connection", func() {
        It("returns a Connection", func() {
            connection := db.Connection()
            Expect(connection).To(BeAssignableToTypeOf(&models.Connection{}))
        })
    })
})
