package mocks

import "time"

type Clock struct {
	NowCall struct {
		Returns struct {
			Time time.Time
		}
	}
}

func NewClock() *Clock {
	return &Clock{}
}

func (c *Clock) Now() time.Time {
	return c.NowCall.Returns.Time
}
