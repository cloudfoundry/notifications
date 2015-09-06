package postal_test

import (
	"bytes"
	"crypto/md5"
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/uaa"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/pivotal-golang/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeliveryWorker", func() {
	var (
		mailClient             *mocks.MailClient
		worker                 postal.DeliveryWorker
		id                     int
		logger                 lager.Logger
		buffer                 *bytes.Buffer
		delivery               postal.Delivery
		queue                  *mocks.Queue
		unsubscribesRepo       *mocks.UnsubscribesRepo
		globalUnsubscribesRepo *mocks.GlobalUnsubscribesRepo
		kindsRepo              *mocks.KindsRepo
		database               *mocks.Database
		strategyDeterminer     *mocks.StrategyDeterminer
		conn                   *mocks.Connection
		userLoader             *mocks.UserLoader
		userGUID               string
		fakeUserEmail          string
		templateLoader         *mocks.TemplatesLoader
		receiptsRepo           *mocks.ReceiptsRepo
		tokenLoader            *mocks.TokenLoader
		messageID              string
		messageStatusUpdater   *mocks.MessageStatusUpdater
		deliveryFailureHandler *mocks.DeliveryFailureHandler
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		id = 1234
		logger = lager.NewLogger("notifications")
		logger.RegisterSink(lager.NewWriterSink(buffer, lager.DEBUG))
		mailClient = mocks.NewMailClient()
		queue = mocks.NewQueue()
		unsubscribesRepo = mocks.NewUnsubscribesRepo()
		globalUnsubscribesRepo = mocks.NewGlobalUnsubscribesRepo()
		kindsRepo = mocks.NewKindsRepo()
		kindsRepo.FindCall.Returns.Kinds = []models.Kind{
			{
				ID:       "some-kind",
				ClientID: "some-client",
				Critical: false,
			},
			{
				ID:       "another-kind",
				ClientID: "some-client",
				Critical: false,
			},
		}

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		strategyDeterminer = mocks.NewStrategyDeterminer()
		userGUID = "user-123"
		sum := md5.Sum([]byte("banana's are so very tasty"))
		encryptionKey := sum[:]
		fakeUserEmail = "user-123@example.com"
		userLoader = mocks.NewUserLoader()
		userLoader.LoadCall.Returns.Users = map[string]uaa.User{
			"user-123": {Emails: []string{fakeUserEmail}},
			"user-456": {Emails: []string{"user-456@example.com"}},
		}
		tokenLoader = mocks.NewTokenLoader()
		templateLoader = mocks.NewTemplatesLoader()
		templateLoader.LoadTemplatesCall.Returns.Templates = postal.Templates{
			Text:    "{{.Text}} {{.Domain}}",
			HTML:    "<p>{{.HTML}}</p>",
			Subject: "{{.Subject}}",
		}
		receiptsRepo = mocks.NewReceiptsRepo()
		messageStatusUpdater = mocks.NewMessageStatusUpdater()
		deliveryFailureHandler = mocks.NewDeliveryFailureHandler()

		worker = postal.NewDeliveryWorker(postal.DeliveryWorkerConfig{
			ID:            id,
			Sender:        "from@example.com",
			Domain:        "example.com",
			UAAHost:       "https://uaa.example.com",
			EncryptionKey: encryptionKey,
			Logger:        logger,
			Queue:         queue,

			Database:               database,
			DBTrace:                false,
			GlobalUnsubscribesRepo: globalUnsubscribesRepo,
			UnsubscribesRepo:       unsubscribesRepo,
			KindsRepo:              kindsRepo,
			UserLoader:             userLoader,
			TemplatesLoader:        templateLoader,
			ReceiptsRepo:           receiptsRepo,
			TokenLoader:            tokenLoader,
			MailClient:             mailClient,
			StrategyDeterminer:     strategyDeterminer,
			MessageStatusUpdater:   messageStatusUpdater,
			DeliveryFailureHandler: deliveryFailureHandler,
		})

		messageID = "randomly-generated-guid"
		delivery = postal.Delivery{
			ClientID: "some-client",
			UserGUID: userGUID,
			Options: postal.Options{
				Subject:    "the subject",
				Text:       "body content",
				ReplyTo:    "thesender@example.com",
				KindID:     "some-kind",
				TemplateID: "some-template-id",
			},
			MessageID:     messageID,
			VCAPRequestID: "some-request-id",
		}
	})

	Describe("Work", func() {
		It("pops Deliveries off the queue, sending emails for each", func() {
			queue.Enqueue(gobble.NewJob(delivery))

			worker.Work()

			delivery2 := postal.Delivery{
				UserGUID: "user-456",
			}
			queue.Enqueue(gobble.NewJob(delivery2))

			<-time.After(10 * time.Millisecond)
			worker.Halt()

			Expect(mailClient.SendCall.CallCount).To(Equal(2))
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
})
