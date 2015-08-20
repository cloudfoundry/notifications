package queue_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQueueSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v2/queue")
}
