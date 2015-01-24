package rainmaker_test

import (
	"testing"

	"github.com/pivotal-golang/rainmaker/internal/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var fakeCloudController *fakes.CloudController

func TestRainmakerSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rainmaker Suite")
}

var _ = BeforeSuite(func() {
	fakeCloudController = fakes.NewCloudController()
	fakeCloudController.Start()
})

var _ = AfterSuite(func() {
	fakeCloudController.Close()
})

var _ = BeforeEach(func() {
	fakeCloudController.Reset()
})
