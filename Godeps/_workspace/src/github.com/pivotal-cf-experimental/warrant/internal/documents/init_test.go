package documents_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDocumentsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Documents Suite")
}
