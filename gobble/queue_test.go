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
    })

    Describe("Enqueue", func() {
        It("sticks the job in the database table", func() {
            job := gobble.NewJob(map[string]bool{
                "testing": true,
            })

            job = queue.Enqueue(job)

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

            job1 = queue.Enqueue(job1)
            job2 = queue.Enqueue(job1)

            jobChannel := queue.Reserve("1")
            reservedJob1 := <-jobChannel

            jobChannel = queue.Reserve("1")
            reservedJob2 := <-jobChannel

            Expect(reservedJob1.ID).To(Equal(job1.ID))
            Expect(reservedJob2.ID).To(Equal(job2.ID))
        })

        It("keeps trying to reserve a job until one becomes available", func() {
            jobChannel := queue.Reserve("my-id")

            Consistently(jobChannel).ShouldNot(Receive())

            job := queue.Enqueue(gobble.Job{
                Payload: "hello",
            })

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

            <-done
            <-done

            results, err := gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs` WHERE `worker_id` = ''")
            if err != nil {
                panic(err)
            }

            Expect(len(results)).To(Equal(0))
        })
    })

    Describe("Dequeue", func() {
        It("deletes the job from the queue", func() {
            job := queue.Enqueue(gobble.Job{})
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
})
