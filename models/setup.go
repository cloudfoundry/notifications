package models

import (
	"github.com/cloudfoundry-incubator/notifications/db"
	v1models "github.com/cloudfoundry-incubator/notifications/v1/models"
	v2models "github.com/cloudfoundry-incubator/notifications/v2/models"
)

func Setup(database *db.DB) {
	database.TableMap().AddTableWithName(v1models.Client{}, "clients").SetKeys(true, "Primary").ColMap("ID").SetUnique(true)
	database.TableMap().AddTableWithName(v1models.Kind{}, "kinds").SetKeys(true, "Primary").SetUniqueTogether("id", "client_id")
	database.TableMap().AddTableWithName(v1models.Receipt{}, "receipts").SetKeys(true, "Primary").SetUniqueTogether("user_guid", "client_id", "kind_id")
	database.TableMap().AddTableWithName(v1models.Unsubscribe{}, "unsubscribes").SetKeys(true, "Primary").SetUniqueTogether("user_id", "client_id", "kind_id")
	database.TableMap().AddTableWithName(v1models.GlobalUnsubscribe{}, "global_unsubscribes").SetKeys(true, "Primary").ColMap("UserID").SetUnique(true)
	database.TableMap().AddTableWithName(v1models.Template{}, "templates").SetKeys(true, "Primary").ColMap("Name").SetUnique(true)
	database.TableMap().AddTableWithName(v1models.Message{}, "messages").SetKeys(false, "ID")
	database.TableMap().AddTableWithName(v2models.Sender{}, "senders").SetKeys(false, "ID").SetUniqueTogether("name", "client_id")
	database.TableMap().AddTableWithName(v2models.CampaignType{}, "campaign_types").SetKeys(false, "ID").SetUniqueTogether("name", "sender_id")
}
