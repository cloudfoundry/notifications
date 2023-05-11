package postal_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPostalSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "postal")
}
