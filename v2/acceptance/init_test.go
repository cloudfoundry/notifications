package acceptance

import (
	"os"
	"strconv"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/docs"
	"github.com/cloudfoundry-incubator/notifications/testing/servers"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	Servers struct {
		Notifications servers.Notifications
		SMTP          *servers.SMTP
		UAA           *servers.UAA
		CC            servers.CC
	}
	Trace, _ = strconv.ParseBool(os.Getenv("TRACE"))

	docCollection *docs.DocGenerator
)

func TestAcceptanceSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v2/acceptance")
}

var _ = BeforeSuite(func() {
	Servers.SMTP = servers.NewSMTP()
	Servers.SMTP.Boot()

	Servers.UAA = servers.NewUAA()
	Servers.UAA.Boot()

	Servers.CC = servers.NewCC()
	Servers.CC.Boot()

	Servers.Notifications = servers.NewNotifications()
	Servers.Notifications.Compile()
	Servers.Notifications.ResetDatabase()
	Servers.Notifications.Boot()

	docCollection = docs.NewDocGenerator(docs.NewRequestInspector())
})

var _ = AfterSuite(func() {
	Servers.Notifications.Close()
	Servers.Notifications.Destroy()

	if docCollection != nil {
		docCollection.GenerateBlueprint(os.Getenv("DOC_FILE"))
	}
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
