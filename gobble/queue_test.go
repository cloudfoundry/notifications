package gobble_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queue", func() {
	var queue *gobble.Queue
	var waitMaxDuration time.Duration

	BeforeEach(func() {
		TruncateTables()

		waitMaxDuration = gobble.WaitMaxDuration
		gobble.WaitMaxDuration = 100 * time.Millisecond
		queue = gobble.NewQueue()
	})

	AfterEach(func() {
		gobble.WaitMaxDuration = waitMaxDuration
		queue.Close()
	})

	Describe("Enqueue", func() {
		It("sticks the job in the database table", func() {
			job := gobble.NewJob(map[string]bool{
				"testing": true,
			})

			job, err := queue.Enqueue(job)
			if err != nil {
				panic(err)
			}

			results, err := gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			if err != nil {
				panic(err)
			}

			jobs := []gobble.Job{}
			for _, result := range results {
				jobs = append(jobs, *(result.(*gobble.Job)))
			}

			Expect(len(jobs)).To(Equal(1))
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
			err = gobble.Database().Connection.SelectOne(&reloadedJob, "SELECT * FROM `jobs` where id = ?", job.ID)
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

			err := gobble.Database().Connection.Insert(&job)
			if err != nil {
				panic(err)
			}

			jobChannel := queue.Reserve("workerId")
			reservedJob := <-jobChannel

			Expect(reservedJob.ID).To(Equal(job.ID))
		})

		It("reserves the next available job", func() {
			job1 := gobble.Job{
				Payload: "first",
			}
			job2 := gobble.Job{
				Payload: "second",
			}

			job1, err := queue.Enqueue(job1)
			if err != nil {
				panic(err)
			}

			job2, err = queue.Enqueue(job1)
			if err != nil {
				panic(err)
			}

			jobChannel := queue.Reserve("1")
			var reservedJob1 gobble.Job
			Eventually(jobChannel).Should(Receive(&reservedJob1))

			jobChannel = queue.Reserve("1")
			var reservedJob2 gobble.Job
			Eventually(jobChannel).Should(Receive(&reservedJob2))

			Expect(reservedJob1.ID).To(Equal(job1.ID))
			Expect(reservedJob2.ID).To(Equal(job2.ID))
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

			Eventually(done, 5*time.Second).Should(Receive())
			Eventually(done, 5*time.Second).Should(Receive())

			results, err := gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs` WHERE `worker_id` = ''")
			if err != nil {
				panic(err)
			}

			Expect(len(results)).To(Equal(0))
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
	})

	Describe("Dequeue", func() {
		It("deletes the job from the queue", func() {
			job, err := queue.Enqueue(gobble.Job{})
			if err != nil {
				panic(err)
			}
			results, err := gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			if err != nil {
				panic(err)
			}
			Expect(len(results)).To(Equal(1))

			queue.Dequeue(job)
			results, err = gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
			if err != nil {
				panic(err)
			}
			Expect(len(results)).To(Equal(0))
		})
	})

	Describe("Unlock", func() {
		It("clears the workerID values for any jobs in the queue", func() {
			queue.Enqueue(gobble.Job{})
			results, err := gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs` WHERE `worker_id` = ''")
			if err != nil {
				panic(err)
			}
			Expect(len(results)).To(Equal(1))

			<-queue.Reserve("my-worker")
			results, err = gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs` WHERE `worker_id` = ''")
			if err != nil {
				panic(err)
			}
			Expect(len(results)).To(Equal(0))

			queue.Unlock()
			results, err = gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs` WHERE `worker_id` = ''")
			if err != nil {
				panic(err)
			}
			Expect(len(results)).To(Equal(1))
		})
	})
})
