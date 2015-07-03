package acceptance

import (
	"database/sql"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/acceptance/servers"
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/gobble"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/pivotal-cf/uaa-sso-golang/uaa"

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
	UAA           *servers.UAA
	ZonedUAA      *servers.UAA
}

func TestAcceptanceSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

var _ = BeforeSuite(func() {
	os.Setenv("VCAP_APPLICATION", `{"instance_index": -1}`)

	MigrateDatabase()

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
	ResetDatabase()
	Servers.SMTP.Reset()
})

func FetchDatabases() (*models.DB, *gobble.DB) {
	env := application.NewEnvironment()
	sqlDB, err := sql.Open("mysql", env.DatabaseURL)
	Expect(err).NotTo(HaveOccurred())

	database := models.NewDatabase(sqlDB, models.Config{DefaultTemplatePath: path.Join(env.RootPath, "templates", "default.json")})
	gobbleDB := gobble.NewDatabase(sqlDB)

	return database, gobbleDB
}

func MigrateDatabase() {
	env := application.NewEnvironment()
	database, gobbleDB := FetchDatabases()

	database.Migrate(env.ModelMigrationsPath)
	gobbleDB.Migrate(env.GobbleMigrationsPath)
}

func ResetDatabase() {
	database, gobbleDB := FetchDatabases()

	database.Setup()
	database.Connection().(*models.Connection).TruncateTables()
	database.Seed()

	gobbleDB.Connection.TruncateTables()
}

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
