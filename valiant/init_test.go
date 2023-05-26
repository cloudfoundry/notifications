package valiant_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestValiantSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "valiant")
}
