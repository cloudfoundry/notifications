package notify_test

import (
	"fmt"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebV1NotifySuite(t *testing.T) {
	fakes.RegisterFastTokenSigningMethod()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Web V1 Notify Suite")
}

func ExpectToContainMiddlewareStack(actualMiddleware []stack.Middleware, expectedMiddleware ...stack.Middleware) {
	if len(actualMiddleware) != len(expectedMiddleware) {
		Fail(fmt.Sprintf("Expected to see a middleware with %d elements, got %d", len(expectedMiddleware), len(actualMiddleware)))
	}

	for i, ware := range expectedMiddleware {
		Expect(actualMiddleware[i]).To(BeAssignableToTypeOf(ware))
	}

}
