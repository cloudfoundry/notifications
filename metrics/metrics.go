package metrics

import (
    "encoding/json"
    "log"
    "os"
)

var Logger = log.New(os.Stdout, "", 0)

type Metric struct {
    kind    string
    payload map[string]interface{}
}

func NewMetric(kind string, payload map[string]interface{}) Metric {
    return Metric{
        kind:    kind,
        payload: payload,
    }
}

func (metric Metric) Log() {
    payload, err := json.Marshal(metric.payload)
    if err != nil {
        panic(err)
    }
    Logger.Printf("[METRIC] {\"kind\":\"%s\",\"payload\":\"%s\"}", metric.kind, payload)
}
