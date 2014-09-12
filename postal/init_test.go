package postal_test

import (
    "testing"

    "github.com/cloudfoundry-incubator/notifications/test_helpers/fakes"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestPostalSuite(t *testing.T) {
    fakes.RegisterFastTokenSigningMethod()

    RegisterFailHandler(Fail)
    RunSpecs(t, "Postal Suite")
}
