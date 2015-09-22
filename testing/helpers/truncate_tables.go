package helpers

import (
	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/db"
	v1models "github.com/cloudfoundry-incubator/notifications/v1/models"
	v2models "github.com/cloudfoundry-incubator/notifications/v2/models"
)

func TruncateTables(database *db.DB) {
	env, err := application.NewEnvironment()
	if err != nil {
		panic(err)
	}

	dbMigrator := v1models.DatabaseMigrator{}
	dbMigrator.Migrate(database.RawConnection(), env.ModelMigrationsPath)
	v1models.Setup(database)
	v2models.Setup(database)

	connection := database.Connection().(*db.Connection)
	err = connection.TruncateTables()
	if err != nil {
		panic(err)
	}
}
