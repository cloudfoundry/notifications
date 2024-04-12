package gobble_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MockHeartbeater struct {
	BeatCall struct {
		Receives struct {
			Job gobble.Job
		}
	}

	HaltCall struct {
		WasCalled bool
	}
}

func (b *MockHeartbeater) Beat(job *gobble.Job) {
	b.BeatCall.Receives.Job = *job
}

func (b *MockHeartbeater) Halt() {
	b.HaltCall.WasCalled = true
}

var _ = Describe("Worker", func() {
	var (
		queue                 *gobble.Queue
		worker                gobble.Worker
		heartbeater           *MockHeartbeater
		callbackWasCalledWith gobble.Job
		callback              func(*gobble.Job)
		database              *gobble.DB
		clock                 *mocks.Clock
	)

	BeforeEach(func() {
		TruncateTables()

		callback = func(job *gobble.Job) {
			callbackWasCalledWith = *job
		}
		database = gobble.NewDatabase(sqlDB)
		clock = &mocks.Clock{}
		clock.NowCall.Returns.Time = time.Now().UTC().Truncate(time.Second)

		queue = gobble.NewQueue(database, clock, gobble.Config{MaxQueueLength: 1000})
		heartbeater = &MockHeartbeater{}
		worker = gobble.NewWorker(1, queue, callback, heartbeater)
	})

	AfterEach(func() {
		queue.Close()
	})

	Describe("Perform", func() {
		It("reserves a job, performs the callback, and then dequeues the completed job", func() {
			job, err := queue.Enqueue(&gobble.Job{
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
			worker = gobble.NewWorker(1, queue, callback, heartbeater)

			job, err := queue.Enqueue(&gobble.Job{}, database.Connection)
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

		It("heartbeats for job ownership while the job executes", func() {
			job, err := queue.Enqueue(&gobble.Job{
				Payload: "the-payload",
			}, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			hold := make(chan struct{})

			callback = func(*gobble.Job) {
				<-hold
			}
			worker = gobble.NewWorker(2, queue, callback, heartbeater)

			go worker.Perform()

			Eventually(func() int {
				return heartbeater.BeatCall.Receives.Job.ID
			}).Should(Equal(job.ID))

			hold <- struct{}{}

			Eventually(func() bool {
				return heartbeater.HaltCall.WasCalled
			}).Should(BeTrue())
		})
	})

	Describe("Work", func() {
		It("works in a loop, and can be stopped", func() {
			worker = gobble.NewWorker(1, queue, callback, &MockHeartbeater{})

			queue.Enqueue(&gobble.Job{
				Payload: "the-payload",
			}, database.Connection)
			queue.Enqueue(&gobble.Job{
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
