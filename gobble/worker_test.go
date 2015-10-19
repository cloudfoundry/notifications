package gobble_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker", func() {
	var (
		queue                 *gobble.Queue
		worker                gobble.Worker
		callbackWasCalledWith gobble.Job
		callback              func(*gobble.Job)
		database              *gobble.DB
	)

	BeforeEach(func() {
		TruncateTables()

		callback = func(job *gobble.Job) {
			callbackWasCalledWith = *job
		}
		database = gobble.NewDatabase(sqlDB)

		queue = gobble.NewQueue(database, gobble.Config{})
		worker = gobble.NewWorker(1, queue, callback)
	})

	AfterEach(func() {
		queue.Close()
	})

	Describe("Perform", func() {
		It("reserves a job, performs the callback, and then dequeues the completed job", func() {
			job, err := queue.Enqueue(gobble.Job{
				Payload: "the-payload",
			}, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			worker.Perform()

			Expect(callbackWasCalledWith.ID).To(Equal(job.ID))

			results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(0))
		})

		It("re-enqueues jobs that are marked for retry", func() {
			callback = func(job *gobble.Job) {
				job.Retry(1 * time.Minute)
			}
			worker = gobble.NewWorker(1, queue, callback)

			job, err := queue.Enqueue(gobble.Job{}, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			worker.Perform()

			results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))

			retriedJob := results[0].(*gobble.Job)
			Expect(retriedJob.ID).To(Equal(job.ID))
			Expect(retriedJob.RetryCount).To(Equal(1))
			Expect(retriedJob.ActiveAt).To(BeTemporally("~", time.Now().Add(1*time.Minute), 1*time.Minute))
		})
	})

	Describe("Work", func() {
		It("works in a loop, and can be stopped", func() {
			queue.Enqueue(gobble.Job{
				Payload: "the-payload",
			}, database.Connection)
			queue.Enqueue(gobble.Job{
				Payload: "the-payload",
			}, database.Connection)

			results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(2))

			worker.Work()

			Eventually(func() (int, error) {
				results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")

				return len(results), err
			}).Should(Equal(0))

			worker.Halt()
		})
	})
})
