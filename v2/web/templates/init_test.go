package templates_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebV2TemplatesSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v2/web/templates")
}
