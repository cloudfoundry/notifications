package preferences_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWebV1PreferencesSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v1/web/preferences")
}
