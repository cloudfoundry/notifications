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
	var (
		gauge  metrics.QueueGauge
		timer  chan time.Time
		queue  *fakes.Queue
		buffer *bytes.Buffer
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		queue = fakes.NewQueue()
		timer = make(chan time.Time, 10)
		gauge = metrics.NewQueueGauge(queue, metrics.NewLogger(buffer), timer)
	})

	It("reports the number of items on the queue as a metric", func() {
		job := gobble.Job{}

		go gauge.Run()

		Expect(buffer.String()).To(BeEmpty())

		timer <- time.Now()

		Eventually(func() []string {
			return strings.Split(buffer.String(), "\n")
		}).Should(Equal([]string{
			`[METRIC] {"kind":"gauge","payload":{"queue-length":0}}`,
			"",
		}))

		job, err := queue.Enqueue(job)
		Expect(err).NotTo(HaveOccurred())
		timer <- time.Now()

		Eventually(func() []string {
			return strings.Split(buffer.String(), "\n")
		}).Should(Equal([]string{
			`[METRIC] {"kind":"gauge","payload":{"queue-length":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"queue-retry-counts.0":1}}`,
			`[METRIC] {"kind":"gauge","payload":{"queue-length":1}}`,
			"",
		}))

		queue.Dequeue(job)
		timer <- time.Now()

		Eventually(func() []string {
			return strings.Split(buffer.String(), "\n")
		}).Should(Equal([]string{
			`[METRIC] {"kind":"gauge","payload":{"queue-length":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"queue-retry-counts.0":1}}`,
			`[METRIC] {"kind":"gauge","payload":{"queue-length":1}}`,
			`[METRIC] {"kind":"gauge","payload":{"queue-length":0}}`,
			"",
		}))
	})

	It("reports the number of jobs grouped by retry count", func() {
		go gauge.Run()

		Expect(buffer.String()).To(BeEmpty())

		_, err := queue.Enqueue(gobble.Job{RetryCount: 4})
		Expect(err).NotTo(HaveOccurred())

		for i := 0; i < 3; i++ {
			_, err = queue.Enqueue(gobble.Job{RetryCount: 1})
			Expect(err).NotTo(HaveOccurred())
		}

		timer <- time.Now()

		Eventually(func() []string {
			return strings.Split(buffer.String(), "\n")
		}).Should(ConsistOf([]string{
			`[METRIC] {"kind":"gauge","payload":{"queue-length":4}}`,
			`[METRIC] {"kind":"gauge","payload":{"queue-retry-counts.4":1}}`,
			`[METRIC] {"kind":"gauge","payload":{"queue-retry-counts.1":3}}`,
			"",
		}))
	})
})
