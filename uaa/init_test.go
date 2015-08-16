package uaa_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUAASuite(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	metricsLogger := metrics.DefaultLogger
	metrics.DefaultLogger = log.New(buffer, "", 0)

	RegisterFailHandler(Fail)
	RunSpecs(t, "uaa")
	metrics.DefaultLogger = metricsLogger
}
