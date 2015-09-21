package v2_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestV2Suite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "postal/v2")
}
