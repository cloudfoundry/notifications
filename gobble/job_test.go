package gobble_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Job", func() {
	Describe("NewJob", func() {
		It("creates an instance of Job with a serialized payload", func() {
			data := map[string]string{
				"test":    "testing a new job",
				"example": "another field",
			}

			job := gobble.NewJob(data)

			Expect(job).To(BeAssignableToTypeOf(&gobble.Job{}))
			Expect(job.Payload).To(Equal(`{"example":"another field","test":"testing a new job"}`))
		})
	})

	Describe("Unmarshal", func() {
		It("unmarshals the payload into the given object", func() {
			data := map[string]string{
				"test":    "testing a new job",
				"example": "another field",
			}

			job := gobble.NewJob(data)

			var payload map[string]string
			err := job.Unmarshal(&payload)
			Expect(err).NotTo(HaveOccurred())

			Expect(payload).To(Equal(data))
		})
	})

	Describe("Retry", func() {
		It("sets up the job to be retried", func() {
			job := gobble.NewJob("the data")
			job.RetryCount = 1
			job.WorkerID = "my-id"
			job.ActiveAt = time.Now().Add(-5 * time.Minute)

			job.Retry(10 * time.Minute)

			Expect(job.WorkerID).To(Equal(""))
			Expect(job.RetryCount).To(Equal(2))
			Expect(job.ActiveAt).To(BeTemporally("~", time.Now().Add(10*time.Minute), 10*time.Second))
			Expect(job.ShouldRetry).To(BeTrue())
		})
	})

	Describe("State", func() {
		It("returns the current retry count and active at values", func() {
			expectedActiveAt := time.Now().Add(-5 * time.Minute)

			job := gobble.NewJob("the data")
			job.RetryCount = 4
			job.ActiveAt = expectedActiveAt

			retryCount, activeAt := job.State()
			Expect(retryCount).To(Equal(4))
			Expect(activeAt).To(Equal(expectedActiveAt))
		})
	})
})
