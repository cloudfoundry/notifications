package metrics_test

import (
	"bytes"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueueGauge", func() {
	It("reports the number of items on the queue as a metric", func() {
		job := gobble.Job{}
		queue := fakes.NewQueue()
		output := []byte{}
		buffer := bytes.NewBuffer(output)
		logger := metrics.NewLogger(buffer)
		timer := make(chan time.Time, 10)
		gauge := metrics.NewQueueGauge(queue, logger, timer)

		go gauge.Run()

		Expect(buffer.String()).To(BeEmpty())

		timer <- time.Now()

		Eventually(func() []string {
			return strings.Split(buffer.String(), "\n")
		}).Should(Equal([]string{
			`[METRIC] {"kind":"gauge","payload":{"length":0}}`,
			"",
		}))

		job, err := queue.Enqueue(job)
		Expect(err).NotTo(HaveOccurred())
		timer <- time.Now()

		Eventually(func() []string {
			return strings.Split(buffer.String(), "\n")
		}).Should(Equal([]string{
			`[METRIC] {"kind":"gauge","payload":{"length":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"length":1}}`,
			"",
		}))

		queue.Dequeue(job)
		timer <- time.Now()

		Eventually(func() []string {
			return strings.Split(buffer.String(), "\n")
		}).Should(Equal([]string{
			`[METRIC] {"kind":"gauge","payload":{"length":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"length":1}}`,
			`[METRIC] {"kind":"gauge","payload":{"length":0}}`,
			"",
		}))
	})
})
