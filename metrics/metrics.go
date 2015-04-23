package metrics

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var DefaultLogger = NewLogger(os.Stdout)

type Metric struct {
	Kind    string                 `json:"kind"`
	Payload map[string]interface{} `json:"payload"`
}

func NewLogger(buffer io.Writer) *log.Logger {
	return log.New(buffer, "[METRIC] ", 0)
}

func NewMetric(kind string, payload map[string]interface{}) Metric {
	return Metric{
		Kind:    kind,
		Payload: payload,
	}
}

func (metric Metric) Log() {
	metric.LogWith(DefaultLogger)
}

func (metric Metric) LogWith(logger *log.Logger) {
	message, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}

	logger.Printf("%s", message)
}
