package utilities_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUtilitiesSuite(t *testing.T) {
	fakes.RegisterFastTokenSigningMethod()

	buffer := bytes.NewBuffer([]byte{})
	metricsLogger := metrics.DefaultLogger
	metrics.DefaultLogger = log.New(buffer, "", 0)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Postal Utilities Suite")

	metrics.DefaultLogger = metricsLogger
}
