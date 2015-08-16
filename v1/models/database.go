package models

import (
	"database/sql"

	"github.com/cloudfoundry-incubator/notifications/db"
)

type DatabaseInterface interface {
	db.DatabaseInterface
}

type ConnectionInterface interface {
	db.ConnectionInterface
}

type Config struct {
	DefaultTemplatePath string
}

func NewDatabase(sqlDB *sql.DB, config Config) DatabaseInterface {
	database := db.NewDatabase(sqlDB, db.Config{
		DefaultTemplatePath: config.DefaultTemplatePath,
	})
	Setup(database)

	return database
}

func Setup(database *db.DB) {
	database.TableMap().AddTableWithName(Client{}, "clients").SetKeys(true, "Primary").ColMap("ID").SetUnique(true)
	database.TableMap().AddTableWithName(Kind{}, "kinds").SetKeys(true, "Primary").SetUniqueTogether("id", "client_id")
	database.TableMap().AddTableWithName(Receipt{}, "receipts").SetKeys(true, "Primary").SetUniqueTogether("user_guid", "client_id", "kind_id")
	database.TableMap().AddTableWithName(Unsubscribe{}, "unsubscribes").SetKeys(true, "Primary").SetUniqueTogether("user_id", "client_id", "kind_id")
	database.TableMap().AddTableWithName(GlobalUnsubscribe{}, "global_unsubscribes").SetKeys(true, "Primary").ColMap("UserID").SetUnique(true)
	database.TableMap().AddTableWithName(Template{}, "templates").SetKeys(true, "Primary").ColMap("Name").SetUnique(true)
	database.TableMap().AddTableWithName(Message{}, "messages").SetKeys(false, "ID")
}
