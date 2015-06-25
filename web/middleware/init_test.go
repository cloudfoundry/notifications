package middleware_test

import (
	"testing"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebMiddlewareSuite(t *testing.T) {
	fakes.RegisterFastTokenSigningMethod()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Middleware Suite")
}
