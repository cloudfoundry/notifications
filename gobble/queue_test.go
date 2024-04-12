package gobble_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queue", func() {
	var (
		queue    *gobble.Queue
		database *gobble.DB
		clock    *mocks.Clock
	)

	BeforeEach(func() {
		TruncateTables()
		database = gobble.NewDatabase(sqlDB)
		clock = &mocks.Clock{}
		clock.NowCall.Returns.Time = time.Now().UTC().Truncate(time.Second)

		queue = gobble.NewQueue(database, clock, gobble.Config{
			WaitMaxDuration: 50 * time.Millisecond,
			MaxQueueLength:  1000,
		})
	})

	AfterEach(func() {
		queue.Close()
	})

	Describe("Enqueue", func() {
		It("sticks the job in the database table using a connection that is passed in", func() {
			job := gobble.NewJob(map[string]bool{
				"testing": true,
			})

			job, err := queue.Enqueue(job, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			jobs := []*gobble.Job{}
			_, err = database.Connection.Select(&jobs, "SELECT * FROM `jobs`")
			Expect(err).NotTo(HaveOccurred())

			Expect(jobs).To(HaveLen(1))
			Expect(jobs).To(ContainElement(job))
		})

		It("doesn't surpase max queue length", func() {
			job2 := gobble.Job{
				Payload:  "something",
				ActiveAt: time.Now().UTC().Truncate(time.Second),
			}

			err := database.Connection.Insert(&job2)
			Expect(err).NotTo(HaveOccurred())

			queue = gobble.NewQueue(database, clock, gobble.Config{
				WaitMaxDuration: 50 * time.Millisecond,
				MaxQueueLength:  1,
			})
			job := gobble.NewJob(map[string]bool{
				"testing": true,
			})

			nilJob, err := queue.Enqueue(job, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			Expect(nilJob).To(BeNil())

			jobs := []*gobble.Job{}
			_, err = database.Connection.Select(&jobs, "SELECT * FROM `jobs`")
			Expect(err).NotTo(HaveOccurred())

			Expect(jobs).To(HaveLen(1))
			Expect(jobs).ToNot(ContainElement(job))
		})

		Context("when the transaction is not commited", func() {
			It("should not put things in the database", func() {
				job := gobble.NewJob(map[string]bool{
					"testing": true,
				})

				transaction := db.NewTransaction(&db.Connection{DbMap: database.Connection})
				transaction.Begin()

				job, err := queue.Enqueue(job, transaction)
				Expect(err).NotTo(HaveOccurred())
				transaction.Rollback()

				results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
				Expect(err).NotTo(HaveOccurred())

				jobs := []gobble.Job{}
				for _, result := range results {
					jobs = append(jobs, *(result.(*gobble.Job)))
				}

				Expect(jobs).To(HaveLen(0))
			})
		})
	})

	Describe("Requeue", func() {
		It("updates the queue in the database", func() {
			job := gobble.NewJob(map[string]bool{
				"testing": true,
			})

			job, err := queue.Enqueue(job, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			job.RetryCount = 5

			queue.Requeue(job)

			reloadedJob := gobble.Job{}
			err = database.Connection.SelectOne(&reloadedJob, "SELECT * FROM `jobs` where id = ?", job.ID)
			Expect(err).NotTo(HaveOccurred())

			Expect(reloadedJob.ID).To(Equal(job.ID))
			Expect(reloadedJob.RetryCount).To(Equal(5))
		})

		It("deletes the job if above max_queue_length", func() {
			queue = gobble.NewQueue(database, clock, gobble.Config{
				WaitMaxDuration: 50 * time.Millisecond,
				MaxQueueLength:  2,
			})
			job := gobble.NewJob(map[string]bool{
				"testing": true,
			})

			job, err := queue.Enqueue(job, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			job2 := gobble.Job{
				Payload:  "something",
				ActiveAt: time.Now().UTC().Truncate(time.Second),
			}

			err = database.Connection.Insert(&job2)
			Expect(err).NotTo(HaveOccurred())

			job.RetryCount = 5

			queue.Requeue(job)

			reloadedJob := gobble.Job{}
			err = database.Connection.SelectOne(&reloadedJob, "SELECT * FROM `jobs` where id = ?", job.ID)
			Expect(err).To(HaveOccurred())
			len, err := queue.Len()
			Expect(err).ToNot(HaveOccurred())
			Expect(len).To(Equal(1))
		})
	})

	Describe("Reserve", func() {
		It("reserves a job in the database", func() {
			job := gobble.Job{
				Payload:  "something",
				ActiveAt: time.Now().UTC().Truncate(time.Second),
			}

			err := database.Connection.Insert(&job)
			Expect(err).NotTo(HaveOccurred())

			jobChannel := queue.Reserve("workerId")
			reservedJob := <-jobChannel

			Expect(reservedJob.ID).To(Equal(job.ID))
			Expect(reservedJob.ActiveAt).To(BeTemporally("~", time.Now(), 250*time.Millisecond))
		})

		It("keeps trying to reserve a job until one becomes available", func() {
			jobChannel := queue.Reserve("my-id")

			Consistently(jobChannel).ShouldNot(Receive())

			job, err := queue.Enqueue(&gobble.Job{
				Payload: "hello",
			}, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			var reservedJob *gobble.Job
			Eventually(jobChannel).Should(Receive(&reservedJob))

			Expect(reservedJob.ID).To(Equal(job.ID))
		})

		It("ensures a job can only be reserved by a single worker", func() {
			for i := 0; i < 100; i++ {
				_, err := queue.Enqueue(&gobble.Job{}, database.Connection)
				Expect(err).ToNot(HaveOccurred())
			}

			done := make(chan bool)
			reserveJob := func(id string) {
				for i := 0; i < 50; i++ {
					<-queue.Reserve(id)
					<-time.After(1 * time.Millisecond)
				}
				done <- true
			}

			go reserveJob("worker-1")
			go reserveJob("worker-2")

			Eventually(done, 30*time.Second).Should(Receive())
			Eventually(done, 30*time.Second).Should(Receive())

			results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs` WHERE `worker_id` = ''")
			Expect(err).NotTo(HaveOccurred())

			Expect(results).To(HaveLen(0))
		})

		It("picks the first job that is active", func() {
			queue.Enqueue(&gobble.Job{
				ActiveAt: time.Now().Add(1 * time.Hour),
			}, database.Connection)
			job2, err := queue.Enqueue(&gobble.Job{}, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			job := <-queue.Reserve("worker-id")

			Expect(job.ID).To(Equal(job2.ID))
		})

		Context("when the worker id is set", func() {
			Context("when active_at is in the future", func() {
				It("should not grab the job", func() {
					_, err := queue.Enqueue(&gobble.Job{
						WorkerID: "some-worker",
						ActiveAt: time.Now().Add(1 * time.Minute),
					}, database.Connection)
					Expect(err).NotTo(HaveOccurred())
					Consistently(queue.Reserve("some-other-worker")).ShouldNot(Receive())
				})
			})

			Context("when active_at is a little bit in the past", func() {
				It("should not grab the job", func() {
					_, err := queue.Enqueue(&gobble.Job{
						WorkerID: "some-worker",
						ActiveAt: time.Now().Add(-1 * time.Minute),
					}, database.Connection)
					Expect(err).NotTo(HaveOccurred())
					Consistently(queue.Reserve("some-other-worker")).ShouldNot(Receive())
				})
			})

			Context("when active_at is very far in the past", func() {
				It("should grab the job", func() {
					_, err := queue.Enqueue(&gobble.Job{
						WorkerID: "some-worker",
						ActiveAt: time.Now().Add(-2 * time.Minute),
					}, database.Connection)
					Expect(err).NotTo(HaveOccurred())
					Eventually(queue.Reserve("some-other-worker")).Should(Receive())
				})
			})
		})

		Context("when the worker id is not set", func() {
			Context("when active_at is in the future", func() {
				It("should not grab the job", func() {
					_, err := queue.Enqueue(&gobble.Job{
						ActiveAt: time.Now().Add(1 * time.Minute),
					}, database.Connection)
					Expect(err).NotTo(HaveOccurred())
					Consistently(queue.Reserve("some-other-worker")).ShouldNot(Receive())
				})
			})

			Context("when active_at is in the past", func() {
				It("should grab the job", func() {
					_, err := queue.Enqueue(&gobble.Job{
						ActiveAt: time.Now().Add(-1 * time.Minute),
					}, database.Connection)
					Expect(err).NotTo(HaveOccurred())
					Eventually(queue.Reserve("some-other-worker")).Should(Receive())
				})
			})
		})
	})

	Describe("Dequeue", func() {
		It("deletes the job from the queue", func() {
			job, err := queue.Enqueue(&gobble.Job{}, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))

			queue.Dequeue(job)

			results, err = database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(0))
		})

		It("ignores errors when the row is gone", func() {
			job, err := queue.Enqueue(&gobble.Job{}, database.Connection)
			Expect(err).NotTo(HaveOccurred())

			results, err := database.Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))

			Expect(func() {
				queue.Dequeue(job)
			}).NotTo(Panic())

			Expect(func() {
				queue.Dequeue(job)
			}).NotTo(Panic())
		})
	})

	Describe("Len", func() {
		It("returns the length of the queue", func() {
			job, err := queue.Enqueue(&gobble.Job{}, database.Connection)
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
})
