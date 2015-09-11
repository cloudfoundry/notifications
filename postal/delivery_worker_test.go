package postal_test

import (
	"bytes"
	"errors"
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeliveryWorker", func() {
	var (
		worker                 postal.DeliveryWorker
		logger                 lager.Logger
		buffer                 *bytes.Buffer
		delivery               postal.Delivery
		queue                  *mocks.Queue
		deliveryFailureHandler *mocks.DeliveryFailureHandler
		v1Workflow             *mocks.Process
		v2Workflow             *mocks.Workflow
		strategyDeterminer     *mocks.StrategyDeterminer
		connection             *mocks.Connection
		messageStatusUpdater   *mocks.MessageStatusUpdater
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))
		queue = mocks.NewQueue()
		deliveryFailureHandler = mocks.NewDeliveryFailureHandler()
		strategyDeterminer = mocks.NewStrategyDeterminer()
		connection = mocks.NewConnection()
		database := mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = connection
		messageStatusUpdater = mocks.NewMessageStatusUpdater()

		config := postal.DeliveryWorkerConfig{
			Logger: logger,
			Queue:  queue,
			DeliveryFailureHandler: deliveryFailureHandler,
			StrategyDeterminer:     strategyDeterminer,
			Database:               database,
			UAAHost:                "my-uaa-host",
			MessageStatusUpdater:   messageStatusUpdater,
		}

		v1Workflow = mocks.NewProcess()
		v2Workflow = mocks.NewWorkflow()
		worker = postal.NewDeliveryWorker(v1Workflow, v2Workflow, config)
	})

	Describe("Work", func() {
		It("pops Deliveries off the queue, sending emails for each", func() {
			reserveChan := make(chan gobble.Job)
			go func() {
				reserveChan <- gobble.NewJob(delivery)
			}()
			queue.ReserveCall.Returns.Chan = reserveChan

			worker.Work()

			<-time.After(10 * time.Millisecond)
			worker.Halt()

			Expect(v1Workflow.DeliverCall.CallCount).To(Equal(1))
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
			j := gobble.NewJob(delivery)
			job = &j
		})

		Context("when the job is not a campaign, and not a v2 delivery", func() {
			It("should hand the job to the v1 workflow", func() {
				worker.Deliver(job)

				Expect(v1Workflow.DeliverCall.Receives.Job).To(Equal(job))
				Expect(v1Workflow.DeliverCall.Receives.Logger).ToNot(BeNil())
			})
		})

		Context("when the job is a campaign", func() {
			BeforeEach(func() {
				j := gobble.NewJob(struct {
					JobType string
				}{
					JobType: "campaign",
				})
				job = &j
			})

			It("uses the strategy determiner", func() {
				worker.Deliver(job)

				Expect(strategyDeterminer.DetermineCall.Receives.Job).To(Equal(*job))
				Expect(strategyDeterminer.DetermineCall.Receives.UAAHost).To(Equal("my-uaa-host"))
				Expect(strategyDeterminer.DetermineCall.Receives.Connection).To(Equal(connection))
			})

			Context("when the strategy fails to determine", func() {
				It("uses the deliveryFailureHandler", func() {
					strategyDeterminer.DetermineCall.Returns.Error = errors.New("some error")

					worker.Deliver(job)

					Expect(deliveryFailureHandler.HandleCall.WasCalled).To(BeTrue())
					Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
					Expect(deliveryFailureHandler.HandleCall.Receives.Logger).ToNot(BeNil())
				})
			})
		})

		Context("when the job is a v2 workflow", func() {
			BeforeEach(func() {
				j := gobble.NewJob(struct {
					JobType    string
					MessageID  string
					CampaignID string
				}{
					JobType:    "v2",
					MessageID:  "some-message-id",
					CampaignID: "some-campaign-id",
				})
				job = &j
			})

			It("should hand the job to the v2 workflow", func() {
				worker.Deliver(job)

				Expect(v2Workflow.DeliverCall.Receives.Delivery).To(Equal(postal.Delivery{
					MessageID:  "some-message-id",
					CampaignID: "some-campaign-id",
				}))
				Expect(v2Workflow.DeliverCall.Receives.Logger).NotTo(BeNil())
				Expect(v1Workflow.DeliverCall.CallCount).To(Equal(0))
			})

			Context("when the workflow encounters an error", func() {
				It("updates the message status to retry if the job should be retried", func() {
					v2Workflow.DeliverCall.Returns.Error = errors.New("delivery failure")
					job.ShouldRetry = true

					worker.Deliver(job)

					Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
					Expect(deliveryFailureHandler.HandleCall.Receives.Logger).NotTo(BeNil())
					Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(connection))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal("some-message-id"))
					Expect(messageStatusUpdater.UpdateCall.Receives.CampaignID).To(Equal("some-campaign-id"))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal("retry"))
				})

				It("updates the message status to failed if the job should not be retried", func() {
					v2Workflow.DeliverCall.Returns.Error = errors.New("delivery failure")
					job.ShouldRetry = false

					worker.Deliver(job)

					Expect(deliveryFailureHandler.HandleCall.Receives.Job).To(Equal(job))
					Expect(deliveryFailureHandler.HandleCall.Receives.Logger).NotTo(BeNil())
					Expect(messageStatusUpdater.UpdateCall.Receives.Connection).To(Equal(connection))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageID).To(Equal("some-message-id"))
					Expect(messageStatusUpdater.UpdateCall.Receives.CampaignID).To(Equal("some-campaign-id"))
					Expect(messageStatusUpdater.UpdateCall.Receives.MessageStatus).To(Equal("failed"))
				})
			})
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
