package postal_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPostalSuite(t *testing.T) {
	helpers.RegisterFastTokenSigningMethod()

	buffer := bytes.NewBuffer([]byte{})
	metricsLogger := metrics.DefaultLogger
	metrics.DefaultLogger = log.New(buffer, "", 0)

	RegisterFailHandler(Fail)
	RunSpecs(t, "postal")

	metrics.DefaultLogger = metricsLogger
}
