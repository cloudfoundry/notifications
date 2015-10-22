package gobble

import "time"

type TickerInterface interface {
	Tick() <-chan time.Time
	Start()
	Stop()
}

type Ticker struct {
	constructor func(time.Duration) *time.Ticker
	duration    time.Duration
	ticker      *time.Ticker
}

func NewTicker(tickerConstructor func(time.Duration) *time.Ticker, duration time.Duration) *Ticker {
	return &Ticker{
		constructor: tickerConstructor,
		duration:    duration,
	}
}

func (t Ticker) Tick() <-chan time.Time {
	if t.ticker == nil {
		return nil
	}

	return t.ticker.C
}

func (t *Ticker) Start() {
	t.ticker = t.constructor(t.duration)
}

func (t Ticker) Stop() {
	t.ticker.Stop()
}
