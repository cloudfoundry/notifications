package acceptance

import (
	"path"
	"regexp"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	TRUE  = true
	FALSE = false
)

var GUIDRegex = regexp.MustCompile(`[0-9a-f]{8}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{4}\-[0-9a-f]{12}`)

var Servers struct {
	Notifications servers.Notifications
	SMTP          *servers.SMTP
	CC            servers.CC
	UAA           servers.UAA
}

func TestAcceptanceSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
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
	TruncateTables()
	Servers.SMTP.Reset()
})

func TruncateTables() {
	env := application.NewEnvironment()
	migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
	models.NewDatabase(env.DatabaseURL, migrationsPath).Connection().(*models.Connection).TruncateTables()
	gobble.Database().Connection.TruncateTables()
}
