package gobble_test

import (
    "github.com/cloudfoundry-incubator/notifications/gobble"

    . "github.com/onsi/ginkgo"
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

            Expect(job).To(BeAssignableToTypeOf(gobble.Job{}))
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
            if err != nil {
                panic(err)
            }

            Expect(payload).To(Equal(data))
        })
    })
})
