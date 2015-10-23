package mocks

type MetricsEmitter struct {
	IncrementCall struct {
		Receives struct {
			Counter string
		}
	}
}

func NewMetricsEmitter() *MetricsEmitter {
	return &MetricsEmitter{}
}

func (e *MetricsEmitter) Increment(counter string) {
	e.IncrementCall.Receives.Counter = counter
}
