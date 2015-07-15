package v1

import (
	"os"
	"regexp"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	TRUE      = true
	FALSE     = false
	GUIDRegex = regexp.MustCompile(`[0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12}`)
	Servers   struct {
		Notifications servers.Notifications
		SMTP          *servers.SMTP
		CC            servers.CC
		UAA           *servers.UAA
		ZonedUAA      *servers.UAA
	}
)

func TestAcceptanceSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V1 Acceptance Suite")
}

var _ = BeforeSuite(func() {
	os.Setenv("VCAP_APPLICATION", `{"instance_index": -1}`)

	Servers.SMTP = servers.NewSMTP()
	Servers.SMTP.Boot()

	Servers.ZonedUAA = servers.NewUAA("testzone1")
	Servers.ZonedUAA.Boot()

	Servers.UAA = servers.NewUAA("uaa")
	Servers.UAA.Boot()

	Servers.CC = servers.NewCC()
	Servers.CC.Boot()

	Servers.Notifications = servers.NewNotifications()
	Servers.Notifications.Compile()
	Servers.Notifications.MigrateDatabase()
	Servers.Notifications.Boot()
})

var _ = AfterSuite(func() {
	Servers.Notifications.Close()
	Servers.Notifications.Destroy()
	Servers.CC.Close()
	Servers.UAA.Close()
	Servers.ZonedUAA.Close()
	Servers.SMTP.Close()
})

var _ = BeforeEach(func() {
	Servers.Notifications.ResetDatabase()
	Servers.SMTP.Reset()
})

func GetClientTokenFor(clientID, zone string) uaa.Token {
	token, err := GetUAAClientFor(clientID, zone).GetClientToken()
	if err != nil {
		panic(err)
	}

	return token
}

func GetUserTokenFor(code string) uaa.Token {
	token, err := GetUAAClientFor("notifications-sender", "uaa").Exchange(code)
	if err != nil {
		panic(err)
	}

	return token
}

func GetUAAClientFor(clientID string, zone string) uaa.UAA {
	var host string
	if zone == "testzone1" {
		host = Servers.ZonedUAA.ServerURL
	} else {
		host = Servers.UAA.ServerURL
	}
	return uaa.NewUAA("", host, clientID, "secret", "")
}
