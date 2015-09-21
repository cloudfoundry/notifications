package v1_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestV1Suite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "postal/v1")
}
