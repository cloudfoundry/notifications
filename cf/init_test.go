package cf_test

import (
    "bytes"
    "log"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/metrics"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestCFSuite(t *testing.T) {
    buffer := bytes.NewBuffer([]byte{})
    metricsLogger := metrics.Logger
    metrics.Logger = log.New(buffer, "", 0)

    RegisterFailHandler(Fail)
    RunSpecs(t, "CF Suite")

    metrics.Logger = metricsLogger
}
