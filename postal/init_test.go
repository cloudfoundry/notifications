package postal_test

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPostalSuite(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	metricsLogger := metrics.DefaultLogger
	metrics.DefaultLogger = log.New(buffer, "", 0)

	RegisterFailHandler(Fail)
	RunSpecs(t, "postal")

	metrics.DefaultLogger = metricsLogger
}

type logLine struct {
	Source   string                 `json:"source"`
	Message  string                 `json:"message"`
	LogLevel int                    `json:"log_level"`
	Data     map[string]interface{} `json:"data"`
}

func parseLogLines(b []byte) ([]logLine, error) {
	var lines []logLine
	for _, line := range bytes.Split(b, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		var ll logLine
		err := json.Unmarshal(line, &ll)
		if err != nil {
			return lines, err
		}

		lines = append(lines, ll)
	}

	return lines, nil
}
