package metrics_test

import (
	"bytes"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"

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
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.length","value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"0"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"1"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"2"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"3"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"4"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"5"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"6"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"7"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"8"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"9"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"10"},"value":0}}`,
			"",
		}))

		job, err := queue.Enqueue(job)
		Expect(err).NotTo(HaveOccurred())
		timer <- time.Now()

		Eventually(func() []string {
			return strings.Split(buffer.String(), "\n")
		}).Should(Equal([]string{
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.length","value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"0"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"1"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"2"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"3"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"4"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"5"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"6"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"7"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"8"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"9"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"10"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.length","value":1}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"0"},"value":1}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"1"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"2"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"3"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"4"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"5"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"6"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"7"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"8"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"9"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"10"},"value":0}}`,
			"",
		}))

		queue.Dequeue(job)
		timer <- time.Now()

		Eventually(func() []string {
			return strings.Split(buffer.String(), "\n")
		}).Should(Equal([]string{
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.length","value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"0"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"1"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"2"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"3"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"4"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"5"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"6"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"7"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"8"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"9"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"10"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.length","value":1}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"0"},"value":1}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"1"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"2"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"3"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"4"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"5"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"6"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"7"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"8"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"9"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"10"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.length","value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"0"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"1"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"2"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"3"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"4"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"5"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"6"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"7"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"8"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"9"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"10"},"value":0}}`,
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
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.length","value":4}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"0"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"1"},"value":3}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"2"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"3"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"4"},"value":1}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"5"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"6"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"7"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"8"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"9"},"value":0}}`,
			`[METRIC] {"kind":"gauge","payload":{"name":"notifications.queue.retry","tags":{"count":"10"},"value":0}}`,
			"",
		}))
	})
})
