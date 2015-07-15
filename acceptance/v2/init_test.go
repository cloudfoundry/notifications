package v2

import (
	"testing"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var Servers struct {
	Notifications servers.Notifications
	SMTP          *servers.SMTP
	UAA           *servers.UAA
}

func TestAcceptanceSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V2 Acceptance Suite")
}

var _ = BeforeSuite(func() {
	Servers.SMTP = servers.NewSMTP()
	Servers.SMTP.Boot()

	Servers.UAA = servers.NewUAA("uaa")
	Servers.UAA.Boot()

	Servers.Notifications = servers.NewNotifications()
	Servers.Notifications.Compile()
	Servers.Notifications.Boot()
})

var _ = AfterSuite(func() {
	Servers.Notifications.Close()
	Servers.Notifications.Destroy()
})
