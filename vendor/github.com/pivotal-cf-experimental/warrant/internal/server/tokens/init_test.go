package tokens_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTokensSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/server/tokens")
}
