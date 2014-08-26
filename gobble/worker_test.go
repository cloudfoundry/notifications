package gobble_test

import (
    "github.com/cloudfoundry-incubator/notifications/gobble"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Worker", func() {
    var queue *gobble.Queue
    var worker gobble.Worker
    var callbackWasCalledWith gobble.Job

    BeforeEach(func() {
        TruncateTables()

        callback := func(job gobble.Job) {
            callbackWasCalledWith = job
        }

        queue = gobble.NewQueue()
        worker = gobble.NewWorker(1, queue, callback)
    })

    Describe("Perform", func() {
        It("reserves a job, performs the callback, and then dequeues the completed job", func() {
            job := queue.Enqueue(gobble.Job{
                Payload: "the-payload",
            })

            worker.Perform()

            Expect(callbackWasCalledWith.ID).To(Equal(job.ID))

            results, err := gobble.Database().Connection.Select(gobble.Job{}, "SELECT * FROM `jobs`")
            if err != nil {
                panic(err)
            }

            Expect(len(results)).To(Equal(0))
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
