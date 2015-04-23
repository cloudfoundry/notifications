package strategies_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStrategiesSuite(t *testing.T) {
	fakes.RegisterFastTokenSigningMethod()

	buffer := bytes.NewBuffer([]byte{})
	metricsLogger := metrics.DefaultLogger
	metrics.DefaultLogger = log.New(buffer, "", 0)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Postal Strategies Suite")

	metrics.DefaultLogger = metricsLogger
}
