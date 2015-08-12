package testing

import (
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/models"
	v1models "github.com/cloudfoundry-incubator/notifications/v1/models"
)

func TruncateTables(database *db.DB) {
	env := application.NewEnvironment()
	dbMigrator := v1models.DatabaseMigrator{}
	dbMigrator.Migrate(database.RawConnection(), env.ModelMigrationsPath)
	models.Setup(database)

	connection := database.Connection().(*db.Connection)
	err := connection.TruncateTables()
	if err != nil {
		panic(err)
	}
}
