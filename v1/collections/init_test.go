package collections_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebHandlersServicesSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v1/collections")
}
