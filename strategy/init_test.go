package strategy_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStrategySuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "strategy")
}
