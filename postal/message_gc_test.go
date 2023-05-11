package postal_test

import (
	"bytes"
	"errors"
	"log"
	"time"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/postal"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageGC", func() {
	var (
		messageGC       postal.MessageGC
		repo            *mocks.MessagesRepo
		database        *mocks.Database
		conn            db.ConnectionInterface
		loggerBuffer    *bytes.Buffer
		lifetime        time.Duration
		pollingInterval time.Duration
	)

	BeforeEach(func() {
		loggerBuffer = bytes.NewBuffer([]byte{})
		logger := log.New(loggerBuffer, "", 0)

		conn = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = conn

		repo = mocks.NewMessagesRepo()

		lifetime = 2 * time.Minute
		pollingInterval = 500 * time.Millisecond

		messageGC = postal.NewMessageGC(lifetime, database, repo, pollingInterval, logger)
	})

	Describe("Run", func() {
		It("It calls collect every passed in duration", func() {
			messageGC.Run()

			Eventually(func() int {
				return repo.DeleteBeforeCall.CallCount
			}).Should(BeNumerically(">=", 2))

			call1 := repo.DeleteBeforeCall.InvocationTimes[0]
			call2 := repo.DeleteBeforeCall.InvocationTimes[1]
			Expect(call2).To(BeTemporally(">", call1.Add(pollingInterval-50*time.Millisecond)))
			Expect(call2).To(BeTemporally("<", call1.Add(pollingInterval+50*time.Millisecond)))
		})
	})

	Describe("Collect", func() {
		It("Deletes message statuses older than the specified time", func() {
			messageGC.Collect()

			Expect(repo.DeleteBeforeCall.Receives.Connection).To(Equal(conn))
			Expect(repo.DeleteBeforeCall.Receives.ThresholdTime).To(BeTemporally("~", time.Now().Add(-2*time.Minute), 10*time.Second))
		})

		Context("When the repo errors unexpectantly", func() {
			It("logs the error", func() {
				repo.DeleteBeforeCall.Returns.Error = errors.New("messages table is totally corrupt")

				messageGC.Collect()

				Expect(loggerBuffer.String()).To(ContainSubstring("messages table is totally corrupt"))
			})
		})

	})
})
