package services_test

import (
	"testing"

	"github.com/cloudfoundry-incubator/notifications/testing/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebHandlersServicesSuite(t *testing.T) {
	helpers.RegisterFastTokenSigningMethod()

	RegisterFailHandler(Fail)
	RunSpecs(t, "v1/services")
}
