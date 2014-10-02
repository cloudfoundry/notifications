package cryptography_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestCryptographySuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Cryptography Suite")
}
