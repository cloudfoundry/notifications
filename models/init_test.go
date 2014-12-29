package models_test

import (
	"path"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestModelsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

func TruncateTables() {
	env := application.NewEnvironment()
	migrationsPath := path.Join(env.RootPath, env.ModelMigrationsDir)
	db := models.NewDatabase(models.Config{
		DatabaseURL:    env.DatabaseURL,
		MigrationsPath: migrationsPath,
	})
	connection := db.Connection().(*models.Connection)
	err := connection.TruncateTables()
	if err != nil {
		panic(err)
	}
}
