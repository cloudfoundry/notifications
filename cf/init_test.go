package cf_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCFSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cf")
}
