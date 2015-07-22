package notificationtypes_test

import (
	"testing"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWebV2NotificationTypesSuite(t *testing.T) {
	fakes.RegisterFastTokenSigningMethod()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Web V2 Notification Types Suite")
}
