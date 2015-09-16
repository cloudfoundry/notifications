package clients_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebV1ClientsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v1/web/clients")
}
