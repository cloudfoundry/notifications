package webutil_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebutilSuite(t *testing.T) {
	fakes.RegisterFastTokenSigningMethod()

	buffer := bytes.NewBuffer([]byte{})
	metricsLogger := metrics.DefaultLogger
	metrics.DefaultLogger = log.New(buffer, "", 0)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Webutil Suite")

	metrics.DefaultLogger = metricsLogger
}
