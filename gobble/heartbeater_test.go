package gobble_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MockTicker struct {
	TickCall struct {
		CallCount int
		Returns   struct {
			TimeChan <-chan time.Time
		}
	}

	StartCall struct {
		CallCount int
	}

	StopCall struct {
		WasCalled bool
	}
}

func (t *MockTicker) Tick() <-chan time.Time {
	t.TickCall.CallCount++
	return t.TickCall.Returns.TimeChan
}

func (t *MockTicker) Start() {
	t.StartCall.CallCount++
}

func (t *MockTicker) Stop() {
	t.StopCall.WasCalled = true
}

var _ = Describe("Heartbeater", func() {
	var (
		queue  *mocks.Queue
		ticker *MockTicker
		beater gobble.Heartbeater
	)

	BeforeEach(func() {
		queue = mocks.NewQueue()
		ticker = &MockTicker{}
		beater = gobble.NewHeartbeater(queue, ticker)
	})

	Describe("Beat", func() {
		It("updates the active_at lease for a job on a given interval", func() {
			timeChan := make(chan time.Time)
			ticker.TickCall.Returns.TimeChan = timeChan
			job := &gobble.Job{}
			Expect(queue.RequeueCall.Receives.Job).To(BeNil())

			go beater.Beat(job)

			Eventually(func() int {
				return ticker.StartCall.CallCount
			}).Should(Equal(1))

			now := time.Now()
			timeChan <- now

			Eventually(func() *gobble.Job {
				return queue.RequeueCall.Receives.Job
			}).Should(Equal(&gobble.Job{
				ActiveAt: now,
			}))

			futureTime := time.Now().Add(10 * time.Minute)
			timeChan <- futureTime

			Eventually(func() *gobble.Job {
				return queue.RequeueCall.Receives.Job
			}).Should(Equal(&gobble.Job{
				ActiveAt: futureTime,
			}))

			Expect(ticker.StartCall.CallCount).To(Equal(1))
		})
	})

	Describe("Halt", func() {
		It("stops beating", func() {
			job := &gobble.Job{}

			go beater.Beat(job)

			Expect(ticker.StopCall.WasCalled).To(BeFalse())
			Eventually(func() int {
				return ticker.TickCall.CallCount
			}).Should(Equal(1))

			beater.Halt()

			Eventually(func() bool {
				return ticker.StopCall.WasCalled
			}).Should(BeTrue())

			Expect(ticker.TickCall.CallCount).To(Equal(1))
		})
	})
})
