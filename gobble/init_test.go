package gobble_test

import (
    "testing"

    "github.com/cloudfoundry-incubator/notifications/gobble"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestGobbleSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Gobble Suite")
}

func TruncateTables() {
    gobble.Database().Connection.TruncateTables()
}
