package horde_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAudiencesSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v2/horde")
}
