package models_test

import (
    "testing"

    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestModelsSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Models Suite")
}

func TruncateTables() {
    err := models.Database().Connection().TruncateTables()
    if err != nil {
        panic(err)
    }
}
