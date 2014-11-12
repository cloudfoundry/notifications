package metrics

import (
	"encoding/json"
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "", 0)

type Metric struct {
	Kind    string                 `json:"kind"`
	Payload map[string]interface{} `json:"payload"`
}

func NewMetric(kind string, payload map[string]interface{}) Metric {
	return Metric{
		Kind:    kind,
		Payload: payload,
	}
}

func (metric Metric) Log() {
	message, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}
	Logger.Printf("[METRIC] %s", message)
}
