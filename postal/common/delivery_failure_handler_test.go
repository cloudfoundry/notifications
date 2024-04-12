package common_test

import (
	"bytes"
	"time"

	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeliveryFailureHandler", func() {
	var (
		job     *mocks.GobbleJob
		buffer  *bytes.Buffer
		logger  lager.Logger
		handler common.DeliveryFailureHandler
	)

	BeforeEach(func() {
		job = mocks.NewGobbleJob()
		buffer = bytes.NewBuffer([]byte{})
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.INFO))

		handler = common.NewDeliveryFailureHandler(10)
	})

	It("retries the job using an exponential backoff algorithm", func() {
		backoffDurations := map[int]time.Duration{
			0: 1 * time.Minute,
			1: 2 * time.Minute,
			2: 4 * time.Minute,
			3: 8 * time.Minute,
			4: 16 * time.Minute,
			5: 32 * time.Minute,
			6: 64 * time.Minute,
			7: 128 * time.Minute,
			8: 256 * time.Minute,
			9: 512 * time.Minute,
		}

		for retryCount, duration := range backoffDurations {
			job.StateCall.Returns.Count = retryCount

			handler.Handle(job, logger)

			Expect(job.RetryCall.Receives.Duration).To(Equal(duration))
		}

		job.StateCall.Returns.Count = 10
		job.RetryCall.WasCalled = false

		handler.Handle(job, logger)

		Expect(job.RetryCall.WasCalled).To(BeFalse())
	})

	It("gives up after 9 retries", func() {
		job.StateCall.Returns.Count = 10

		handler.Handle(job, logger)

		Expect(job.RetryCall.WasCalled).To(BeFalse())
	})

	It("logs the retry attempt", func() {
		expectedActiveAt := time.Now().Truncate(time.Second)
		job.StateCall.Returns.Time = expectedActiveAt
		job.StateCall.Returns.Count = 4

		handler.Handle(job, logger)

		lines, err := parseLogLines(buffer.Bytes())
		Expect(err).NotTo(HaveOccurred())
		Expect(lines).To(HaveLen(1))

		line := lines[0]
		Expect(line.Source).To(Equal("notifications"))
		Expect(line.Message).To(Equal("notifications.delivery-failed-retrying"))
		Expect(line.LogLevel).To(Equal(int(lager.INFO)))
		Expect(line.Data).To(HaveKeyWithValue("retry_count", float64(4)))

		Expect(line.Data).To(HaveKey("active_at"))
		activeAt, err := time.Parse(time.RFC3339, line.Data["active_at"].(string))
		Expect(err).NotTo(HaveOccurred())
		Expect(activeAt.UTC()).To(Equal(expectedActiveAt.UTC()))
	})
})
