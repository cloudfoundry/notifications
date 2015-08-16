package campaigntypes_test

import (
	"testing"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebV2CampaignTypesSuite(t *testing.T) {
	fakes.RegisterFastTokenSigningMethod()

	RegisterFailHandler(Fail)
	RunSpecs(t, "v2/web/campaigntypes")
}
