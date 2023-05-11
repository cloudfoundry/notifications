package uaa_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUAASuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "uaa")
}
