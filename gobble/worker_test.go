package gobble_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker", func() {
	var queue *gobble.Queue
	var worker gobble.Worker
	var callbackWasCalledWith gobble.Job
	var callback func(*gobble.Job)

	BeforeEach(func() {
		TruncateTables()

		callback = func(job *gobble.Job) {
			callbackWasCalledWith = *job
		}

		queue = gobble.NewQueue()
		worker = gobble.NewWorker(1, queue, callback)
	})

	AfterEach(func() {
		queue.Close()
	})

	Describe("Perform", func() {
		It("reserves a job, performs the callback, and then dequeues the completed job", func() {
			job, err := queue.Enqueue(gobble.Job{
				Payload: "the-payload",
			})
			if err != nil {
				panic(err)
			}

			worker.Perform()

			Expect(callbackWasCalledWith.ID).To(Equal(job.ID))

			results, err := gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			if err != nil {
				panic(err)
			}

			Expect(len(results)).To(Equal(0))
		})

		It("re-enqueues jobs that are marked for retry", func() {
			callback = func(job *gobble.Job) {
				job.Retry(1 * time.Minute)
			}
			worker = gobble.NewWorker(1, queue, callback)

			job, err := queue.Enqueue(gobble.Job{})
			if err != nil {
				panic(err)
			}

			worker.Perform()

			results, err := gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			if err != nil {
				panic(err)
			}

			Expect(len(results)).To(Equal(1))
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
			})
			queue.Enqueue(gobble.Job{
				Payload: "the-payload",
			})

			results, err := gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			if err != nil {
				panic(err)
			}

			Expect(len(results)).To(Equal(2))

			worker.Work()

			Eventually(func() int {
				results, err := gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
				if err != nil {
					panic(err)
				}

				return len(results)
			}).Should(Equal(0))

			worker.Halt()
		})
	})
})
