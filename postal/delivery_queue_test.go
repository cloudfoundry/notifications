package postal_test

import (
    "time"

    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("DeliveryQueue", func() {
    var queue *postal.DeliveryQueue
    var delivery postal.Delivery

    BeforeEach(func() {
        queue = postal.NewDeliveryQueue()
        delivery = postal.Delivery{
            User:         uaa.User{},
            Options:      postal.Options{},
            UserGUID:     "user-guid",
            Space:        "the-space",
            Organization: "the-organization",
            ClientID:     "client-id",
            Templates:    postal.Templates{},
        }
    })

    Describe("Enqueue/Dequeue", func() {
        It("enqueues delivery jobs onto the queue", func() {
            queue.Enqueue(delivery)

            Expect(<-queue.Dequeue()).To(Equal(delivery))
        })

        It("waits for delivery jobs to be enqueued", func() {
            go func() {
                <-time.After(100 * time.Millisecond)
                queue.Enqueue(delivery)
            }()

            Expect(<-queue.Dequeue()).To(Equal(delivery))
        })
    })
})
