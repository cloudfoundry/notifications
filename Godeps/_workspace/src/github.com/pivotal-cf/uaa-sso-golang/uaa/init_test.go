package uaa_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestUAASuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "UAA Suite")
}
