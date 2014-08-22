package metrics_test

import (
    "bytes"
    "log"

    "github.com/cloudfoundry-incubator/notifications/metrics"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Metric", func() {
    var packageLogger *log.Logger
    var buffer *bytes.Buffer

    BeforeEach(func() {
        buffer = bytes.NewBuffer([]byte{})
        packageLogger = metrics.Logger
        metrics.Logger = log.New(buffer, "", 0)
    })

    AfterEach(func() {
        metrics.Logger = packageLogger
    })

    It("can log itself", func() {
        metric := metrics.NewMetric("counter", map[string]interface{}{
            "name": "test",
        })
        metric.Log()

        message, err := buffer.ReadString('\n')
        if err != nil {
            panic(err)
        }

        Expect(message).To(Equal(`[METRIC] {"kind":"counter","payload":{"name":"test"}}` + "\n"))
    })
})
