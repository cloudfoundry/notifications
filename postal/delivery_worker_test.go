package postal_test

import (
	"bytes"
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/postal/common"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeliveryWorker", func() {
	var (
		worker                 postal.DeliveryWorker
		logger                 lager.Logger
		buffer                 *bytes.Buffer
		delivery               common.Delivery
		queue                  *mocks.Queue
		deliveryFailureHandler *mocks.DeliveryFailureHandler
		v1DeliveryJobProcessor *mocks.V1DeliveryJobProcessor
		connection             *mocks.Connection
		messageStatusUpdater   *mocks.MessageStatusUpdater
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))
		queue = mocks.NewQueue()
		deliveryFailureHandler = mocks.NewDeliveryFailureHandler()
		connection = mocks.NewConnection()
		database := mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = connection
		messageStatusUpdater = mocks.NewMessageStatusUpdater()

		config := postal.DeliveryWorkerConfig{
			ID:                     42,
			Logger:                 logger,
			Queue:                  queue,
			DeliveryFailureHandler: deliveryFailureHandler,
			Database:               database,
			UAAHost:                "my-uaa-host",
			MessageStatusUpdater:   messageStatusUpdater,
		}

		v1DeliveryJobProcessor = mocks.NewV1DeliveryJobProcessor()
		worker = postal.NewDeliveryWorker(v1DeliveryJobProcessor, config)
	})

	Describe("Work", func() {
		It("pops Deliveries off the queue, sending emails for each", func() {
			reserveChan := make(chan *gobble.Job)
			go func() {
				reserveChan <- gobble.NewJob(delivery)
			}()
			queue.ReserveCall.Returns.Chan = reserveChan

			worker.Work()

			<-time.After(10 * time.Millisecond)
			worker.Halt()

			Expect(v1DeliveryJobProcessor.ProcessCall.CallCount).To(Equal(1))
		})

		It("can be halted", func() {
			go func() {
				worker.Halt()
			}()

			Eventually(func() bool {
				worker.Work()
				return true
			}).Should(BeTrue())
		})
	})

	Describe("Deliver", func() {
		var job *gobble.Job

		BeforeEach(func() {
			job = gobble.NewJob(delivery)
		})

		It("should hand the job to the v1 workflow", func() {
			worker.Deliver(job)

			Expect(v1DeliveryJobProcessor.ProcessCall.Receives.Job).To(Equal(job))
			Expect(v1DeliveryJobProcessor.ProcessCall.Receives.Logger).ToNot(BeNil())
		})

		Context("when the job cannot be unmarshalled", func() {
			BeforeEach(func() {
				j := gobble.Job{
					Payload: "%%",
				}
				job = &j

				worker.Deliver(job)
			})

			It("should use the deliveryFailureHandler", func() {
				Expect(deliveryFailureHandler.HandleCall.WasCalled).To(BeTrue())
				Expect(deliveryFailureHandler.HandleCall.Receives.Job).ToNot(BeNil())
				Expect(deliveryFailureHandler.HandleCall.Receives.Logger).ToNot(BeNil())
			})
		})
	})
})
