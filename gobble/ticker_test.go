package gobble_test

import (
	"time"

	"github.com/cloudfoundry-incubator/notifications/gobble"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ticker", func() {
	var (
		ticker           *gobble.Ticker
		timeChan         chan time.Time
		recordedDuration time.Duration
	)

	BeforeEach(func() {
		recordedDuration = time.Duration(0)
		timeChan = make(chan time.Time)

		constructor := func(duration time.Duration) *time.Ticker {
			recordedDuration = duration
			return &time.Ticker{C: timeChan}
		}
		ticker = gobble.NewTicker(constructor, 1*time.Millisecond)
	})

	Describe("Start", func() {
		It("calls the constructor with the interval", func() {
			Expect(recordedDuration).To(Equal(time.Duration(0)))

			ticker.Start()

			Expect(recordedDuration).To(Equal(1 * time.Millisecond))
		})
	})

	Describe("Tick", func() {
		It("returns a channel if the ticker has been started", func() {
			var receiveTimeChan <-chan time.Time
			receiveTimeChan = timeChan

			ticker.Start()
			Expect(ticker.Tick()).To(Equal(receiveTimeChan))
		})

		It("returns nil if the ticker has not been started", func() {
			Expect(ticker.Tick()).To(BeNil())
		})
	})

	Describe("Stop", func() {
		It("stops the ticker", func() {
			ticker = gobble.NewTicker(time.NewTicker, 1*time.Millisecond)

			ticker.Start()
			ticker.Stop()

			Consistently(ticker.Tick(), 50*time.Millisecond).ShouldNot(Receive())
		})
	})
})
