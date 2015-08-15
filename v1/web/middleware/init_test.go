package middleware_test

import (
	"testing"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebMiddlewareSuite(t *testing.T) {
	fakes.RegisterFastTokenSigningMethod()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Web Middleware Suite")
}
