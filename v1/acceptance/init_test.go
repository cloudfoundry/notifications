package v1

import (
	"os"
	"regexp"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/testing/servers"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	GUIDRegex = regexp.MustCompile(`[0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12}`)
	Servers   struct {
		Notifications servers.Notifications
		SMTP          *servers.SMTP
		CC            servers.CC
		UAA           *servers.UAA
	}
)

func TestAcceptanceSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v1/acceptance")
}

var _ = BeforeSuite(func() {
	os.Setenv("VCAP_APPLICATION", `{"instance_index": -1}`)

	Servers.SMTP = servers.NewSMTP()
	Servers.SMTP.Boot()

	Servers.UAA = servers.NewUAA()
	Servers.UAA.Boot()

	users := map[string]string{
		"user-123":          "user-123",
		"user-456":          "user-456",
		"user-789":          "user-789",
		"unauthorized-user": "unauthorized-user",
	}

	Servers.CC = servers.NewCC(users)
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
	Servers.SMTP.Close()
})

var _ = BeforeEach(func() {
	Servers.Notifications.ResetDatabase()
	Servers.SMTP.Reset()
})

var _ = AfterEach(func() {
	err := Servers.Notifications.WaitForJobsQueueToEmpty()
	Expect(err).NotTo(HaveOccurred())
})

func GetClientTokenFor(clientID string) uaa.Token {
	token, err := GetUAAClientFor(clientID).GetClientToken()
	if err != nil {
		panic(err)
	}

	return token
}

func GetUserTokenFor(code string) uaa.Token {
	token, err := GetUAAClientFor("notifications-sender").Exchange(code)
	if err != nil {
		panic(err)
	}

	return token
}

func GetUAAClientFor(clientID string) uaa.UAA {
	return uaa.NewUAA("", Servers.UAA.ServerURL, clientID, "secret", "")
}
