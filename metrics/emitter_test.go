package metrics_test

import (
	"bytes"
	"log"

	"github.com/cloudfoundry-incubator/notifications/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Emitting metrics", func() {
	var logger *log.Logger
	var buffer *bytes.Buffer

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = metrics.NewLogger(buffer)
	})

	It("can log itself", func() {
		emitter := metrics.NewEmitter(logger)
		emitter.Increment("some.counter")

		message, err := buffer.ReadString('\n')
		Expect(err).NotTo(HaveOccurred())
		Expect(message).To(Equal(`[METRIC] {"kind":"counter","payload":{"name":"some.counter"}}` + "\n"))
	})
})
