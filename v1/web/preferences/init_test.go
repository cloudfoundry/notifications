package preferences_test

import (
	"fmt"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebV1PreferencesSuite(t *testing.T) {
	helpers.RegisterFastTokenSigningMethod()

	RegisterFailHandler(Fail)
	RunSpecs(t, "v1/web/preferences")
}

func ExpectToContainMiddlewareStack(actualMiddleware []stack.Middleware, expectedMiddleware ...stack.Middleware) {
	if len(actualMiddleware) != len(expectedMiddleware) {
		Fail(fmt.Sprintf("Expected to see a middleware with %d elements, got %d", len(expectedMiddleware), len(actualMiddleware)))
	}

	for i, ware := range expectedMiddleware {
		Expect(actualMiddleware[i]).To(BeAssignableToTypeOf(ware))
	}

}
