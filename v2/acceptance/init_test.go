package acceptance

import (
	"fmt"
	"io/ioutil"
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

	roundtripRecorder *docs.RoundTripRecorder
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

	roundtripRecorder = docs.NewRoundTripRecorder()
})

var _ = AfterSuite(func() {
	Servers.Notifications.Close()
	Servers.Notifications.Destroy()

	if roundtripRecorder != nil {
		context, err := docs.BuildTemplateContext(docs.Structure, roundtripRecorder.RoundTrips)
		Expect(err).NotTo(HaveOccurred())

		markdown, err := docs.GenerateMarkdown(context)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(fmt.Sprintf("%s/V2_API.md", os.Getenv("ROOT_PATH")), []byte(markdown), 0644)
		Expect(err).NotTo(HaveOccurred())
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
