package params_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebHandlersParamsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Params Suite")
}
