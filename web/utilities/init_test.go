package utilities_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestUtilitiesSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Utilities Suite")
}
