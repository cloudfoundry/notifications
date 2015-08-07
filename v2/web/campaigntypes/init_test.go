package campaigntypes_test

import (
	"testing"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebV2CampaignTypesSuite(t *testing.T) {
	fakes.RegisterFastTokenSigningMethod()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Web V2 Campaign types Suite")
}
