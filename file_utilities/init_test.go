package file_utilities_test

import (
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestFileUtilitiesSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "File Utilities Suite")
}
