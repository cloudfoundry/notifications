package models_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestModelsSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Models Suite")
}
