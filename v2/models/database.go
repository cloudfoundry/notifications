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
	database.TableMap().AddTableWithName(Sender{}, "senders").SetKeys(false, "ID").SetUniqueTogether("name", "client_id")
	database.TableMap().AddTableWithName(CampaignType{}, "campaign_types").SetKeys(false, "ID").SetUniqueTogether("name", "sender_id")
	database.TableMap().AddTableWithName(Template{}, "v2_templates").SetKeys(false, "ID").SetUniqueTogether("name", "client_id")
	database.TableMap().AddTableWithName(Campaign{}, "campaigns").SetKeys(false, "ID")
	database.TableMap().AddTableWithName(Message{}, "messages").SetKeys(false, "ID")
	database.TableMap().AddTableWithName(Unsubscriber{}, "unsubscribers").SetKeys(false, "ID").SetUniqueTogether("campaign_type_id", "user_guid")
}
