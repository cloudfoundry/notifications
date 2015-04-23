package gobble_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queue", func() {
	var (
		queue    *gobble.Queue
		database *gobble.DB
	)

	BeforeEach(func() {
		TruncateTables()
		database = gobble.NewDatabase(sqlDB)

		queue = gobble.NewQueue(database, gobble.Config{
			WaitMaxDuration: 50 * time.Millisecond,
		})
	})

	AfterEach(func() {
		queue.Close()
	})

	Describe("Enqueue", func() {
		It("sticks the job in the database table", func() {
			job := gobble.NewJob(map[string]bool{
				"testing": true,
			})

			job, err := queue.Enqueue(job)
			Expect(err).NotTo(HaveOccurred())

			results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			if err != nil {
				panic(err)
			}

			jobs := []gobble.Job{}
			for _, result := range results {
				jobs = append(jobs, *(result.(*gobble.Job)))
			}

			Expect(jobs).To(HaveLen(1))
			Expect(jobs).To(ContainElement(job))
		})
	})

	Describe("Requeue", func() {
		It("updates the queue in the database", func() {
			job := gobble.NewJob(map[string]bool{
				"testing": true,
			})

			job, err := queue.Enqueue(job)
			if err != nil {
				panic(err)
			}

			job.RetryCount = 5

			queue.Requeue(job)

			reloadedJob := gobble.Job{}
			err = database.Connection.SelectOne(&reloadedJob, "SELECT * FROM `jobs` where id = ?", job.ID)
			if err != nil {
				panic(err)
			}

			Expect(reloadedJob.ID).To(Equal(job.ID))
			Expect(reloadedJob.RetryCount).To(Equal(5))
		})
	})

	Describe("Reserve", func() {
		It("reserves a job in the database", func() {
			job := gobble.Job{
				Payload: "something",
			}

			err := database.Connection.Insert(&job)
			if err != nil {
				panic(err)
			}

			jobChannel := queue.Reserve("workerId")
			reservedJob := <-jobChannel

			Expect(reservedJob.ID).To(Equal(job.ID))
			Expect(reservedJob.ActiveAt).To(BeTemporally("~", time.Now(), 100*time.Millisecond))
		})

		It("keeps trying to reserve a job until one becomes available", func() {
			jobChannel := queue.Reserve("my-id")

			Consistently(jobChannel).ShouldNot(Receive())

			job, err := queue.Enqueue(gobble.Job{
				Payload: "hello",
			})
			if err != nil {
				panic(err)
			}

			var reservedJob gobble.Job
			Eventually(jobChannel).Should(Receive(&reservedJob))

			Expect(reservedJob.ID).To(Equal(job.ID))
		})

		It("ensures a job can only be reserved by a single worker", func() {
			for i := 0; i < 100; i++ {
				queue.Enqueue(gobble.Job{})
			}

			done := make(chan bool)
			reserveJob := func(id string) {
				for i := 0; i < 50; i++ {
					jobChan := queue.Reserve(id)
					<-jobChan
					<-time.After(1 * time.Millisecond)
				}
				done <- true
			}

			go reserveJob("worker-1")
			go reserveJob("worker-2")

			Eventually(done, 1*time.Second).Should(Receive())
			Eventually(done, 1*time.Second).Should(Receive())

			results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs` WHERE `worker_id` = ''")
			if err != nil {
				panic(err)
			}

			Expect(results).To(HaveLen(0))
		})

		It("picks the first job that is active", func() {
			queue.Enqueue(gobble.Job{
				ActiveAt: time.Now().Add(1 * time.Hour),
			})
			job2, err := queue.Enqueue(gobble.Job{})
			if err != nil {
				panic(err)
			}

			job := <-queue.Reserve("worker-id")

			Expect(job.ID).To(Equal(job2.ID))
		})

		Context("when the worker id is set", func() {
			Context("when active_at is in the future", func() {
				It("should not grab the job", func() {
					_, err := queue.Enqueue(gobble.Job{
						WorkerID: "some-worker",
						ActiveAt: time.Now().Add(1 * time.Minute),
					})
					Expect(err).NotTo(HaveOccurred())
					Consistently(queue.Reserve("some-other-worker")).ShouldNot(Receive())
				})
			})

			Context("when active_at is a little bit in the past", func() {
				It("should not grab the job", func() {
					_, err := queue.Enqueue(gobble.Job{
						WorkerID: "some-worker",
						ActiveAt: time.Now().Add(-1 * time.Minute),
					})
					Expect(err).NotTo(HaveOccurred())
					Consistently(queue.Reserve("some-other-worker")).ShouldNot(Receive())
				})
			})

			Context("when active_at is very far in the past", func() {
				It("should grab the job", func() {
					_, err := queue.Enqueue(gobble.Job{
						WorkerID: "some-worker",
						ActiveAt: time.Now().Add(-2 * time.Minute),
					})
					Expect(err).NotTo(HaveOccurred())
					Eventually(queue.Reserve("some-other-worker")).Should(Receive())
				})
			})
		})

		Context("when the worker id is not set", func() {
			Context("when active_at is in the future", func() {
				It("should not grab the job", func() {
					_, err := queue.Enqueue(gobble.Job{
						ActiveAt: time.Now().Add(1 * time.Minute),
					})
					Expect(err).NotTo(HaveOccurred())
					Consistently(queue.Reserve("some-other-worker")).ShouldNot(Receive())
				})
			})

			Context("when active_at is in the past", func() {
				It("should grab the job", func() {
					_, err := queue.Enqueue(gobble.Job{
						ActiveAt: time.Now().Add(-1 * time.Minute),
					})
					Expect(err).NotTo(HaveOccurred())
					Eventually(queue.Reserve("some-other-worker")).Should(Receive())
				})
			})
		})
	})

	Describe("Dequeue", func() {
		It("deletes the job from the queue", func() {
			job, err := queue.Enqueue(gobble.Job{})
			Expect(err).NotTo(HaveOccurred())

			results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))

			queue.Dequeue(job)

			results, err = database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(0))
		})
	})

	Describe("Len", func() {
		It("returns the length of the queue", func() {
			job, err := queue.Enqueue(gobble.Job{})
			Expect(err).NotTo(HaveOccurred())

			length, err := queue.Len()
			Expect(err).NotTo(HaveOccurred())
			Expect(length).To(Equal(1))

			queue.Dequeue(job)

			length, err = queue.Len()
			Expect(err).NotTo(HaveOccurred())
			Expect(length).To(Equal(0))
		})
	})

	Describe("RetryQueueLengths", func() {
		It("returns information about the length of the queue grouped by retry count", func() {
			_, err := queue.Enqueue(gobble.Job{})
			Expect(err).NotTo(HaveOccurred())

			for i := 0; i < 3; i++ {
				_, err = queue.Enqueue(gobble.Job{RetryCount: 1})
				Expect(err).NotTo(HaveOccurred())
			}

			_, err = queue.Enqueue(gobble.Job{RetryCount: 4})
			Expect(err).NotTo(HaveOccurred())

			for i := 0; i < 2; i++ {
				_, err = queue.Enqueue(gobble.Job{RetryCount: 5})
				Expect(err).NotTo(HaveOccurred())
			}

			lengths, err := queue.RetryQueueLengths()
			Expect(err).NotTo(HaveOccurred())
			Expect(lengths).To(Equal(map[int]int{
				0: 1,
				1: 3,
				4: 1,
				5: 2,
			}))
		})
	})
})
