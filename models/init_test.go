package models_test

import (
    "testing"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestModelsSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Models Suite")
}

func TruncateTables() {
    env := config.NewEnvironment()
    db := models.NewDatabase(env.DatabaseURL)
    connection := db.Connection().(*models.Connection)
    err := connection.TruncateTables()
    if err != nil {
        panic(err)
    }
}
