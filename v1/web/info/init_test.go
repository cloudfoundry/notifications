package info_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebV1InfoSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v1/web/info")
}
