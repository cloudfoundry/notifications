package models

import "github.com/cloudfoundry-incubator/notifications/db"

func Setup(database *db.DB) {
	database.TableMap().AddTableWithName(Client{}, "clients").SetKeys(true, "Primary").ColMap("ID").SetUnique(true)
	database.TableMap().AddTableWithName(Kind{}, "kinds").SetKeys(true, "Primary").SetUniqueTogether("id", "client_id")
	database.TableMap().AddTableWithName(Receipt{}, "receipts").SetKeys(true, "Primary").SetUniqueTogether("user_guid", "client_id", "kind_id")
	database.TableMap().AddTableWithName(Unsubscribe{}, "unsubscribes").SetKeys(true, "Primary").SetUniqueTogether("user_id", "client_id", "kind_id")
	database.TableMap().AddTableWithName(GlobalUnsubscribe{}, "global_unsubscribes").SetKeys(true, "Primary").ColMap("UserID").SetUnique(true)
	database.TableMap().AddTableWithName(Template{}, "templates").SetKeys(true, "Primary").ColMap("Name").SetUnique(true)
	database.TableMap().AddTableWithName(Message{}, "messages").SetKeys(false, "ID")
	database.TableMap().AddTableWithName(Sender{}, "senders").SetKeys(false, "ID").SetUniqueTogether("name", "client_id")
	database.TableMap().AddTableWithName(CampaignType{}, "campaign_types").SetKeys(false, "ID").SetUniqueTogether("name", "sender_id")
}
