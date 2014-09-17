package postal_test

import (
    "bytes"
    "log"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    "github.com/cloudfoundry-incubator/notifications/test_helpers/fakes"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestPostalSuite(t *testing.T) {
    fakes.RegisterFastTokenSigningMethod()

    buffer := bytes.NewBuffer([]byte{})
    metricsLogger := metrics.Logger
    metrics.Logger = log.New(buffer, "", 0)

    RegisterFailHandler(Fail)
    RunSpecs(t, "Postal Suite")

    metrics.Logger = metricsLogger
}
