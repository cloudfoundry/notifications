package viron_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestViron(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Viron Suite")
}
