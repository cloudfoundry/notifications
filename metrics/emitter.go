package metrics

import "log"

type Emitter struct {
	logger *log.Logger
}

func NewEmitter(logger *log.Logger) Emitter {
	return Emitter{
		logger: logger,
	}
}

func (e Emitter) Increment(counter string) {
	NewMetric("counter", map[string]interface{}{
		"name": counter,
	}).LogWith(e.logger)
}
