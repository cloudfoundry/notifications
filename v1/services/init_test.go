package services_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWebHandlersServicesSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v1/services")
}
