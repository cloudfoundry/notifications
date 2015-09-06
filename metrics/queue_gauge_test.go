package metrics_test

import (
	"bytes"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/metrics"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueueGauge", func() {
	var (
		gauge  metrics.QueueGauge
		timer  chan time.Time
		queue  *mocks.Queue
		buffer *bytes.Buffer
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		queue = mocks.NewQueue()
		timer = make(chan time.Time, 10)
		gauge = metrics.NewQueueGauge(queue, metrics.NewLogger(buffer), timer)
	})

	It("reports the number of items on the queue as a metric", func() {
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

		queue.RetryQueueLengthsCall.Returns.Lengths = map[int]int{
			0: 1,
		}
		queue.LenCall.Returns.Length = 1
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

		queue.RetryQueueLengthsCall.Returns.Lengths = map[int]int{}
		queue.LenCall.Returns.Length = 0
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

		queue.RetryQueueLengthsCall.Returns.Lengths = map[int]int{
			1: 3,
			4: 1,
		}
		queue.LenCall.Returns.Length = 4
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
