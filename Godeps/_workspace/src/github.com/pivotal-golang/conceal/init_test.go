package conceal_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestConcealSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Conceal Suite")
}
