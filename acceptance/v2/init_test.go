package v2

import (
	"testing"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

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

func GetClientTokenFor(clientID, zone string) uaa.Token {
	token, err := GetUAAClientFor(clientID, zone).GetClientToken()
	if err != nil {
		panic(err)
	}

	return token
}

func GetUAAClientFor(clientID string, zone string) uaa.UAA {
	return uaa.NewUAA("", Servers.UAA.ServerURL, clientID, "secret", "")
}
