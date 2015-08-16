package senders_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebV2SendersSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v2/web/senders")
}
